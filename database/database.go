package database

import (
	"jobsheet-go/constant"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() {
	Db = dbInit()
}
func dbInit() *gorm.DB {
	// gormのログ設定
	now := time.Now()
	f, err := os.OpenFile("log/gorm-"+now.Format(time.DateOnly)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	newLogger := logger.New(
		log.New(f, "", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)

	// データベースオープン
	db, err := gorm.Open(mysql.Open(constant.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	// コネクションプールの設定
	sqldb, err := db.DB()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Hour)

	return db
}
