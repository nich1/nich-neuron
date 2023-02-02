package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const noteDatabase = "database/notes/"

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

func generateTestJson() {
	file1, err1 := os.Create(noteDatabase + "1.json")
	file2, err2 := os.Create(noteDatabase + "3.json")
	file3, err3 := os.Create(noteDatabase + "5.json")

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("[Error] Could not generate test files")
	}

	file1.WriteString(`{"ID": "1", "Title": "Homework Plan", "Body": "Do Linear Algebra on Saturday and everything else on Friday"}`)
	file2.WriteString(`{"ID": "3", "Title": "Make daily checklist", "Body": ""}`)
	file3.WriteString(`{"ID": "5", "Title": "Sunday Reminder", "Body": "Go to church during the day!"}`)

}

func noteById(context *gin.Context) {
	id := context.Param("id")

	fileList, err := ioutil.ReadDir("database/notes/")

	if err != nil {
		return
	}

	for _, file := range fileList {
		if file.Name() == id+".json" {
			fmt.Println(`Retreiving file "` + file.Name() + `"`)
			fileContent, _ := ioutil.ReadFile(noteDatabase + id + ".json")
			var data interface{}
			err2 := json.Unmarshal(fileContent, &data)
			if err2 != nil {
				return
			}
			fmt.Println(data)
			context.IndentedJSON(http.StatusOK, data)

			return
		}
	}

	fmt.Println("File not found")

}

func getNotes(context *gin.Context) {
	fileList, err := ioutil.ReadDir("database/notes/")

	if err != nil {
		return
	}

	var arr []interface{}

	for _, file := range fileList {
		fmt.Println(`Retreiving file "` + file.Name() + `"`)
		fileContent, _ := ioutil.ReadFile(noteDatabase + file.Name())
		var data interface{}
		err2 := json.Unmarshal(fileContent, &data)
		fmt.Println(data)

		if err2 != nil {
			return
		}
		arr = append(arr, data)
		//context.IndentedJSON(http.StatusOK, data)

	}

	//json.Marshal(listOfObjects)

	context.IndentedJSON(http.StatusOK, arr)

}

// Needs to be redone without the array
func getNotes1(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, notes)
}

func createNote(context *gin.Context) {
	var newNote note

	// MAKE IT SO ID'S ARE UNIQUE

	// If there is an error
	if err := context.BindJSON(&newNote); err != nil {
		return
	}

	// Check if file exists already
	if _, err := ioutil.ReadFile("database/notes/" + newNote.ID + ".json"); err == nil {
		fmt.Println("ERROR: File already exists")
		return
	}

	notes = append(notes, newNote)

	file, err := os.Create("database/notes/" + newNote.ID + ".json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	//write the json to the file. please make this better in the future.
	file.WriteString(
		`{"id": "` + newNote.ID + `",` + `"title": "` + newNote.Title + `",` + `"body": "` + newNote.Body + `"}`)

	context.IndentedJSON(http.StatusCreated, newNote)

}

// Add status return
func deleteNoteById(context *gin.Context) {
	id := context.Param("id")

	fileList, err := ioutil.ReadDir("database/notes/")

	if err != nil {
		fmt.Println("ERROR from delete request: " + err.Error())
		return
	}

	for _, file := range fileList {
		if file.Name() == id+".json" {
			fmt.Println(`Deleted file "` + file.Name() + `"`)
			os.Remove("database/notes/" + id + ".json")
			return
		}
	}

	fmt.Println(`File "` + id + `.json" does not exist in the database`)

	context.IndentedJSON(http.StatusOK, "note not found")

}

// Add status return
func putNote(context *gin.Context) {
	// create note w payload
	var newNote note
	if err := context.BindJSON(&newNote); err != nil {
		return
	}

	id := newNote.ID

	os.Remove(noteDatabase + id + ".json")
	file, err2 := os.Create(noteDatabase + id + ".json")
	if err2 != nil {
		return
	}

	file.WriteString(
		`{"id": "` + newNote.ID + `",` + `"title": "` + newNote.Title + `",` + `"body": "` + newNote.Body + `"}`)

}

func main() {
	// Generate test files
	generateTestJson()

	// Initialize server
	server := gin.Default()

	server.GET("/nich-neuron/notes", getNotes)
	server.POST("/nich-neuron/notes", createNote)
	server.GET("/nich-neuron/notes/:id", noteById)
	server.DELETE("/nich-neuron/notes/:id", deleteNoteById)
	server.PUT("/nich-neuron/notes", putNote)

	server.Run(":8080")

}
