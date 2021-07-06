package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
	"student-scope-send/app"
	"student-scope-send/controller"
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
	app.SetEmail(email.NewEmail(
		os.Getenv("EMAIL_HOST"),
		port,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	))
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
	r.POST("api/upload", controller.Upload)
	r.GET("api/query", controller.Query)
	r.GET("api/send", controller.Send)
	r.Static("api/export", "files/export")
	if err := r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}
