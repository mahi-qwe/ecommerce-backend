package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends a plain-text email
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || host == "" || port == "" {
		return fmt.Errorf("smtp config missing in .env")
	}

	addr := host + ":" + port

	// Message format
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	// Auth
	auth := smtp.PlainAuth("", from, password, host)

	// Send mail
	if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
