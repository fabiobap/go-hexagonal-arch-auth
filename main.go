package main

import (
	"github.com/go-hexagonal-arch-auth/app"
	"github.com/go-hexagonal-arch-auth/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	logger.Info("Starting the application")
	app.Start()
}
