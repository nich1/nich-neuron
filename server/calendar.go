package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

const calendarFile = "database/time/calendar/calendar.json"

// TODO: Add chronological ordering the calendar file?
type timeInstance struct {
	annual bool //Labels an event as occuring annually
	year   uint16
	month  uint8
	day    uint8 // Disregard if 0
	hour   uint8 // 25 To Disregard, Auto Disregard if day = 0
	minute uint8 // 0 To Disregard, Auto Disregard if hour = 25
}
type reminder struct {
	title string
	when  timeInstance
}
type calendarEvent struct {
	title     string
	enabled   bool
	date      timeInstance
	reminders []timeInstance
}

func getCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray []calendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Return json
	context.IndentedJSON(http.StatusOK, calendarArray)

}

func addCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray []calendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to calendarEvent struct
	var convertedJson calendarEvent
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

	// TODO VERY IMPORTANT ADD REMINDERS HERE- needs to be tested now

	// Add payloads date and reminders to an array,
	// add dates to server calendar
	var eventQueue []timeInstance
	for i := 0; i < len(eventQueue); i++ {
		// Open / create proper directories and file
		var eventYear string
		if eventQueue[i].annual == true {
			eventYear = "annual"
		} else {
			eventYear = string(eventQueue[i].year)
		}
		if _, dirErr1 := os.Stat("database/time/calendar/" + eventYear); dirErr1 != nil {
			if os.IsNotExist(dirErr1) {
				// Dir does not exist, create year dir
				dirErr2 := os.Mkdir(("database/time/calendar/" + eventYear), 0755)
				if dirErr2 != nil {
					fmt.Println("Error in making directory")
				}
			}
		}
		var eventMonth string
		eventMonth = string(eventQueue[i].month)
		if _, dirErr3 := os.Stat("database/time/calendar/" + eventYear + "/" + eventMonth); dirErr3 != nil {
			if os.IsNotExist(dirErr3) {
				// Dir does not exist, create month dir
				dirErr4 := os.Mkdir(("database/time/calendar/" + eventYear + "/" + eventMonth), 0755)
				if dirErr4 != nil {
					fmt.Println("Error in making directory")
				}
			}
		}
		var filename string
		filename = string(eventQueue[i].day) + ".json"
		if _, fileErr1 := os.Stat("database/time/calendar/" + eventYear + "/" + eventMonth + "/" + filename); fileErr1 != nil {
			if os.IsNotExist(fileErr1) {
				// File does not exist, create file
				_, fileErr2 := os.Create(("database/time/calendar/" + eventYear + "/" + eventMonth + "/" + filename))
				if fileErr2 != nil {
					fmt.Println("Error in making file")
				}
			}
		}
		// Read in targeted file for calendar
		fileContent, fileErr := ioutil.ReadFile("database/time/calendar/" + eventYear + "/" + eventMonth + "/" + filename)
		if fileErr != nil {
			fmt.Println(fileErr)
			return
		}

		// Convert file into an array of calendarEvent
		var eventArray []timeInstance
		unmarshalErr1 := json.Unmarshal([]byte(fileContent), &eventArray)
		if unmarshalErr1 != nil {
			fmt.Println(unmarshalErr1)
			return
		}
		eventArray = append(eventArray, eventQueue[i])

		// Convert array back to json
		writeJson2, writeJsonErr := json.Marshal(eventArray)
		if writeJsonErr != nil {
			fmt.Println(writeJsonErr)
		}

		// Write eventArray to schedules.json file
		writeErr := ioutil.WriteFile("database/time/calendar/"+eventYear+"/"+eventMonth+"/"+filename, writeJson2, 0644)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
	}

}

func deleteCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray []calendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to calendarEvent struct
	var convertedJson calendarEvent
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Delete convertedJson from calendarArray
	for i := 0; i < len(calendarArray); i++ {
		if reflect.DeepEqual(convertedJson, calendarArray[i]) /*convertedJson == calendarArray[i]*/ {
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

	// TODO VERY IMPORTANT ADD REMINDERS HERE- needs to be tested now

	// Add payloads date and reminders to an array,
	// removes dates from server calendar
	var eventQueue []timeInstance
	for i := 0; i < len(eventQueue); i++ {
		var eventYear string
		if eventQueue[i].annual == true {
			eventYear = "annual"
		} else {
			eventYear = string(eventQueue[i].year)
		}
		var filename string
		filename = string(eventQueue[i].day) + ".json"
		// Read in targeted file for calendar
		fileContent, fileErr := ioutil.ReadFile("database/time/calendar/" + eventYear + "/" + string(eventQueue[i].month) + "/" + filename)
		if fileErr != nil {
			fmt.Println(fileErr)
			return
		}

		// Convert file into an array of calendarEvent
		var eventArray []timeInstance
		unmarshalErr1 := json.Unmarshal([]byte(fileContent), &eventArray)
		if unmarshalErr1 != nil {
			fmt.Println(unmarshalErr1)
			return
		}
		for j := 0; j < len(eventArray); j++ {
			if eventQueue[i] == eventArray[j] {
				eventArray = append(eventArray[:j], eventArray[j+1:]...)
				break
			}
		}

		// Convert array back to json
		writeJson2, writeJsonErr := json.Marshal(eventArray)
		if writeJsonErr != nil {
			fmt.Println(writeJsonErr)
		}

		// Write eventArray to schedules.json file
		writeErr := ioutil.WriteFile("database/time/calendar/"+eventYear+string(eventQueue[i].month)+"/"+filename, writeJson2, 0644)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
	}
}

func putCalendar(context *gin.Context) {
	// Create two arrays for comparison
	var oldArray []timeInstance
	var newArray []timeInstance
	oldArray = append(oldArray)

}

/*
func putCalendar(context *gin.Context) {
	// Assumes that an array is passed in. The first object in the array (both are calendarEvent struct-able)
	// is the object to be replaced, and the second is the new object data

	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray []calendarEvent
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to array of Schedule
	var convertedJson [2]calendarEvent
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Change matching object in schedulesArray to the new data
	for i := 0; i < len(calendarArray); i++ {
		if reflect.DeepEqual(convertedJson[0], calendarArray[i]) /*convertedJson[0] == calendarArray[i] {
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
*/
