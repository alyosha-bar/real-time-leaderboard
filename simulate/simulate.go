package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Strcture for fake data
type CodingChallenge struct {
	ID     int
	Points int
	Topic  string
}

type User struct {
	ID       int
	Username string
	Email    string
}

type Submission struct {
	ID             int
	User           User
	Challenge      CodingChallenge
	TimeToComplete int // in seconds
	SubmittedAt    time.Time
}

// Sample data pools
var topics = []string{
	"Sorting Algorithms", "REST API Design", "Concurrency", "Data Structures",
	"Graph Theory", "Dynamic Programming", "File I/O", "Unit Testing",
}

var users = []User{
	{1, "alice", "alice@example.com"},
	{2, "bob", "bob@example.com"},
	{3, "charlie", "charlie@example.com"},
	{4, "diana", "diana@example.com"},
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

		submissionID++

		// Wait random duration between 1–10 seconds
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
	}
}

func main() {
	fmt.Println("Generating fake submissions...")
	generateFakeData()
}
