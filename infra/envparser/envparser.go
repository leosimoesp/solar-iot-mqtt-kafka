package envparser

import (
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	err = godotenv.Load(currentDir + "/.env")
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}
}

type envparser struct{}

// Parse implements dto.EnvParser.
func (e *envparser) Parse(v interface{}) error {
	return env.Parse(v)
}

func New() dto.EnvParser {
	return &envparser{}
}
