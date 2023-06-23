package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type mostRecent struct {
	year  int
	month int
	day   int
}

func cleanCalendarHistory() {
	fileContent, fileErr := ioutil.ReadFile("database/time/listener.json")
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	// Convert file into struct
	var mostRecentReadDates mostRecent
	unmarshalErr1 := json.Unmarshal(fileContent, &mostRecentReadDates)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}
	// If the year changed
	if mostRecentReadDates.year < int(time.Now().Year()) {
		// Delete year folder
		os.Remove("database/time/calendar/" + strconv.FormatUint(uint64(time.Now().Year()), 10))
		mostRecentReadDates.year = int(time.Now().Year())
		mostRecentReadDates.month = 1
		mostRecentReadDates.day = 1
	}
	// If the month changed
	if mostRecentReadDates.month < int(time.Now().Month()) {
		// Delete month folder
		os.Remove("database/time/calendar/" + strconv.FormatUint(uint64(time.Now().Year()), 10) + "/" + strconv.FormatUint(uint64(time.Now().Month()), 10))
		mostRecentReadDates.month = int(time.Now().Month())
		mostRecentReadDates.day = 1
	}
	// If the day changed
	if mostRecentReadDates.day < int(time.Now().Day()) {
		// Delete day file
		os.Remove("database/time/calendar/" + strconv.FormatUint(uint64(time.Now().Year()), 10) + "/" + strconv.FormatUint(uint64(time.Now().Month()), 10) + "/" + strconv.FormatUint(uint64(time.Now().Day()), 10) + ".json")
		mostRecentReadDates.day = int(time.Now().Day())
	}

}

func checkCalendar(filename string) {
	dt := time.Now()
	fmt.Println("Reading " + filename)

	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of timeInstances
	var todaysArray timeInstances
	unmarshalErr1 := json.Unmarshal(fileContent, &todaysArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}
	for event := 0; event < len(todaysArray); event++ {
		if uint8(todaysArray[event].Hour) == uint8(dt.Hour()) && uint8(todaysArray[event].Minute) == uint8(dt.Minute()) {
			// Send email
			sendEmail(todaysArray[event].Title, "Calendar event \""+todaysArray[event].Title+"\"")
			// Delete
			todaysArray = append(todaysArray[:event], todaysArray[event+1:]...)
			// Convert array back to json
			writeJson, writeJsonErr := json.Marshal(todaysArray)
			if writeJsonErr != nil {
				fmt.Println(writeJsonErr)
			}
			// Write eventArray to schedules.json file
			writeErr := ioutil.WriteFile(filename, writeJson, 0644)
			if writeErr != nil {
				fmt.Println(writeErr)
			}
		}
	}
}

func checkSchedule(filename string) {
	dt := time.Now()
	fmt.Println("Reading " + filename)

	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of timeInstances
	var scheduleArray []Schedule
	unmarshalErr1 := json.Unmarshal(fileContent, &scheduleArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}
	var weekday string = time.Now().Weekday().String()
	for event := 0; event < len(scheduleArray); event++ {
		correctWeekday := false
		switch weekday {
		case "Sunday":
			if scheduleArray[event].Sunday {
				correctWeekday = true
			}
			break
		case "Monday":
			if scheduleArray[event].Monday {
				correctWeekday = true
			}
			break
		case "Tuesday":
			if scheduleArray[event].Tuesday {
				correctWeekday = true
			}
			break
		case "Wednesday":
			if scheduleArray[event].Wednesday {
				correctWeekday = true
			}
			break
		case "Thursday":
			if scheduleArray[event].Thursday {
				correctWeekday = true
			}
			break
		case "Friday":
			if scheduleArray[event].Friday {
				correctWeekday = true
			}
			break
		case "Saturday":
			if scheduleArray[event].Saturday {
				correctWeekday = true
			}
			break
		}

		if correctWeekday && scheduleArray[event].Enabled && uint8(scheduleArray[event].Hour) == uint8(dt.Hour()) && uint8(scheduleArray[event].Minute) == uint8(dt.Minute()) {
			// Send email
			sendEmail(scheduleArray[event].Title, "Schedule event \""+scheduleArray[event].Title+"\"")
		}
	}
}

func timeListener() {
	for true {
		//fmt.Println("Listener going off")
		dt := time.Now()
		var year string = strconv.FormatInt(int64(dt.Year()), 10)
		var month string = strconv.FormatInt(int64(dt.Month()), 10)
		var day string = strconv.FormatInt(int64(dt.Day()), 10)

		calendarReadFile := "database/time/calendar/" + year + "/" + month + "/" + day + ".json"
		annualReadFile := "database/time/calendar/annual/" + month + "/" + day + ".json"
		fmt.Println(annualReadFile)
		fmt.Println(calendarReadFile)

		// Check if we can clean up files or directories
		cleanCalendarHistory()

		// See if annualReadFile is an actual file
		if _, err2 := os.Stat(annualReadFile); !os.IsNotExist(err2) {
			// File exists, can read
			checkCalendar(annualReadFile)
		}

		// See if calendarReadFile is an actual file
		if _, err2 := os.Stat(calendarReadFile); !os.IsNotExist(err2) {
			// File exists, can read
			checkCalendar(annualReadFile)
		}

		// See if calendarReadFile is an actual file
		var scheduleFile string = "database/time/schedules/schedules.json"
		if _, err2 := os.Stat(scheduleFile); !os.IsNotExist(err2) {
			// File exists, can read
			checkSchedule(scheduleFile)
		}

		// Read every 10 seconds
		time.Sleep(time.Millisecond * 10000)
	}
}
