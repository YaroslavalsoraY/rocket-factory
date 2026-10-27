package config

import "time"

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type SessionConfig interface {
	TTL() time.Duration
}

type IamGRPCConfig interface {
	Address() string
}

type RedisConfig interface {
	Address() string
	MaxIdle() int
	Timeout() time.Duration
	IdleTimeout() time.Duration
}
