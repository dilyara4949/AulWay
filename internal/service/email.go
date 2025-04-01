package service

import (
	"aulway/internal/domain"
	"aulway/internal/utils/config"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(to, subject, body string, smtpConfig config.SMTP) error {
	from := smtpConfig.Username
	msg := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + body)
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendEmailWithQR(to, subject string, tickets []domain.Ticket, smtpConfig config.SMTP, body string) error {
	e := email.NewEmail()
	e.From = smtpConfig.Username
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	for i, t := range tickets {
		if t.QRCode != "" {
			qrBytes, err := base64.StdEncoding.DecodeString(t.QRCode)
			if err != nil {
				return fmt.Errorf("failed to decode QR: %w", err)
			}
			cid := fmt.Sprintf("qr%d.png", i+1)
			_, err = e.Attach(bytes.NewReader(qrBytes), cid, "image/png")
			if err != nil {
				return fmt.Errorf("failed to attach QR image: %w", err)
			}
			// CID is automatically assigned by filename, so use `cid:qr1.png` in the body
		}
	}

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	addr := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)
	if err := e.Send(addr, auth); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
