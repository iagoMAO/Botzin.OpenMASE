package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MASE_PORT       string
	HB_PORT         string
	BUDDY_PORT      string
	SERVERLIST_PORT string
	XTEA_KEY        string
}

func GetConfig() Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to load .env file.")
		return Config{}
	}

	return Config{
		MASE_PORT:       os.Getenv("MASE_PORT"),
		HB_PORT:         os.Getenv("HB_PORT"),
		BUDDY_PORT:      os.Getenv("BUDDY_PORT"),
		SERVERLIST_PORT: os.Getenv("SERVERLIST_PORT"),
		XTEA_KEY:        os.Getenv("XTEA_KEY"),
	}
}
