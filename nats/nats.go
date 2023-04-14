package nats

import (
	"balance/config"
	"balance/logging"
	"github.com/nats-io/nats.go"
	"time"
)

// InitNats подключение к брокеру Nats
func InitNats(cfg config.BrokerConnConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		logging.GetLogger().Fatal(err)
		return nil, err
	}

	return nc, nil
}

// SendMessage отправка сообщений
func SendMessage(nc *nats.Conn, message string) error {
	err := nc.Publish("messages", []byte(message))
	if err != nil {
		logging.GetLogger().Error(err)
		return err
	}
	logging.GetLogger().Printf("Sent message: %s", message)
	return nil
}

// ReceiveMessage получение сообщения
func ReceiveMessage(nc *nats.Conn) (*nats.Msg, error) {
	msg, err := nc.Request("messages", nil, 1000*time.Millisecond)
	if err != nil {
		logging.GetLogger().Error(err)
		return nil, err
	}

	return msg, nil
}

// SubscribeToMessages подписка на сообщения
func SubscribeToMessages(nc *nats.Conn) error {
	sub, err := nc.SubscribeSync("messages")
	if err != nil {
		logging.GetLogger().Error(err)
		return err
	}
	defer sub.Unsubscribe()

	// Ожидание сообщения
	msg, err := sub.NextMsg(1000 * time.Millisecond)
	if err != nil {
		logging.GetLogger().Error(err)
		return err
	}

	logging.GetLogger().Printf("Received message: %s", string(msg.Data))

	return nil
}
