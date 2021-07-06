package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type Email struct {
	D *gomail.Dialer
}

func NewEmail(host string, port int, username string, password string) *Email {
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Email{
		D: d,
	}
}
