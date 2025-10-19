package main

import (
	"problum/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to create application")
		return
	}

	if err = app.Run(); err != nil {
		log.Error().Err(err).Msg("failed to run application")
		return
	}
}
