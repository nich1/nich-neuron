package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func httpListener() {
	// Generate test files
	generateTestJson()

	// Initialize server
	server := gin.Default()

	// Routing For Notes
	server.GET("/nich-neuron/notes", getNotes)
	server.POST("/nich-neuron/notes", createNote)
	server.GET("/nich-neuron/notes/:id", noteById)
	server.DELETE("/nich-neuron/notes/:id", deleteNoteById)
	server.PUT("/nich-neuron/notes", putNote)

	// Routing For Time
	server.GET("/nich-neuron/schedule", getSchedule)
	server.POST("/nich-neuron/schedule", addSchedule)
	server.DELETE("/nich-neuron/schedule", deleteSchedule)
	//server.PUT("/nich-neuron/schedule", putSchedule)

	server.Run(":8080")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go httpListener()
	wg.Add(1)
	go timeListener()
	wg.Wait()
}
