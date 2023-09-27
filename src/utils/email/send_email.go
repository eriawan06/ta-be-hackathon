package email

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"strconv"
)

type TemplateData struct {
	Name        string
	Link        string
	SenderEmail string
	Title       string
	Type        string
}

type TeamInvitationTemplateData struct {
	Title       string
	InvitedName string
	SenderName  string
	DetailLink  string
	AcceptLink  string
}

type TeamRequestTemplateData struct {
	Title           string
	TeamCreatorName string
	SenderName      string
	DetailLink      string
	AcceptLink      string
}

type Request struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Body = buf.String()
	return nil
}

func (r *Request) SendEmail() (bool, error) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_SENDER_EMAIL"))
	mailer.SetHeader("To", r.To...)
	mailer.SetHeader("Subject", r.Subject)
	mailer.SetBody("text/html", r.Body)

	SmtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		SmtpPort,
		os.Getenv("SMTP_SENDER_EMAIL"),
		os.Getenv("SMTP_SENDER_PASSWORD"),
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return false, err
	}

	return true, nil
}
