package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpPort        = "587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		body string,
		to []string,
		cc []string,
		bcc []string,
		attachments []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	body string,
	to []string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.To = []string{"developer.levantung@gmail.com"}
	e.Subject = subject
	e.HTML = []byte(body)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, attachment := range attachments {
		if _, err := e.AttachFile(attachment); err != nil {
			return fmt.Errorf("cannot attach file: %w", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	if err := e.Send(smtpAuthAddress+":"+smtpPort, smtpAuth); err != nil {
		return fmt.Errorf("cannot send email: %w", err)
	}

	return nil
}
