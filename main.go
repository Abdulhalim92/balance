package main

import (
	"balance/config"
	gorm_db "balance/gorm-db"
	"balance/internal/handler"
	"balance/internal/repository"
	"balance/internal/service"
	"balance/logging"
	"balance/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
)

func main() {
	log := logging.GetLogger()

	router := gin.Default()

	cfg := config.GetConfig()
	fmt.Println(cfg) // TODO
	redisClient, err := redis.InitRedis(cfg.CacheConn)
	if err != nil {
		log.Errorf("failed to connection to redis due error: %v", err)
		return
	}

	connection, err := gorm_db.GetDBConnection(cfg.DatabaseConn)
	if err != nil {
		log.Errorf("failed to connection to database due error: %v", err)
		return
	}

	newRepository := repository.NewRepository(connection, log)

	newService := service.NewService(newRepository, log, redisClient)

	newHandler := handler.NewHandler(router, log, newService)
	newHandler.Init()

	address := net.JoinHostPort(cfg.Listen.BindIP, cfg.Listen.Port)

	log.Fatal(router.Run(address))
}
