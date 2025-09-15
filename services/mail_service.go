package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends a plain-text email
func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	// Setup SMTP server
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	addr := host + ":" + port

	// Message format
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Auth
	auth := smtp.PlainAuth("", from, password, host)

	// Send
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
