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
	"student-scope-send/email"
	"student-scope-send/transcript"
)

var templateFilePath = "./0002.jpg"
var fontFilePath = "./fonts/MSYH.TTC"
var envFile = ".env"
var mail *email.Email
var sendEmailsCache = "send_email.json"
var sendEmails []string

func init() {
	if !isFileExists(envFile) {
		panic(fmt.Sprintf("%s 未配置\n", envFile))
	}
	err := godotenv.Load(envFile)
	if err != nil {
		panic(fmt.Sprintf("%s 读取失败: %+v\n", envFile, err))
	}

	if err = InitRedis(); err != nil {
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
	r.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("upload", UploadTranscript)
	r.GET("query", DownloadTranscriptImg)
	r.Static("export", "files/export")
	if err := r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}

func isFileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
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
