package nats

import (
	"balance/config"
	"balance/logging"
	"github.com/nats-io/nats.go"
)

func InitNats(cfg config.BrokerConnConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		logging.GetLogger().Fatal(err)
		return nil, err
	}

	return nc, nil
}
