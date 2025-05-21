package common

import (
	"fmt"
	"log"
	"os"

	gomail "gopkg.in/gomail.v2"
)

var (
	EMAIL_ID       = os.Getenv("EMAIL_ID")
	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	EMAIL_HOST     = os.Getenv("EMAIL_HOST")
	EMAIL_PORT     = os.Getenv("EMAIL_PORT")
)

type Email struct {
	EmailFrom, Password, Subject, Body, Host, Address string
	EmailTo                                           []string
}

type EmailInterface interface {
	SendEmailForPasswordReset(emailId, subject, password string) error
}

func (e *Email) SendEmailForPasswordReset(emailId, subject, password string) error {
	e.Subject = subject
	e.Body = fmt.Sprintf("Your New Password is: %s", password)
	e.EmailTo = []string{emailId}

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.EmailFrom)
	msg.SetHeader("To", e.EmailTo...)
	msg.SetHeader("Subject", e.Subject)
	msg.SetBody("text/html", e.Body)

	n := gomail.NewDialer(e.Host, 587, e.EmailFrom, e.Password)
	err := n.DialAndSend(msg)

	if err != nil {
		return err
	}

	log.Println("Email Send Successfully.")
	return nil
}
