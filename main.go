package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/handler/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Start")
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort = ":8080"
		//		defaultDBPath = ".sqlite3/todo.db"
	)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = defaultPort
	// }

	// dbPath := os.Getenv("DB_PATH")
	// if dbPath == "" {
	// 	dbPath = defaultDBPath
	// }

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}
	DBMS := "mysql"
	USER := "root"
	PASS := "root"
	PROTOCOL := "tcp(ca-mission:3306)"
	DBNAME := "ca_mission"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME //+ "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	// set up mysql
	userDB, err1 := sql.Open(DBMS, CONNECT)
	defer userDB.Close()
	err1 = userDB.Ping()

	if err != nil {
		fmt.Println("Fail to connect db")
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		fmt.Println("Try again to connect db")
		userDB, err1 = sql.Open(DBMS, CONNECT)
		err1 = userDB.Ping()
		if err1 != nil {
			fmt.Println("Fail to connect db")
			fmt.Println(err1)
		}

		//return err
	} else {
		fmt.Println("db connection success")
	}
	result, err0 := userDB.Exec("INSERT INTO users(name, token) VALUES('hi1', 'hi1')")
	if err0 != nil {
		fmt.Println("Fail exec")
		fmt.Println(err0)
		fmt.Println(result)
	}
	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(userDB)
	//	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする
	// サーバーをポート8080で起動
	//http.ListenAndServe(defaultPort, mux)
	srv := &http.Server{
		Addr:    defaultPort,
		Handler: mux,
	}

	go func() {
		log.Println("Go Routine")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalln("Server closed with error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create Signal monitoring
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt) //Monitoring signals
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Println("Failed to gracefully shutdown:", err)
	}
	log.Println("Server shutdown")

	return nil
}
