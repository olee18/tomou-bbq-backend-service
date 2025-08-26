package database

import (
	"fmt"
	"laotop_final/config"
	"laotop_final/logs"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

var (
	openConnectionDB *gorm.DB
)

var err error

func PostgresConnection() (*gorm.DB, error) {
	myDSN := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Bangkok",
		config.Env("postgres.host"),
		config.Env("postgres.user"),
		config.Env("postgres.password"),
		config.Env("postgres.database"),
		config.Env("postgres.port"),
	)
	fmt.Println("Postgres connecting")
	openConnectionDB, err = gorm.Open(postgres.Open(myDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation("Asia/Bangkok")
			return time.Now().In(ti)
		},
	})
	if err != nil {
		logs.Error(err)
		log.Fatal("ERROR_PING_POSTGRES", err)
		return nil, err
	}

	sqlDB, err := openConnectionDB.DB()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(200)
	fmt.Println("Postgres connected")
	return openConnectionDB, nil
}
func CloseConnectionPostgres(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			logs.Error(err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			logs.Error(err)
			return
		}
	}
}
