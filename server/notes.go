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

// When the server starts, these files are automatically created if they don't already exist. This is used for testing
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

// Function called during a get + /id request. Returns json of specified id
func noteById(context *gin.Context) {
	id := context.Param("id")

	fileList, err := ioutil.ReadDir(noteDatabase)

	// Send response back saying something went wrong
	if err != nil {
		context.IndentedJSON(http.StatusConflict, "Note does not exist")
		return
	}

	for _, file := range fileList {
		if file.Name() == id+".json" {
			fmt.Println(`Retreiving file "` + file.Name() + `"`)
			fileContent, _ := ioutil.ReadFile(noteDatabase + id + ".json")
			var data interface{}
			err2 := json.Unmarshal(fileContent, &data)
			// Send response back saying something went wrong
			if err2 != nil {
				context.IndentedJSON(http.StatusConflict, "Error in converting data")
				return
			}
			fmt.Println(data)
			context.IndentedJSON(http.StatusOK, data)

			return
		}
	}

	fmt.Println("File not found")

}

// Called during basic get request. Returns json array of all json objects in the database
func getNotes(context *gin.Context) {
	fileList, err := ioutil.ReadDir(noteDatabase)

	// Send response back saying something went wrong
	if err != nil {
		context.IndentedJSON(http.StatusConflict, "Database directory error")
		return
	}

	var arr []interface{}

	for _, file := range fileList {
		fmt.Println(`Retreiving file "` + file.Name() + `"`)
		fileContent, _ := ioutil.ReadFile(noteDatabase + file.Name())
		var data interface{}
		err2 := json.Unmarshal(fileContent, &data)
		fmt.Println(data)

		// Send response back saying something went wrong
		if err2 != nil {
			context.IndentedJSON(http.StatusConflict, "Error in converting data")
			return
		}
		arr = append(arr, data)

	}

	context.IndentedJSON(http.StatusOK, arr)

}

// Called during post request, creates new json file and fills it with payload
func createNote(context *gin.Context) {
	var newNote note

	// Send response back saying something went wrong
	if err := context.BindJSON(&newNote); err != nil {
		context.IndentedJSON(http.StatusConflict, "Error in converting data")
		return
	}

	// Check if file exists already
	if _, err := ioutil.ReadFile(noteDatabase + newNote.ID + ".json"); err == nil {
		context.IndentedJSON(http.StatusConflict, "File already exists")
		return
	}

	file, err := os.Create(noteDatabase + newNote.ID + ".json")

	if err != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured in creating database file")
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
// Called by delete /id request. Deletes json file in database
func deleteNoteById(context *gin.Context) {
	id := context.Param("id")
	fileList, err := ioutil.ReadDir(noteDatabase)
	if err != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured in reading database directory")
		fmt.Println("ERROR from delete request: " + err.Error())
		return
	}

	for _, file := range fileList {
		if file.Name() == id+".json" {
			fmt.Println(`Deleted file "` + file.Name() + `"`)
			os.Remove(noteDatabase + id + ".json")
			return
		}
	}

	fmt.Println(`File "` + id + `.json" does not exist in the database`)
	context.IndentedJSON(http.StatusConflict, "Note not found")
}

// Add status return
// Called by put request. Changes payload if file exists, otherwise creates file + injects payload
func putNote(context *gin.Context) {
	// create note w payload
	var newNote note

	if err := context.BindJSON(&newNote); err != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		return
	}

	id := newNote.ID

	os.Remove(noteDatabase + id + ".json")
	file, err2 := os.Create(noteDatabase + id + ".json")

	if err2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured in creating database file")
		return
	}

	// I believe this can be improved upon in the future
	file.WriteString(
		`{"id": "` + newNote.ID + `",` + `"title": "` + newNote.Title + `",` + `"body": "` + newNote.Body + `"}`)

}
