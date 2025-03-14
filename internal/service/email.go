package service

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func SendEmail(to, subject, body string) error {
	smtpConfig := SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "your-email@example.com",
		Password: "your-email-password",
	}

	from := smtpConfig.Username
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
