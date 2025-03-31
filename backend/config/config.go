package config

import (
	"os"
)

type Config struct {
	Port          string
	DatabasePath  string
	JWTSecret     string
	Blockchain    BlockchainConfig
}

type BlockchainConfig struct {
	ProviderURL string
	ContractAddr string
}

func NewConfig() (*Config, error) {
	return &Config{
		Port:          os.Getenv("PORT"),
		DatabasePath:  "db/insights.db",
		JWTSecret:     os.Getenv("JWT_SECRET"),
		Blockchain: BlockchainConfig{
			ProviderURL: os.Getenv("BLOCKCHAIN_PROVIDER"),
			ContractAddr: os.Getenv("CONTRACT_ADDRESS"),
		},
	}, nil
}
