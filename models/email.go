package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "apasdauren70@gmail.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}

	return &es
}

type EmailService struct {
	// used as the default sender when one is not provided for an email.
	// Also used in func where email is a predetermined, like the forgotten password email.
	DefaultSender string

	// unexported fields
	dialer *mail.Dialer
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	es.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.PlainText != "" && email.HTML != "":
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		msg.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)

	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		PlainText: "To reset your password, please visit the following link: " + resetURL,
		HTML: `<p>To reset your password, please visit the following link: <a
		href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return nil
}
