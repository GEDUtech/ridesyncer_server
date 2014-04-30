package email

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Config      *Config
	ContentType string
	Charset     string
	Subject     string
	To          []string
}

func NewMailer(config *Config) *Mailer {
	return &Mailer{Config: config, To: make([]string, 0)}
}

func (m *Mailer) AddTo(to ...string) *Mailer {
	for _, t := range to {
		m.To = append(m.To, t)
	}
	return m
}

func (m *Mailer) SetSubject(subject string) *Mailer {
	m.Subject = subject
	return m
}

func (m *Mailer) SetContentType(contentType string) *Mailer {
	m.ContentType = contentType
	return m
}

func (m *Mailer) SetCharset(charset string) *Mailer {
	m.Charset = charset
	return m
}

func (m *Mailer) Send(body string) error {
	content := fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: %s; charset=\"%s\"\n\n%s",
		m.Subject, m.ContentType, m.Charset, body)
	return smtp.SendMail(m.Config.Addr, m.Config.Auth, "RideSyncer", m.To, []byte(content))
}
