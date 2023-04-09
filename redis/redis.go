package redis

import (
	"balance/config"
	"balance/logging"
	"github.com/go-redis/redis"
	"net"
)

func InitRedis(cfg config.CacheConnConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logging.GetLogger().Fatal(err)
		return nil, err
	}

	return client, nil
}
