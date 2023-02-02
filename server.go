package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var notes = []note{
	{ID: "1", Title: "Homework Plan", Body: "Do Linear Algebra on Saturday and everything else on Friday"},
	{ID: "5", Title: "Sunday Reminder", Body: "Go to church during the day!"},
	{ID: "3", Title: "Make daily checklist", Body: ""},
}

func getNotes(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, notes)
}

func createNote(context *gin.Context) {
	var newNote note

	// MAKE IT SO ID'S ARE UNIQUE

	// If there is an error
	if err := context.BindJSON(&newNote); err != nil {
		return
	}

	notes = append(notes, newNote)

	file, err := os.Create("database/notes/" + newNote.ID + ".json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	//write the json to the file
	file.WriteString(
		`{"id": "` + newNote.ID + `",` + `"title": "` + newNote.Title + `",` + `"body": "` + newNote.Body + `"}`)

	context.IndentedJSON(http.StatusCreated, newNote)

}

func deleteNoteById(context *gin.Context) {
	id := context.Param("id")

	for i, note := range notes {
		if note.ID == id {

			message := note
			notes = append(notes[:i], notes[i+1])
			context.IndentedJSON(http.StatusOK, message)
			return
		}
	}

	context.IndentedJSON(http.StatusOK, "note not found")

}

func putNote(context *gin.Context) {
	// create note w payload
	var newNote note
	if err := context.BindJSON(&newNote); err != nil {
		return
	}

	id := newNote.ID

	// if note with ID exists
	for i := range notes {
		if notes[i].ID == id {

			notes[i] = newNote

			return
		}
	}

	// create new note if a note with the ID didn't exist
	notes = append(notes, newNote)

}

func noteById(context *gin.Context) {
	id := context.Param("id")
	note, err := getNoteById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Note not found."})
		return
	}

	context.IndentedJSON(http.StatusOK, note)

}

func getNoteById(id string) (*note, error) {
	for i, note := range notes {
		if note.ID == id {
			return &notes[i], nil
		}
	}

	return nil, errors.New("note not found")
}

func main() {
	// Initialize server
	server := gin.Default()

	server.GET("/nich-neuron/notes", getNotes)
	server.POST("/nich-neuron/notes", createNote)
	server.GET("/nich-neuron/notes/:id", noteById)
	server.DELETE("/nich-neuron/delete/:id", deleteNoteById)
	server.PUT("/nich-neuron/notes", putNote)

	server.Run(":8080")

}
