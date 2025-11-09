package main

import (
	"problum/internal/worker"

	"github.com/rs/zerolog/log"
)

func main() {
	w, err := worker.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create worker")
		return
	}

	if err = w.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run worker")
		return
	}
}
