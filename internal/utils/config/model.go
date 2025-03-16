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
	Redis
	SMTP
}

type Redis struct {
	Host     string
	Port     string
	Password string
	Timeout  int
	PoolSize int
	Database int
	Duration time.Duration
}

type SMTP struct {
	Host     string
	Port     string
	Username string
	Password string
}
