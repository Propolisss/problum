package main

import (
	"problum/internal/app"

	"github.com/rs/zerolog/log"
)

// @title Problum API
// @version 1.0
// @description Problum backend API.
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app, err := app.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create application")
		return
	}

	if err = app.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run application")
		return
	}
}
