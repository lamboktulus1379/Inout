package persistence

import (
	"fmt"
	"inout/domain/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewRepositories() (*gorm.DB, error) {
	cfg := &Config{}

	cfg.User = os.Getenv("SQL_SERVER_USER")
	cfg.Host = os.Getenv("SQL_SERVER_HOST")
	cfg.Port = os.Getenv("SQL_SERVER_PORT")
	cfg.Password = os.Getenv("SQL_SERVER_PASSWORD")
	cfg.Database = os.Getenv("SQL_SERVER_DATABASE")
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	fmt.Println("DSN: ", dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		),
	})
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	log.Printf("INFO: Connected to DB")
	db.AutoMigrate(&model.User{}, &model.ActivityType{}, &model.UserActivity{})
	return db, nil
}
