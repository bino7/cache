package cache

import "time"

type Config struct {
	Host        string
	Port        int
	Password    string
	DB          int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}
