package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

const calendarFile = "database/time/calendar/calendar.json"

// TODO: Add chronological ordering the calendar file?
type timeInstance struct {
	title  string
	annual bool //Labels an event as occuring annually
	year   uint16
	month  uint8
	day    uint8 // Disregard if 0
	hour   uint8 // 25 To Disregard, Auto Disregard if day = 0
	minute uint8 // 0 To Disregard, Auto Disregard if hour = 25
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
			eventYear = strconv.FormatUint(uint64(eventQueue[i].year), 10)
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
		eventMonth = strconv.FormatUint(uint64(eventQueue[i].month), 10)
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
		filename = strconv.FormatUint(uint64(eventQueue[i].day), 10) + ".json"
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
			eventYear = strconv.FormatUint(uint64(eventQueue[i].year), 10)
		}
		var filename string
		filename = strconv.FormatUint(uint64(eventQueue[i].day), 10) + ".json"
		// Read in targeted file for calendar
		fileContent, fileErr := ioutil.ReadFile("database/time/calendar/" + eventYear + "/" + strconv.FormatUint(uint64(eventQueue[i].month), 10) + "/" + filename)
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
		writeErr := ioutil.WriteFile("database/time/calendar/"+eventYear+strconv.FormatUint(uint64(eventQueue[i].month), 10)+"/"+filename, writeJson2, 0644)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
	}
}

func putCalendar(context *gin.Context) {
	// Create two arrays for comparison
	var oldArray []timeInstance
	var newArray []timeInstance

	// Extract payload data
	var payload []calendarEvent
	if unmarshalErr1 := context.BindJSON(&payload); unmarshalErr1 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during conversion")
		fmt.Println(unmarshalErr1)
		return
	}
	// Add old data (first item) to oldArray which keeps track of all the time instances
	oldArray = append(oldArray, payload[0].date)
	for i := 0; i < len(payload[0].reminders); i++ {
		oldArray = append(oldArray, payload[0].reminders[i])
	}
	// Add new data (second item) to newArray which keeps track of all the time instances
	newArray = append(newArray, payload[1].date)
	for i := 0; i < len(payload[1].reminders); i++ {
		newArray = append(newArray, payload[1].reminders[i])
	}

	// For every old timeInstance, compare and update data to cooresponding new timeInstance
	for i := 0; i < len(oldArray); i++ {
		oldFileName := "database/time/calendar/" + strconv.FormatUint(uint64(oldArray[i].year), 10) + "/" + strconv.FormatUint(uint64(oldArray[i].month), 10) + "/" + strconv.FormatUint(uint64(oldArray[i].day), 10) + ".json"
		newFileName := "database/time/calendar/" + strconv.FormatUint(uint64(newArray[i].year), 10) + "/" + strconv.FormatUint(uint64(newArray[i].month), 10) + "/" + strconv.FormatUint(uint64(newArray[i].day), 10) + ".json"

		// If filenames are the same, then we can simply modify the object in the file
		if oldFileName == newFileName {
			// Read in database file for calendar
			fileContent, fileErr := ioutil.ReadFile(oldFileName)
			if fileErr != nil {
				fmt.Println(fileErr)
				return
			}
			// Convert file into an array of timeInstance
			var dayArray []timeInstance
			unmarshalErr2 := json.Unmarshal([]byte(fileContent), &dayArray)
			if unmarshalErr2 != nil {
				fmt.Println(unmarshalErr2)
				return
			}
			// Find object in file, replace with equivalent newTimeInstance
			for j := 0; j < len(dayArray); j++ {
				if oldArray[i] == dayArray[j] {
					dayArray[j] = newArray[i]
				}
			}

		} else { // If filenames are not the same, we must delete from old file and create in new file
			//[1] First delete object from old file
			// Read in file
			fileContent1, fileErr1 := ioutil.ReadFile(oldFileName)
			if fileErr1 != nil {
				fmt.Println(fileErr1)
				return
			}
			// Convert file into an array of timeInstances
			var arrayWithDel []timeInstance
			unmarshalErr3 := json.Unmarshal([]byte(fileContent1), &arrayWithDel)
			if unmarshalErr3 != nil {
				fmt.Println(unmarshalErr3)
				return
			}
			// Delete match from arrayWithDel
			for j := 0; j < len(arrayWithDel); j++ {
				if oldArray[i] == arrayWithDel[j] {
					arrayWithDel = append(arrayWithDel[:j], arrayWithDel[j+1:]...)
				}
			}
			// Convert array back to json
			writeJson, writeJsonErr := json.Marshal(arrayWithDel)
			if writeJsonErr != nil {
				fmt.Println(writeJsonErr)
			}
			// Write calendarArray to calendar.json file
			writeErr := ioutil.WriteFile(oldFileName, writeJson, 0644)
			if writeErr != nil {
				fmt.Println(writeErr)
			}

			//[2] Then add object to new file
			// Read in file
			fileContent2, fileErr2 := ioutil.ReadFile(newFileName)
			if fileErr2 != nil {
				fmt.Println(fileErr2)
				return
			}
			// Convert file into an array of timeInstances
			var arrayToAdd []timeInstance
			unmarshalErr4 := json.Unmarshal([]byte(fileContent2), &arrayToAdd)
			if unmarshalErr4 != nil {
				fmt.Println(unmarshalErr4)
				return
			}
			// Add new object to arrayToAdd
			arrayToAdd = append(arrayToAdd, newArray[i])
			// Convert array back to json
			writeJson2, writeJsonErr2 := json.Marshal(arrayToAdd)
			if writeJsonErr2 != nil {
				fmt.Println(writeJsonErr2)
			}
			// Write calendarArray to calendar.json file
			writeErr2 := ioutil.WriteFile(newFileName, writeJson2, 0644)
			if writeErr2 != nil {
				fmt.Println(writeErr2)
			}
		}
	}

}
