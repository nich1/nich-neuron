package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const dataDirectory string = "database/daily/"

func main() {
	year, month, day := time.Now().Date()
	var fileName string = strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + "-" + strconv.Itoa(year)

	if _, err := os.Stat(dataDirectory + fileName + ".json"); err == nil {
		fmt.Println("Day already logged, now update")
	} else {
		fmt.Printf("Day not logged, create file")
	}

}
