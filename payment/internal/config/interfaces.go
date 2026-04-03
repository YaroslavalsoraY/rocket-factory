package config

type PaymentGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
