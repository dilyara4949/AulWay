package config

import (
	"strconv"
	"time"
)

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

func (s SMTP) PortAsInt() int {
	p, err := strconv.Atoi(s.Port)
	if err != nil {
		return 8080
	}
	return p
}
