package email

import (
	"errors"
	"net/smtp"
)

type Config struct {
	Auth smtp.Auth
	Addr string
}

func NewConfig(username, password, host, port string) (*Config, error) {
	if username == "" || password == "" || host == "" || port == "" {
		return nil, errors.New("Invalid email configuration")
	}
	return &Config{smtp.PlainAuth("", username, password, host), host + ":" + port}, nil
}
