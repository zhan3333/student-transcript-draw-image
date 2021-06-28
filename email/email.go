package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

var EMAIL *Email

type Email struct {
	D *gomail.Dialer
}

func InitEmail(host string, port int, username string, password string) error {
	EMAIL = NewEmail(host, port, username, password)
	return nil
}

func NewEmail(host string, port int, username string, password string) *Email {
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Email{
		D: d,
	}
}
