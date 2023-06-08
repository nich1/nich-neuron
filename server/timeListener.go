package main

import (
	"fmt"
	"time"
)

func timeListener() {
	dt := time.Now()
	for true {
		time.Sleep(time.Millisecond * 3000)
		fmt.Println(dt.String())
	}
}
