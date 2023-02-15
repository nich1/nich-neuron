package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func sendEmail() {
	from := os.Getenv("MAIL")
	password := os.Getenv("PASSWD")

	toList := []string{"example@gmail.com"}

	host := "smtp.gmail.com"

	port := "587"

	msg := "Test message"

	body := []byte(msg)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, toList, body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully sent mail to all user in toList")
}
