package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	router "github.com/Masa4240/go-mission-catechdojo/handler"
	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	rankmodel "github.com/Masa4240/go-mission-catechdojo/model/rankratio"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"

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
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			logger.Info("zap err")
			fmt.Println(err)
		}
	}(logger)
	logger.Info("Start", zap.String("key", "value"), zap.Time("now", time.Now()))

	// config values
	const (
		defaultPort = ":8080"
	)

	// set time zone
	time.Local, _ = time.LoadLocation("Asia/Tokyo")

	dbms := "mysql"
	user := "root"
	pass := "xxx"
	protocol := "tcp(ca-mission:3306)"
	dbname := "ca_mission"

	connect := user + ":" + pass + "@" + protocol + "/" + dbname + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	// set up mysql
	// userDB, err1 := sql.Open(DBMS, CONNECT)
	userDB, err := gorm.Open(dbms, connect)
	if err != nil {
		logger.Info("Fail to connect DB", zap.Time("now", time.Now()), zap.Error(err))
		const wait = 3
		for i := 0; i < 5; i++ {
			time.Sleep(wait * time.Second)
			userDB, err = gorm.Open(dbms, connect)
			if err == nil {
				break
			}
		}
	}
	defer userDB.Close()

	// Master Data initialization
	if err := gachaservice.NewGachaService(
		charactermodel.NewCharacterModel(userDB, logger),
		ucmodel.NewUcModel(userDB, logger),
		usermodel.NewUserModel(userDB, logger),
		rankmodel.NewRankModel(userDB, logger),
		logger,
	).InitMasterData(); err != nil {
		logger.Info("Fail to get master data", zap.Time("now", time.Now()), zap.Error(err))
		return err
	}

	mux := router.NewRouter(userDB, logger)
	const serverTimeout = 10
	srv := &http.Server{
		Addr:              defaultPort,
		Handler:           mux,
		ReadHeaderTimeout: serverTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			logger.Info("Server closed with err", zap.Time("now", time.Now()), zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create Signal monitoring
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt) // Monitoring signals
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	const gsTimer = 10
	timer := gsTimer * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timer)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		logger.Info("Fail graceful shutdown", zap.Time("now", time.Now()), zap.Error(err))
	}
	return nil
}
