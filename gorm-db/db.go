package gorm_db

import (
	"balance/config"
	"balance/logging"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection(cfg config.DatabaseConnConfig) (*gorm.DB, error) {
	log := logging.GetLogger()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Dushanbe",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname, cfg.Sslmode)

	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Errorf("%s GoPostgresConnection -> Open error", err.Error())
		return nil, err
	}

	log.Infof("Connection success-->> host:port -->> %s:%s", cfg.Host, cfg.Port)

	return conn, nil
}
