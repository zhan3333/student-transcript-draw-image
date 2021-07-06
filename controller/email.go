package controller

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"student-scope-send/app"
	"time"
)

func SendEmail(email string, name string, transcriptFile string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("%s成绩单", name))
	m.Attach(transcriptFile)
	return app.GetEmail().D.DialAndSend(m)
}

func IsHasSend(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return app.GetRedis().SIsMember(ctx, "send:emails", email).Result()
}

func SetHasSendEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return app.GetRedis().SAdd(ctx, "send:emails", email).Err()
}
