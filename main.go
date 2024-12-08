package main

import (
	"jobsheet-go/database"
	"jobsheet-go/logger"
	"jobsheet-go/util/file"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// ログフォルダがない場合は作成
	if !file.IsDir("log") {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			log.Print(err)
		}
	}

	// システムログ設定
	logger.LogInit()

	// データベース接続
	database.Init()

	// echoを起動
	e := echo.New()

	// ルーティング設定
	SetRoute(e)

	// echoのログ取得
	now := time.Now()
	echo_log, err := os.OpenFile("log/echo-"+now.Format(time.DateOnly)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: echo_log,
	}))

	// パニックが発生してもサーバを継続して稼働させる
	e.Use(middleware.Recover())

	// CORS設定
	e.Use(middleware.CORS())

	// サーバスタート
	slog.Info("Server Start")
	err = e.Start(":8081")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
}
