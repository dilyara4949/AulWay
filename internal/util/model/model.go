package model

import "time"

type Config struct {
	JWTTokenSecret    string
	AccessTokenExpire int
	RestPort          string
	Address           string
	HeaderTimeout     time.Duration
	Postgres
}

type Postgres struct {
	Host           string
	Port           string
	User           string
	Password       string
	DB             string
	Timeout        int
	MaxConnections int
}
