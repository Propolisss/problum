package nats

import (
	"fmt"

	"problum/internal/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func New(cfg *config.Nats) (*nats.Conn, error) {
	return nats.Connect(fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port))
}

func NewStream(conn *nats.Conn) (jetstream.JetStream, error) {
	return jetstream.New(conn)
}
