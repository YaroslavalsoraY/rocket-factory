package config

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
