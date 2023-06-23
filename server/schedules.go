package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const schedulesFile = "database/time/schedules/schedules.json"

type Schedule struct {
	Title     string
	Enabled   bool
	Hour      int8
	Minute    int8
	Sunday    bool
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
}

func putSchedule(context *gin.Context) {
	// Assumes that an array is passed in. The first object in the array (both are Schedule struct-able)
	// is the object to be replaced, and the second is the new object data

	// Read in database file for schedules
	fileContent, fileErr := ioutil.ReadFile(schedulesFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of Schedule
	var schedulesArray []Schedule
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &schedulesArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to array of Schedule
	var convertedJson [2]Schedule
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Change matching object in schedulesArray to the new data
	for i := 0; i < len(schedulesArray); i++ {
		if convertedJson[0] == schedulesArray[i] {
			schedulesArray[i] = convertedJson[1]
		}
	}

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(schedulesArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write schedulesArray to schedules.json file
	writeErr := ioutil.WriteFile(schedulesFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}
}

func getSchedule(context *gin.Context) {
	// Read in database file for schedules
	fileContent, fileErr := ioutil.ReadFile(schedulesFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of Schedule
	var schedulesArray []Schedule
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &schedulesArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Return json
	context.IndentedJSON(http.StatusOK, schedulesArray)

}

func addSchedule(context *gin.Context) {
	// Read in database file for schedules
	fileContent, fileErr := ioutil.ReadFile(schedulesFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of Schedule
	var schedulesArray []Schedule
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &schedulesArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to Schedule struct
	var convertedJson Schedule
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Add convertedJson to schedulesArray
	schedulesArray = append(schedulesArray, convertedJson)

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(schedulesArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write schedulesArray to schedules.json file
	writeErr := ioutil.WriteFile(schedulesFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}
}

func deleteSchedule(context *gin.Context) {
	// Read in database file for schedules
	fileContent, fileErr := ioutil.ReadFile(schedulesFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of Schedule
	var schedulesArray []Schedule
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &schedulesArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to Schedule struct
	var convertedJson Schedule
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Delete convertedJson from schedulesArray
	for i := 0; i < len(schedulesArray); i++ {
		if convertedJson == schedulesArray[i] {
			schedulesArray = append(schedulesArray[:i], schedulesArray[i+1:]...)
		}
	}

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(schedulesArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write schedulesArray to schedules.json file
	writeErr := ioutil.WriteFile(schedulesFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}

}
