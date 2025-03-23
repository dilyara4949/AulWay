package service

import (
	"aulway/internal/utils/config"
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body string, smtpConfig config.SMTP) error {
	from := smtpConfig.Username
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
