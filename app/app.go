package app

import (
	"github.com/go-redis/redis/v8"
	"student-scope-send/email"
)

var mail *email.Email
var rds *redis.Client

func SetEmail(email2 *email.Email) {
	mail = email2
}

func GetEmail() *email.Email {
	return mail
}

func GetRedis() *redis.Client {
	return rds
}

func SetRedis(r *redis.Client) {
	rds = r
}
