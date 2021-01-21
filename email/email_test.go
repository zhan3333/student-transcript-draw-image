package email

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
	"testing"
)

func TestSend(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.Nil(t, err)
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	assert.Nil(t, err)
	user := os.Getenv("EMAIL_USER")
	from := os.Getenv("EMAIL_FROM")
	mail := NewEmail(
		os.Getenv("EMAIL_HOST"),
		port,
		user,
		os.Getenv("EMAIL_PASSWORD"),
	)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", "390961827@qq.com")
	m.SetHeader("Subject", "成绩单")
	m.Attach("testdata/0001-new.jpg")
	if err := mail.D.DialAndSend(m); err != nil {
		assert.Nil(t, err)
	}
}
