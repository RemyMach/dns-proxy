package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() {

	myVar := os.Getenv("REDIS_HOST")
	if myVar == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println(err.Error())
			log.Fatal("Error loading .env file")
		}
	}
}
