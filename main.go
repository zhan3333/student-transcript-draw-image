package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"student-scope-send/controller"
	"student-scope-send/email"
	"student-scope-send/redis"
	"student-scope-send/transcript"
	"student-scope-send/util"
)

var envFile = ".env"
var mail *email.Email
var sendEmailsCache = "send_email.json"
var sendEmails []string

func init() {
	if !util.IsFileExists(envFile) {
		panic(fmt.Sprintf("%s 未配置\n", envFile))
	}
	err := godotenv.Load(envFile)
	if err != nil {
		panic(fmt.Sprintf("%s 读取失败: %+v\n", envFile, err))
	}

	if err = redis.InitRedis(); err != nil {
		panic(err)
	}

	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if email.InitEmail(
		os.Getenv("EMAIL_HOST"),
		port,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	) != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "x-requested-with")
	r.Use(cors.New(config))
	r.GET("api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("api/upload", controller.UploadTranscript)
	r.GET("api/query", controller.DownloadTranscriptImg)
	r.Static("api/export", "files/export")
	if err := r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}

func sendEmail(t transcript.Transcript, transcriptFile string) error {
	if t.Email == "" {
		return fmt.Errorf("%s 未配置邮箱", t.Name)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", t.Email)
	m.SetHeader("Subject", fmt.Sprintf("%s成绩单", t.Name))
	m.Attach(transcriptFile)
	return mail.D.DialAndSend(m)
}

func isSendEmail(name string) bool {
	for _, cacheName := range sendEmails {
		if cacheName == name {
			return true
		}
	}
	return false
}

func setSendEmail(name string) {
	sendEmails = append(sendEmails, name)
	b, _ := json.Marshal(sendEmails)
	_ = ioutil.WriteFile(sendEmailsCache, b, os.ModePerm)
}
