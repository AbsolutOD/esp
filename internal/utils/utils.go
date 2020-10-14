package utils

import (
	"log"
	"os"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func  GetCwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}