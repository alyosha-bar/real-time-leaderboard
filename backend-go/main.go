package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Username string `json:"username"`
}

type Challenge struct {
	ID     int `json:"id"`
	Points int `json:"points"`
}

type Submission struct {
	User      User      `json:"user"`
	Challenge Challenge `json:"challenge"`
}

var (
	ctx              = context.Background()
	isDirty          = false
	isDirtyMutex     sync.Mutex // preventing race conditions
	broadcastChannel = make(chan string)
)

func main() {
	// initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	hub := NewHub()

	// WS handler
	go hub.Run()

	// start background goroutines
	go listenToRedis(rdb)
	go startTicker(rdb, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Go server starting on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
			fmt.Printf("\tGoroutines = %v\n", runtime.NumGoroutine())
			time.Sleep(5 * time.Second)
		}
	}()

}

func listenToRedis(rdb *redis.Client) {
	pubsub := rdb.Subscribe(ctx, "score_updates")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for range ch {
		isDirtyMutex.Lock()
		isDirty = true
		isDirtyMutex.Unlock()
	}
}

func startTicker(rdb *redis.Client, hub *Hub) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for range ticker.C {
		isDirtyMutex.Lock()
		if isDirty {
			// Fetch Top 10 from Redis
			vals, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()

			// Marshal into JSON
			payload, _ := json.Marshal(vals)

			// Push to the Hub's broadcast channel
			hub.broadcast <- payload

			isDirty = false
		}
		isDirtyMutex.Unlock()
	}
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// Hub logic
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the Ticker.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// WebSocket handler
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for local development (Careful in Production!)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new Client for this user
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	// Register the client with the Hub
	hub.register <- client

	// Start the writing goroutine for this specific user
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// The Hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break // Exit if connection is closed
		}
	}
}
