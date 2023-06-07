package main

import (
	"fmt"
	"net/smtp"
	"os"
)

/*
func main() {
	sendEmail("The Nich Neuron Can Now Send Emails.", "The Nich Neuron can now send emails.")
}*/

// getenv not working properly right now
func sendEmail(subject string, body string) {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("NeuronEmail"),
		os.Getenv("NeuronPassword"),
		"smtp.gmail.com")

	msg := "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("NeuronEmail"),
		[]string{os.Getenv("PERSONAL_EMAIL")},
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}
