package main

import (
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"os"
	"strconv"

	"student-scope-send/app"
	"student-scope-send/cmd"
	"student-scope-send/email"
	"student-scope-send/redis"
	"student-scope-send/util"
)

var envFile = ".env"

func init() {
	if !util.IsFileExists(envFile) {
		panic(fmt.Sprintf("%s 未配置\n", envFile))
	}
	err := godotenv.Load(envFile)
	if err != nil {
		panic(fmt.Sprintf("%s 读取失败: %+v\n", envFile, err))
	}

	var rdb *redis2.Client
	if rdb, err = redis.NewRedis(); err != nil {
		panic(err)
	}

	app.SetRedis(rdb)

	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	fmt.Println(os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT"), os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))
	app.SetEmail(email.NewEmail(
		os.Getenv("EMAIL_HOST"),
		port,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	))
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
