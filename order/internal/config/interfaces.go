package config

import "time"

type OrderHttpConfig interface {
	HttpAdress() string
	HeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type InventoryGRPCConfig interface {
	InventoryAdress() string
}

type PaymentGRPCConfig interface {
	PaymentAdress() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
