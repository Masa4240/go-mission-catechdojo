package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/form3tech-oss/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type UserService struct {
	db *sql.DB
}

func NewTODOService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name string) (*model.UserInfo, error) {
	println("Start Create User Process")
	// const (
	// 	insert  = `INSERT INTO users(id, name) VALUES(?, ?)`
	// 	confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	// )

	fmt.Println(name)
	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "taro"

	// Digital Signature
	//	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))
	fmt.Println("Token generated")
	fmt.Println(tokenString)
	ins, err := s.db.Prepare("INSERT INTO users VALUES(?, ?)")
	if err != nil {
		println("Err prepare")
		log.Fatal(err)
		return nil, err
	}
	defer ins.Close()
	// SQLの実行
	res, err := ins.Exec(name, "")
	if err != nil {
		println("Err 1")
		log.Fatal(err)
	}

	// 結果の取得
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		println("Err 2")
		log.Fatal(err)
	}
	fmt.Println(lastInsertID)
	println("Finish create user process")
	return nil, nil
}
