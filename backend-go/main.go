package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CodingChallenge struct {
	ID     int    `json:"id"`
	Points int    `json:"points"`
	Topic  string `json:"topic"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ScoreUpdate struct {
	Username string `json:"username"`
	Points   int    `json:"points"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type Leaderboard struct {
	Entities []LeaderboardEntry `json:"entities"`
}

type Analytics struct {
	Total_Submissions int            `json:"total_submissions"`
	Avg_Completion    int            `json:"avg_completion"`
	Topics            map[string]int `json:"topics"`
}

var (
	GlobalLeaderboard = Leaderboard{Entities: []LeaderboardEntry{}}
	mu                sync.Mutex
)

// upgrade http server to websockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var broadcast = make(chan struct{})

//

func main() {

	// HTTP server with Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Leaderboard": GlobalLeaderboard})
	})

	router.GET("/analytics", func(c *gin.Context) {

		// analytics object
		var analytics Analytics

		// bind to json

		// fetch analytics from python service
		resp, err := http.Get("http://backend-python:8000/analytics")
		if err != nil {
			fmt.Println("Error, ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach analytics service."})
			return
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&analytics)
		if err != nil {
			fmt.Println("Error decoding analytics:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode analytics"})
			return
		}

		// return analytics to frontend
		c.JSON(http.StatusOK, analytics)
	})

	// python service calls this endpoint
	router.POST("/score", func(c *gin.Context) {

		var score ScoreUpdate

		err := c.ShouldBindJSON(&score)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received score: %+v\n", score)

		// update leaderboard based on new score
		// add points to user whose score was received (maybe add them to leaderbaord if not already present)
		addOrUpdateLeaderboard(score.Username, score.Points)

		// reshuffle leaderboard based on points
		SortLeaderboard()

		c.JSON(200, gin.H{"status": "Score received"})
	})

	// Websockets server
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		defer conn.Close()

		for {
			// <-broadcast
			data, _ := json.Marshal(GlobalLeaderboard)
			conn.WriteMessage(websocket.TextMessage, data)
			time.Sleep(time.Second)
		}
	})

	router.Run()
}

func addOrUpdateLeaderboard(username string, points int) {
	mu.Lock()
	defer mu.Unlock()

	// Check if user already exists in leaderboard
	for i, entry := range GlobalLeaderboard.Entities {
		if entry.Username == username {
			// Update points
			GlobalLeaderboard.Entities[i].Score += points
			return
		}
	}

	// If user doesn't exist, add them to the leaderboard
	newEntry := LeaderboardEntry{
		Username: username,
		Score:    points,
	}
	GlobalLeaderboard.Entities = append(GlobalLeaderboard.Entities, newEntry)

	// broadcast <- struct{}{}
}

func SortLeaderboard() {
	mu.Lock()
	defer mu.Unlock()

	sort.Slice(GlobalLeaderboard.Entities, func(i, j int) bool {
		return GlobalLeaderboard.Entities[i].Score > GlobalLeaderboard.Entities[j].Score
	})

}
