package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	BlockfrostProjectID string
	WalletMnemonic      string
}

func Load() (AppConfig, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return AppConfig{}, err
	}

	cfg := AppConfig{
		BlockfrostProjectID: os.Getenv("BLOCKFROST_PROJECT_ID"),
		WalletMnemonic:      os.Getenv("WALLET_MNEMONIC"),
	}

	if cfg.BlockfrostProjectID == "" {
		return AppConfig{}, fmt.Errorf("missing required env var: BLOCKFROST_PROJECT_ID")
	}
	if cfg.WalletMnemonic == "" {
		return AppConfig{}, fmt.Errorf("missing required env var: WALLET_MNEMONIC")
	}

	return cfg, nil
}
