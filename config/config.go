package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	MidtransServerKey  string
	MidtransClientKey  string
	MidtransProduction bool
}

func Load() *Config {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	isProduction, _ := strconv.ParseBool(os.Getenv("MIDTRANS_IS_PRODUCTION"))

	if serverKey == "" || clientKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY and MIDTRANS_CLIENT_KEY must be set")
	}

	env := "Sandbox"
	if isProduction {
		env = "Production"
	}
	log.Printf("[Config] Midtrans environment: %s", env)

	return &Config{
		MidtransServerKey:  serverKey,
		MidtransClientKey:  clientKey,
		MidtransProduction: isProduction,
	}
}
