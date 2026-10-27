package config

type PaymentGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OtelTracerConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
