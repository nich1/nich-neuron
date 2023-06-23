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
	Title  string `json:"title"`
	Annual bool   `json:"annual"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Day    int    `json:"day"`
	Hour   int    `json:"hour"`
	Minute int    `json:"minute"`
}
type timeInstances []struct {
	Title  string `json:"title"`
	Annual bool   `json:"annual"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Day    int    `json:"day"`
	Hour   int    `json:"hour"`
	Minute int    `json:"minute"`
}

/*
type calendarEvent struct {
	title     string
	enabled   bool
	date      timeInstance
	reminders []timeInstance
} */

type calendarEvent struct {
	Title   string `json:"title"`
	Enabled bool   `json:"enabled"`
	Date    struct {
		Title  string `json:"title"`
		Annual bool   `json:"annual"`
		Year   int    `json:"year"`
		Month  int    `json:"month"`
		Day    int    `json:"day"`
		Hour   int    `json:"hour"`
		Minute int    `json:"minute"`
	} `json:"date"`
	Reminders []struct {
		Title  string `json:"title"`
		Annual bool   `json:"annual"`
		Year   int    `json:"year"`
		Month  int    `json:"month"`
		Day    int    `json:"day"`
		Hour   int    `json:"hour"`
		Minute int    `json:"minute"`
	} `json:"reminders"`
}

type calendarEvents []struct {
	Title   string `json:"title"`
	Enabled bool   `json:"enabled"`
	Date    struct {
		Title  string `json:"title"`
		Annual bool   `json:"annual"`
		Year   int    `json:"year"`
		Month  int    `json:"month"`
		Day    int    `json:"day"`
		Hour   int    `json:"hour"`
		Minute int    `json:"minute"`
	} `json:"date"`
	Reminders []struct {
		Title  string `json:"title"`
		Annual bool   `json:"annual"`
		Year   int    `json:"year"`
		Month  int    `json:"month"`
		Day    int    `json:"day"`
		Hour   int    `json:"hour"`
		Minute int    `json:"minute"`
	} `json:"reminders"`
}

func getCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray calendarEvents
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

	// Convert file into array calendarEvents
	var calendarArray calendarEvents
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

	// Add payloads date and reminders to an array,
	// add dates to server calendar
	var eventQueue timeInstances

	// Populate eventQueue
	var time0 timeInstance
	time0.Title = convertedJson.Title
	time0.Annual = convertedJson.Date.Annual
	time0.Year = convertedJson.Date.Year
	time0.Month = convertedJson.Date.Month
	time0.Day = convertedJson.Date.Day
	time0.Hour = convertedJson.Date.Hour
	time0.Minute = convertedJson.Date.Minute
	eventQueue = append(eventQueue, time0)
	// Populate reminders
	for i := 0; i < len(convertedJson.Reminders); i++ {
		var time1 timeInstance
		time1.Title = convertedJson.Reminders[i].Title
		time1.Annual = convertedJson.Reminders[i].Annual
		time1.Year = convertedJson.Reminders[i].Year
		time1.Month = convertedJson.Reminders[i].Month
		time1.Day = convertedJson.Reminders[i].Day
		time1.Hour = convertedJson.Reminders[i].Hour
		time1.Minute = convertedJson.Reminders[i].Minute
		eventQueue = append(eventQueue, time1)

	}

	for i := 0; i < len(eventQueue); i++ {
		// Open / create proper directories and file
		fmt.Println(eventQueue[i])
		fmt.Println(i)
		var eventYear string
		if eventQueue[i].Annual == true {
			eventYear = "annual"
		} else {
			eventYear = strconv.FormatUint(uint64(eventQueue[i].Year), 10)
		}

		// Create year folder if does not exist
		if _, dirErr1 := os.Stat("database/time/calendar/" + eventYear); dirErr1 != nil {
			if os.IsNotExist(dirErr1) {
				// Dir does not exist, create year dir
				dirErr2 := os.Mkdir(("database/time/calendar/" + eventYear), 0755)
				if dirErr2 != nil {
					fmt.Println("Error in making directory")
				}
			}
		}

		// Create month folder if does not exist
		var eventMonth string
		eventMonth = strconv.FormatUint(uint64(eventQueue[i].Month), 10)
		if _, dirErr3 := os.Stat("database/time/calendar/" + eventYear + "/" + eventMonth); dirErr3 != nil {
			if os.IsNotExist(dirErr3) {
				// Dir does not exist, create month dir
				dirErr4 := os.Mkdir(("database/time/calendar/" + eventYear + "/" + eventMonth), 0755)
				if dirErr4 != nil {
					fmt.Println("Error in making directory")
				}
			}
		}

		// Create file named the day if does not exist
		var filename string
		filename = strconv.FormatUint(uint64(eventQueue[i].Day), 10) + ".json"
		if _, fileErr1 := os.Stat("database/time/calendar/" + eventYear + "/" + eventMonth + "/" + filename); fileErr1 != nil {
			if os.IsNotExist(fileErr1) {
				// File does not exist, create file

				var eventArray timeInstances
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
		} else {

			// Read in targeted file for calendar
			fileContent, fileErr := ioutil.ReadFile("database/time/calendar/" + eventYear + "/" + eventMonth + "/" + filename)
			if fileErr != nil {
				fmt.Println(fileErr)
				return
			}

			// Convert file into an array of calendarEvent
			var eventArray []timeInstance
			unmarshalErr1 := json.Unmarshal(fileContent, &eventArray)
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

}

func deleteCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray calendarEvents
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
		var tempEvent calendarEvent
		tempEvent.Title = calendarArray[i].Title
		tempEvent.Enabled = calendarArray[i].Enabled
		tempEvent.Date = calendarArray[i].Date
		tempEvent.Reminders = calendarArray[i].Reminders
		var check1, check2, check3, check4 bool
		check1 = tempEvent.Title == convertedJson.Title
		check2 = tempEvent.Enabled == convertedJson.Enabled
		check3 = tempEvent.Date == convertedJson.Date
		check4 = reflect.DeepEqual(tempEvent.Reminders, convertedJson.Reminders)

		if check1 && check2 && check3 && check4 {

			calendarArray = append(calendarArray[:i], calendarArray[i+1:]...)
			break
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

	// Add payloads date and reminders to an array,
	// removes dates from server calendar
	var eventQueue timeInstances

	// Populate eventQueue
	var time0 timeInstance
	time0.Title = convertedJson.Title
	time0.Annual = convertedJson.Date.Annual
	time0.Year = convertedJson.Date.Year
	time0.Month = convertedJson.Date.Month
	time0.Day = convertedJson.Date.Day
	time0.Hour = convertedJson.Date.Hour
	time0.Minute = convertedJson.Date.Minute
	eventQueue = append(eventQueue, time0)
	// Populate reminders
	for i := 0; i < len(convertedJson.Reminders); i++ {
		var time1 timeInstance
		time1.Title = convertedJson.Reminders[i].Title
		time1.Annual = convertedJson.Reminders[i].Annual
		time1.Year = convertedJson.Reminders[i].Year
		time1.Month = convertedJson.Reminders[i].Month
		time1.Day = convertedJson.Reminders[i].Day
		time1.Hour = convertedJson.Reminders[i].Hour
		time1.Minute = convertedJson.Reminders[i].Minute
		eventQueue = append(eventQueue, time1)

	}

	for i := 0; i < len(eventQueue); i++ {
		var eventYear string
		if eventQueue[i].Annual == true {
			eventYear = "annual"
		} else {
			eventYear = strconv.FormatUint(uint64(eventQueue[i].Year), 10)
		}
		fmt.Println("Debug 1")
		var filename string
		filename = strconv.FormatUint(uint64(eventQueue[i].Day), 10) + ".json"
		// Read in targeted file for calendar
		fileContent, fileErr := ioutil.ReadFile("database/time/calendar/" + eventYear + "/" + strconv.FormatUint(uint64(eventQueue[i].Month), 10) + "/" + filename)
		if fileErr != nil {
			fmt.Println(fileErr)
			return
		}
		fmt.Println("Debug 2")

		// Convert file into an array of calendarEvent
		var eventArray []timeInstance
		unmarshalErr1 := json.Unmarshal([]byte(fileContent), &eventArray)
		if unmarshalErr1 != nil {
			fmt.Println(unmarshalErr1)
			return
		}
		fmt.Println("Debug 3")

		fmt.Println(eventArray)

		for j := 0; j < len(eventArray); j++ {
			if eventQueue[i] == eventArray[j] {
				fmt.Println("MATCH")
				eventArray = append(eventArray[:j], eventArray[j+1:]...)
				break
			}
		}
		fmt.Println("Debug 4")

		// Convert array back to json
		writeJson2, writeJsonErr := json.Marshal(eventArray)
		if writeJsonErr != nil {
			fmt.Println(writeJsonErr)
		}

		fmt.Println("Debug 5")

		// Write eventArray to schedules.json file
		writeErr := ioutil.WriteFile("database/time/calendar/"+eventYear+"/"+strconv.FormatUint(uint64(eventQueue[i].Month), 10)+"/"+filename, writeJson2, 0644)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
	}
}

func putCalendar(context *gin.Context) {
	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of calendarEvent
	var calendarArray calendarEvents
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &calendarArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}

	// Convert rawJson to calendarEvents struct
	// ASSUMES FIRST ITEM IS TO REPLACE, AND SECOND ITEM IS TO REPLACE WITH
	var convertedJson calendarEvents
	if unmarshalErr2 := context.BindJSON(&convertedJson); unmarshalErr2 != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(unmarshalErr2)
		return
	}

	// Replace matching calendarArray index with the replacement data
	for i := 0; i < len(calendarArray); i++ {
		var tempEvent calendarEvent
		tempEvent.Title = calendarArray[i].Title
		tempEvent.Enabled = calendarArray[i].Enabled
		tempEvent.Date = calendarArray[i].Date
		tempEvent.Reminders = calendarArray[i].Reminders
		var check1, check2, check3, check4 bool
		check1 = tempEvent.Title == convertedJson[0].Title
		check2 = tempEvent.Enabled == convertedJson[0].Enabled
		check3 = tempEvent.Date == convertedJson[0].Date
		check4 = reflect.DeepEqual(tempEvent.Reminders, convertedJson[0].Reminders)

		if check1 && check2 && check3 && check4 {

			calendarArray[i] = convertedJson[1]
			break
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

	// EVERYTHING WORKS TIL HERE LOL

	// Create two arrays for comparison
	var oldArray timeInstances
	var newArray timeInstances

	// Add old data (first item) to oldArray which keeps track of all the time instances
	oldArray = append(oldArray, convertedJson[0].Date)
	for i := 0; i < len(convertedJson[0].Reminders); i++ {
		oldArray = append(oldArray, convertedJson[0].Reminders[i])
	}
	// Add new data (second item) to newArray which keeps track of all the time instances
	newArray = append(newArray, convertedJson[1].Date)
	for i := 0; i < len(convertedJson[1].Reminders); i++ {
		newArray = append(newArray, convertedJson[1].Reminders[i])
	}

	// For every old timeInstance, compare and update data to cooresponding new timeInstance
	for i := 0; i < len(oldArray); i++ {
		var yearOld string
		if newArray[i].Annual {
			yearOld = "annual"
		} else {
			yearOld = strconv.FormatUint(uint64(oldArray[i].Year), 10)
		}
		var yearNew string
		if newArray[i].Annual {
			yearNew = "annual"
		} else {
			yearNew = strconv.FormatUint(uint64(oldArray[i].Year), 10)
		}
		oldFileName := "database/time/calendar/" + yearOld + "/" + strconv.FormatUint(uint64(oldArray[i].Month), 10) + "/" + strconv.FormatUint(uint64(oldArray[i].Day), 10) + ".json"
		newFileName := "database/time/calendar/" + yearNew + "/" + strconv.FormatUint(uint64(newArray[i].Month), 10) + "/" + strconv.FormatUint(uint64(newArray[i].Day), 10) + ".json"

		// If filenames are the same, then we can simply modify the object in the file
		if oldFileName == newFileName {
			// Read in database file for calendar
			fileContent, fileErr := ioutil.ReadFile(oldFileName)
			if fileErr != nil {
				fmt.Println(fileErr)
				return
			}
			// Convert file into an array of timeInstance
			var dayArray timeInstances
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
			var arrayWithDel timeInstances
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
			// Before we read the file, we need to check if directories
			// exist (year and month), and if not, create them
			if yearNew != "annual" {
				yearNew = strconv.FormatUint(uint64(newArray[i].Year), 10)
				// Create year folder if does not exist
				if _, dirErr1 := os.Stat("database/time/calendar/" + yearNew); dirErr1 != nil {
					if os.IsNotExist(dirErr1) {
						// Dir does not exist, create year dir
						dirErr2 := os.Mkdir(("database/time/calendar/" + yearNew), 0755)
						if dirErr2 != nil {
							fmt.Println("Error in making directory")
						}
					}
				}
			}
			var month string = strconv.FormatUint(uint64(newArray[i].Month), 10)
			// Create month folder if does not exist
			if _, dirErr3 := os.Stat("database/time/calendar/" + yearNew + "/" + month); dirErr3 != nil {
				if os.IsNotExist(dirErr3) {
					// Dir does not exist, create month dir
					dirErr4 := os.Mkdir(("database/time/calendar/" + yearNew + "/" + month), 0755)
					if dirErr4 != nil {
						fmt.Println("Error in making directory")
					}
				}
			}

			if _, fileErr1 := os.Stat(newFileName); fileErr1 != nil {
				if os.IsNotExist(fileErr1) {
					// File does not exist, create file

					var eventArray timeInstances
					eventArray = append(eventArray, newArray[i])
					// Convert array back to json
					writeJson2, writeJsonErr := json.Marshal(eventArray)
					if writeJsonErr != nil {
						fmt.Println(writeJsonErr)
					}

					// Write eventArray to schedules.json file
					writeErr := ioutil.WriteFile(newFileName, writeJson2, 0644)
					if writeErr != nil {
						fmt.Println(writeErr)
					}

				}
			} else { // File exists, can read

				// Read in file
				fileContent2, fileErr2 := ioutil.ReadFile(newFileName)
				if fileErr2 != nil {
					fmt.Println(fileErr2)
					return
				}
				// Convert file into an array of timeInstances
				var arrayToAdd timeInstances
				unmarshalErr4 := json.Unmarshal(fileContent2, &arrayToAdd)
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

}
