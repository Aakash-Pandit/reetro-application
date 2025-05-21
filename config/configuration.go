package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ReadEnvironmentVariables() error {
	log.Println("Reading all the env variables")

	dir, _ := os.Open("./.envs/.env")
	err := godotenv.Load(dir.Name())
	if err != nil {
		log.Println("err:", err)
		return err
	}

	return nil
}
