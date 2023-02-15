package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const dataDirectory string = "database/daily/"

func logFile(payload string) {
	year, month, day := time.Now().Date()
	var fileName string = strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + "-" + strconv.Itoa(year)

	if _, err := os.Stat(dataDirectory + fileName + ".json"); err == nil {
		fmt.Println("Day already logged, now update")
		os.Remove(dataDirectory + fileName + ".json")

	} else {
		fmt.Printf("Day not logged, create file")
	}
	file, err := os.Create(dataDirectory + fileName + ".json")
	if err != nil {
		fmt.Printf("Error in creating file")
	}

	file.WriteString(payload)
	return

}
