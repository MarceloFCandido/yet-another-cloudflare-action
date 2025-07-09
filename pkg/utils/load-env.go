package utils

import (
	"github.com/joho/godotenv"
)

var LoadEnv = loadEnv

func loadEnv() error {
	err := godotenv.Load()

	return err
}
