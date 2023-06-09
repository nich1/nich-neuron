package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const calendarFile = "database/time/schedules/calendar.json"

// TODO: Add chronological ordering the calendar file?

type CalendarEvent struct {
	Title            string
	Enabled          bool
	Annual           bool
	Year             uint16
	Month            uint8
	Day              uint8  // 0 To Disregard (ie, to label a month as an event like Black History month for example)
	Hour             int8   // Negative to disregard
	Minute           int8   // 0 To Disregard
	DaysInAdvance    uint16 // For reminders 	If days and minutes are 0, no reminder
	MinutesInAdvance uint16 // For reminders
}

func getCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of CalendarEvent
	var calendarArray []CalendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Return json
	context.IndentedJSON(http.StatusOK, calendarArray)

}

// TODO VERY IMPORTANT ADD REMINDERS HERE
func addCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of CalendarEvent
	var calendarArray []CalendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to CalendarEvent struct
	var convertedJson CalendarEvent
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Add convertedJson to calendarArray
	calendarArray = append(calendarArray, convertedJson)

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(calendarArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write calendarArray to schedules.json file
	writeErr := ioutil.WriteFile(calendarFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}

	// TODO VERY IMPORTANT ADD REMINDERS HERE
}

// TODO VERY IMPORTANT ADD REMINDERS HERE
func deleteCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of CalendarEvent
	var calendarArray []CalendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to CalendarEvent struct
	var convertedJson CalendarEvent
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Delete convertedJson from calendarArray
	for i := 0; i < len(calendarArray); i++ {
		if convertedJson == calendarArray[i] {
			calendarArray = append(calendarArray[:i], calendarArray[i+1:]...)
		}
	}

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(calendarArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write calendarArray to calendar.json file
	writeErr := ioutil.WriteFile(calendarFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}
	// TODO VERY IMPORTANT ADD REMINDERS HERE
}

// TODO VERY IMPORTANT ADD REMINDERS HERE
func putCalendar(context *gin.Context) {
	// Assumes that an array is passed in. The first object in the array (both are CalendarEvent struct-able)
	// is the object to be replaced, and the second is the new object data

	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of CalendarEvent
	var calendarArray []CalendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to array of Schedule
	var convertedJson [2]CalendarEvent
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Change matching object in schedulesArray to the new data
	for i := 0; i < len(calendarArray); i++ {
		if convertedJson[0] == calendarArray[i] {
			calendarArray[i] = convertedJson[1]
		}
	}

	// Convert array back to json
	writeJson, writeJsonErr := json.Marshal(calendarArray)
	if writeJsonErr != nil {
		fmt.Println(writeJsonErr)
	}

	// Write schedulesArray to schedules.json file
	writeErr := ioutil.WriteFile(calendarFile, writeJson, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
	}
	// TODO VERY IMPORTANT ADD REMINDERS HERE
}
