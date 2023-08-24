package utils

import (
	"crypto/tls"
	"fmt"

	config "github.com/aditansh/go-notes/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) (string, error) {
	config, _ := config.LoadEnvVariables(".")
	m := gomail.NewMessage()

	m.SetHeader("From", config.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.Email, config.EmailPassword)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return "", err
}
