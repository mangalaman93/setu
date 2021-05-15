package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func email(slots int, err error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Println("no api key found, not emailing")
		return
	}

	m := mail.NewV3Mail()
	p := mail.NewPersonalization()

	fromEmail := os.Getenv("EMAIL_FROM")
	m.SetFrom(mail.NewEmail(fromEmail, fromEmail))

	toEmails := os.Getenv("EMAIL_TO")
	toEmailList := strings.Split(toEmails, ",")
	for _, to := range toEmailList {
		p.AddTos(mail.NewEmail(to, to))
	}

	if err != nil {
		content := mail.NewContent("text/plain", fmt.Sprintf("error occurred: %v", err))
		m.AddContent(content)
	} else {
		content := mail.NewContent("text/plain", fmt.Sprintf("found available slots: %v", slots))
		m.AddContent(content)
	}

	p.Subject = "Setu Updates"
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	if _, err := sendgrid.API(request); err != nil {
		log.Printf("[ERROR] error sending email: %v", err)
	}
	log.Println("email sent")
}
