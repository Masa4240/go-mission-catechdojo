package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/form3tech-oss/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, new_name string) (*model.UserInfo, error) {
	println("Start Create User Process")
	fmt.Println(new_name)
	rows, err := s.db.Query("SELECT * FROM userlist")
	if err != nil {
		fmt.Println("Err in query")
		fmt.Println(err)
	}
	if utf8.RuneCountInString(new_name) > 10 {
		fmt.Println("Too long name.")
		err = errors.New("Invalid name")
		return nil, err
	}
	var name string
	var id int
	duplicateFlag := false
	idnumber := 1
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		//fmt.Println(name, tokken)
		//fmt.Println(name == new_name)
		if name == new_name {
			duplicateFlag = true
		}
		if id >= idnumber {
			idnumber = id + 1
		}
	}
	if duplicateFlag {
		fmt.Println("Provided user already exists. Try with another name")
		err = errors.New("Duplicated request")
		return nil, err
	}
	//fmt.Println(duplicateFlag)
	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	//claims["name"] = new_name
	claims["id"] = idnumber

	// Digital Signature
	//	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))
	fmt.Println("Token generated")
	fmt.Println(tokenString)
	userinfo := model.UserInfo{
		ID:    int64(idnumber),
		Name:  new_name,
		Token: tokenString,
	}
	ins, err := s.db.Prepare("INSERT INTO userlist VALUES(?, ?)")
	if err != nil {
		println("Err prepare")
		log.Fatal(err)
		return nil, err
	}
	defer ins.Close()
	//Duplication check

	// Start to update DB
	res, err := ins.Exec(idnumber, new_name)
	println("Exec")
	println(idnumber)
	println(new_name)
	if err != nil {
		println("Fail to execute DB")
		log.Fatal(err)
	}

	// Confirm the result
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		println("Fail to update db correctly")
		log.Fatal(err)
	}
	fmt.Println(lastInsertID)
	println("Finish create user process")
	return &userinfo, nil
}

func (s *UserService) GetUser(ctx context.Context, reqID int) (*model.UserGetReponse, error) {
	rows, err := s.db.Query("SELECT * FROM userlist")
	if err != nil {
		fmt.Println("Err in query")
		fmt.Println(err)
	}
	var name string
	var id int
	returnName := ""
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if id == reqID {
			returnName = name
		}
	}
	if returnName == "" {
		fmt.Println("INVALID Token")
		err = errors.New("INVALID Token")
		return nil, err
	}
	res := model.UserGetReponse{
		Name: returnName,
	}
	return &res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, newname string, reqID int) (*model.UserGetReponse, error) {

	//confirm token
	rows, err := s.db.Query("SELECT * FROM userlist")
	if err != nil {
		fmt.Println("Err in query")
		fmt.Println(err)
	}
	if utf8.RuneCountInString(newname) > 10 {
		fmt.Println("Too long name.")
		err = errors.New("Invalid name")
		return nil, err
	}

	var name string
	var id int
	getUserStatus := false
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if name == newname {
			return nil, errors.New("This name already exists")
		}
		if id == reqID {
			getUserStatus = true
			fmt.Println("Token ID found")
		}
	}
	if !getUserStatus {
		fmt.Println("INVALID Token")
		err = errors.New("INVALID Token")
		return nil, err
	}
	// update DB

	upd, err := s.db.Prepare("UPDATE userlist SET name = ? WHERE userID = ?")
	if err != nil {
		fmt.Println("Fail in db preparation")
		return nil, err
	}
	result, err := upd.Exec(newname, reqID)
	if err != nil {
		return nil, err
	}
	affected_rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected_rows == 0 {
		return nil, errors.New("Update fail")
	}
	return nil, nil
}
