package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Strcture for fake data
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

type Submission struct {
	ID             int             `json:"id"`
	User           User            `json:"user"`
	Challenge      CodingChallenge `json:"challenge"`
	TimeToComplete int             `json:"time_to_complete"`
	SubmittedAt    time.Time       `json:"submitted_at"`
}

// Sample data pools
var topics = []string{
	"Sorting Algorithms", "REST API Design", "Concurrency", "Data Structures",
	"Graph Theory", "Dynamic Programming", "File I/O", "Unit Testing",
}

var users = []User{
	{1, "viserys", "viserys@example.com"},
	{2, "daemon", "daemon@example.com"},
	{3, "corlys", "corlys@example.com"},
	{4, "rheanys", "rhenys@example.com"},
	{5, "leanor", "leanor@example.com"},
	{6, "aemond", "aemond@example.com"},
	{7, "otto", "otto@example.com"},
	{8, "jaeherys", "jaeherys@example.com"},
	{9, "bealor", "bealor@example.com"},
	{10, "mealor", "mealor@example.com"},
}

var challenges = []CodingChallenge{
	{1, 100, "Sorting Algorithms"},
	{2, 150, "Graph Theory"},
	{3, 120, "REST API Design"},
	{4, 180, "Concurrency"},
	{5, 200, "Dynamic Programming"},
}

// Function to generate fake (structured) data every 1–10 seconds
func generateFakeData() {
	rand.Seed(time.Now().UnixNano())
	submissionID := 1

	for {
		// Randomly select a user and challenge
		user := users[rand.Intn(len(users))]
		challenge := challenges[rand.Intn(len(challenges))]

		// Random time to complete (30–600 seconds)
		timeToComplete := rand.Intn(570) + 30

		// Create submission
		submission := Submission{
			ID:             submissionID,
			User:           user,
			Challenge:      challenge,
			TimeToComplete: timeToComplete,
			SubmittedAt:    time.Now(),
		}

		// Print or send somewhere
		fmt.Printf("[%s] New Submission: %+v\n", submission.SubmittedAt.Format(time.RFC3339), submission)

		// Send to Python backend
		sendToPython(submission)

		submissionID++

		// Wait random duration between 1–10 seconds
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)

	}
}

func main() {
	fmt.Println("Generating fake submissions...")
	generateFakeData()
}

func sendToPython(submission Submission) {
	jsonData, _ := json.Marshal(submission)
	http.Post("http://localhost:8000/submit", "application/json", bytes.NewBuffer(jsonData))
}
