package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func checkEvent(filename string) {
	dt := time.Now()
	fmt.Println("Reading " + filename)

	// Read in database file for calendar
	fileContent, fileErr := ioutil.ReadFile(calendarFile)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// Convert file into an array of timeInstances
	var todaysArray []timeInstance
	unmarshalErr1 := json.Unmarshal([]byte(fileContent), &todaysArray)
	if unmarshalErr1 != nil {
		fmt.Println(unmarshalErr1)
		return
	}
	for event := 0; event < len(todaysArray); event++ {
		if todaysArray[event].hour == uint8(dt.Hour()) && todaysArray[event].minute == uint8(dt.Minute()) {
			sendEmail(todaysArray[event].title, "Calendar event \""+todaysArray[event].title+"\"")
		}
	}
}

func timeListener() {
	for true {
		fmt.Println("Listener going off")
		dt := time.Now()
		calendarReadFile := "database/time/calendar/" + strconv.FormatInt(int64(dt.Year()), 10) + "/" + strconv.FormatInt(int64(dt.Month()), 10) + "/" + strconv.FormatInt(int64(dt.Day()), 10) + ".json"
		annualReadFile := "database/time/calendar/" + "annual/" + strconv.FormatInt(int64(dt.Month()), 10) + "/" + strconv.FormatInt(int64(dt.Day()), 10) + ".json"

		// See if calendarReadFile is an actual file
		if _, dirErr := os.Stat(calendarReadFile); dirErr != nil {
			if !os.IsNotExist(dirErr) {
				// Dir exists, can read
				checkEvent(calendarReadFile)
				if dirErr != nil {
					fmt.Println(dirErr)
					return
				}
			}
		}
		// See if calendarReadFile is an actual file
		if _, dirErr2 := os.Stat(annualReadFile); dirErr2 != nil {
			if !os.IsNotExist(dirErr2) {
				// Dir exists, can read
				checkEvent(annualReadFile)
				if dirErr2 != nil {
					fmt.Println(dirErr2)
					return
				}
			}
		}

		// Read every 10 seconds
		time.Sleep(time.Millisecond * 10000)
	}
}
