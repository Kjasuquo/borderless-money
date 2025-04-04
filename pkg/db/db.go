package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"operation-borderless/internal/domain/model"
	"operation-borderless/pkg/config"
)

func Init(cfg *config.Config) (*gorm.DB, error) {

	var dbURL string

	if cfg.DatabaseUrl == "" {
		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	} else {
		dbURL = cfg.DatabaseUrl
	}

	if dbURL == "" {
		return nil, errors.New("invalid database url")
	}

	log.Println("dbURL: ", dbURL)

	dbConn, err := gorm.Open(
		postgres.Open(dbURL),
		defaultGormConfig(),
	)
	if err != nil {
		log.Println("Init dbConn error: ", err)
		return nil, err
	}

	db, err := dbConn.DB()
	if err != nil {
		log.Println("Init db error: ", err)
		return nil, errors.New("couldn't get DB object")
	}

	err = dbConn.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{}, &model.AuditLog{})
	if err != nil {
		log.Println("error migrating tables: ", err)
		log.Panic(err)
	}

	db.SetMaxOpenConns(60)                  // maximum number of open connection to database
	db.SetMaxIdleConns(40)                  // maximum number of connections in the idle connection pool.
	db.SetConnMaxLifetime(10 * time.Minute) // maximum amount of time a connection may be reused

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return dbConn, nil
}

// defaultGormConfig returns a default GORM configuration
// with default settings suited for this backend project
func defaultGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:      logger.Silent, // log level
				Colorful:      true,          // enable colored output
				SlowThreshold: time.Second * 5,
				// ^ threshold for what should be considered a slow query
			},
		),
		FullSaveAssociations: true,
	}
}
