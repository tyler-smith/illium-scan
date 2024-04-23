package config

import "github.com/tyler-smith/env"

const (
	envPrefix            = "IEXP_"
	envGRPCServerAddrKey = envPrefix + "GRPC_SERVER_ADDR"
	envDBDSNKey          = envPrefix + "DB_DSN"
	envDBDriverKey       = envPrefix + "DB_DRIVER"

	defaultEnvGRPCServerAddr = "0.0.0.0:5002"
	defaultEnvDBDSN          = "file:dev.db"
	defaultEnvDBDriver       = "sqlite3"
)

type Config struct {
	DB
	Indexer
}

type Indexer struct {
	GRPCServerAddr string
}

type DB struct {
	DSN    string
	Driver string
}

func NewFromEnv() Config {
	return Config{
		DB: DB{
			DSN:    env.GetString(envDBDSNKey, defaultEnvDBDSN),
			Driver: env.GetString(envDBDriverKey, defaultEnvDBDriver),
		},
		Indexer: Indexer{
			GRPCServerAddr: env.GetString(envGRPCServerAddrKey, defaultEnvGRPCServerAddr),
		},
	}
}
