package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnvVar(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getIntEnvVar(key string) int64 {

	coonvertVal, convErr := strconv.ParseInt(getEnvVar(key), 10, 64)
	if convErr != nil {
		fmt.Println(convErr)
		os.Exit(1)
	}

	return coonvertVal
}
