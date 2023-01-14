package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/handler/router"
	"github.com/Masa4240/go-mission-catechdojo/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	logger, _ := zap.NewDevelopment()
	logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))

	// config values
	const (
		defaultPort = ":8080"
	)

	// set time zone
	time.Local, _ = time.LoadLocation("Asia/Tokyo")

	DBMS := "mysql"
	USER := "root"
	PASS := "root"
	PROTOCOL := "tcp(ca-mission:3306)"
	DBNAME := "ca_mission"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	// set up mysql
	//userDB, err1 := sql.Open(DBMS, CONNECT)
	userDB, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		logger.Info("Fail to connect DB", zap.Time("now", time.Now()), zap.Error(err))
	}
	defer userDB.Close()

	if !userDB.HasTable("user_lists") {
		logger.Info("No target table. Start to create table", zap.Time("now", time.Now()))
		return nil
	}

	// Monster Lists

	if err := service.NewGachaService(userDB).GetChars(); err != nil {
		logger.Info("Fail to get master character list", zap.Time("now", time.Now()))
		return nil
	}

	mux := router.NewRouter(userDB)
	srv := &http.Server{
		Addr:    defaultPort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			logger.Info("Server closed with err", zap.Time("now", time.Now()), zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create Signal monitoring
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt) //Monitoring signals
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		logger.Info("Fail graceful shutdown", zap.Time("now", time.Now()), zap.Error(err))
	}
	return nil
}
