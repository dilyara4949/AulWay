package config

import "time"

type Postgres struct {
	Host           string
	Port           string
	User           string
	Password       string
	DB             string
	Timeout        int
	MaxConnections int
}

type Config struct {
	Port              string
	Address           string
	JWTTokenSecret    string
	AccessTokenExpire int
	HeaderTimeout     time.Duration
	Postgres
}
