package main

import "github.com/gin-gonic/gin"

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

	// Routing For Schedules
	server.GET("/nich-neuron/schedule", getSchedule)
	server.POST("/nich-neuron/schedule", addSchedule)
	server.DELETE("/nich-neuron/schedule", deleteSchedule)
	server.PUT("/nich-neuron/schedule", putSchedule)

	// Routing For Calendar
	server.GET("/nich-neuron/calendar", getCalendar)
	server.POST("/nich-neuron/calendar", addCalendar)
	server.DELETE("/nich-neuron/calendar", deleteCalendar)
	server.PUT("/nich-neuron/calendar", putCalendar)

	// Routing For Drive
	server.MaxMultipartMemory = 8 << 20 // Max file upload is 8 MiB
	server.POST("/nich-neuron/drive", postDrive)
	server.GET("/nich-neuron/drive", getDrive)
	server.DELETE("/nich-neuron/drive", deleteDrive)

	// Start server on port 8080
	server.Run(":8080")
}
