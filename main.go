package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Mail struct {
	host     string
	port     string
	from     string
	password string
	username string
}

func NewMail() *Mail {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	username := os.Getenv("SMTP_USERNAME")
	return &Mail{
		host:     host,
		port:     port,
		from:     from,
		password: password,
		username: username,
	}
}

// Send sends an email
func (m *Mail) Send(toEmails []string, subject string, body string) error {
	address := fmt.Sprintf("%s:%s", m.host, m.port)
	message := []byte(subject + body)
	auth := smtp.CRAMMD5Auth(m.username, m.password)
	err := smtp.SendMail(address, auth, m.from, toEmails, message)
	return err
}

func main() {
	mail := NewMail()
	err := mail.Send(
		[]string{"to@example.com"},
		"Test Email",
		"This is a test email for running integration tests using Docker",
	)
	if err != nil {
		panic(err)
	}
}
