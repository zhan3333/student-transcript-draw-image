package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"strconv"
	"student-scope-send/email"
	"student-scope-send/read"
	"student-scope-send/transcript"
	"time"
)

var readExcelPath = "./期末成绩单.xlsx"
var templateFilePath = "./0001.jpg"
var fontFilePath = "./fonts/MSYH.TTC"
var outDir = "out"
var envFile = ".env"
var mail *email.Email
var sendEmailsCache = "send_email.json"
var sendEmailsFile *os.File
var sendEmails []string

func init() {
	if !isFileExists(envFile) {
		panic(fmt.Sprintf("%s 未配置\n", envFile))
	}
	err := godotenv.Load(envFile)
	if err != nil {
		panic(fmt.Sprintf("%s 读取失败: %+v\n", envFile, err))
	}
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	mail = email.NewEmail(
		os.Getenv("EMAIL_HOST"),
		port,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	)
	var f *os.File
	if isFileExists("send_email.json") {
		f, err = os.OpenFile(sendEmailsCache, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("打开缓存文件失败: %+v\n", err))
		}
		text, err := ioutil.ReadAll(f)
		if err != nil {
			panic(fmt.Sprintf("读取缓存失败: %+v\n", err))
		}
		err = json.Unmarshal(text, &sendEmails)
		if err != nil {
			panic(fmt.Sprintf("缓存格式不正确: %+v\n", err))
		}
	} else {
		f, err = os.Create(sendEmailsCache)
		if err != nil {
			panic(fmt.Sprintf("创建缓存文件失败: %+v\n", err))
		}
	}
	sendEmailsFile = f
}

func main() {
	var err error
	var failedSendEmailNames []string

	ts, err := read.Read(readExcelPath)
	if err != nil {
		fmt.Printf("读取成绩失败: %+v\n", err)
		return
	}
	fmt.Printf("读取到%d条成绩单数据\n\n\n", len(*ts))
	//fmt.Printf("读取到的成绩为: %+v\n", *ts)
	for i, t := range *ts {
		fmt.Printf("开始处理第 %d 条数据: %s\n", i+1, t.Name)
		outFile := fmt.Sprintf("%s/%s.jpg", outDir, t.Name)
		d := transcript.NewDrawTranscript(
			templateFilePath,
			outFile,
			fontFilePath,
			t,
		)
		if isFileExists(outFile) {
			fmt.Printf("成绩单 %s 已经存在, 如需重新生成, 请先删除对应成绩单文件\n", outFile)
		} else {
			// 生成成绩单
			err = d.ReadTemplate()
			if err != nil {
				fmt.Printf("读取第 %d 个学生 %s 时读取模板失败: %+v\n", i+1, t.Name, err)
				continue
			}
			err = d.Draw()
			if err != nil {
				fmt.Printf("绘制第 %d 个学生 %s 成绩单失败: %+v\n", i+1, t.Name, err)
				continue
			}
			err = d.Save()
			if err != nil {
				fmt.Printf("保存第 %d 个学生 %s 成绩单失败: %+v\n", i+1, t.Name, err)
				continue
			}
		}

		fmt.Println("开始发送邮件")
		if isSendEmail(t.Name) {
			fmt.Println("缓存发现已经发送过邮件, 跳过发送")
		} else {
			time.Sleep(5 * time.Second)
			err = sendEmail(t, outFile)
			if err != nil {
				fmt.Printf("邮件发送失败: %+v\n", err)
				failedSendEmailNames = append(failedSendEmailNames, t.Name)
			} else {
				setSendEmail(t.Name)
				fmt.Println("邮件发送完毕")
			}
		}

		fmt.Printf("第 %d 条数据: %s 处理完成\n", i+1, t.Name)
		fmt.Println("---------------------------------")
	}
	fmt.Printf("发送邮件失败的学生名单: %v\n", failedSendEmailNames)
	fmt.Println("程序执行完毕")
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
