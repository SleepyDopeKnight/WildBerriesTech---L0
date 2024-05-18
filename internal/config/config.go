package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Load(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Println(err)
	}
}
