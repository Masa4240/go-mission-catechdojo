package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/form3tech-oss/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserService struct {
	db *gorm.DB
}
type userLists struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	Name string
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, newName string) (*model.UserInfo, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.String("new name is", newName))

	if utf8.RuneCountInString(newName) > 10 {
		logger.Info("Too long name.", zap.String("New name is", newName), zap.Time("now", time.Now()))
		err := errors.New("Invalid name")
		return nil, err
	}

	userList := []userLists{}
	if err := s.db.Table("user_lists").Find(&userList, "name=?", newName).Error; err != nil {
		logger.Info("Fail to add new name to DB", zap.Time("now", time.Now()), zap.String("NewName is", newName), zap.Error(err))
		fmt.Println(err)
		err := errors.New("Fail to read data to DB")
		return nil, err
	}

	if len(userList) != 0 {
		err := errors.New("Duplicated user name")
		return nil, err
	}

	newUser := userLists{}
	newUser.Name = newName
	res := s.db.Create(&newUser)
	if res.Error != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()))
		return nil, res.Error
	}
	if res.RowsAffected != 1 {
		logger.Info("Fail to add new name to DB", zap.Time("now", time.Now()), zap.String("NewName is", newName), zap.Int64("Affected row", res.RowsAffected))
		err := errors.New("Fail to create data to DB")
		return nil, err
	}

	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = newUser.ID
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))

	userinfo := model.UserInfo{
		ID:    int64(newUser.ID),
		Name:  newName,
		Token: tokenString,
	}
	return &userinfo, nil
}

func (s *UserService) GetUser(ctx context.Context, reqID int) (*model.UserGetResponse, error) {
	userList := userLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User name", zap.Time("now", time.Now()), zap.Int("Requested ID", reqID))
	res := s.db.Table("user_lists").Find(&userList, "id=?", reqID)
	if res.Error != nil {
		logger.Info("ID not found", zap.Time("now", time.Now()), zap.Int("req ID", reqID), zap.Error(res.Error))
		err := errors.New("ID Not Found")
		return nil, err
	}
	response := model.UserGetResponse{
		Name: userList.Name,
	}
	return &response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, newName string, reqID int) (*model.UserGetResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	if utf8.RuneCountInString(newName) > 10 {
		logger.Info("Too long name", zap.Time("now", time.Now()), zap.String("New Name", newName), zap.Int("ID", reqID))
		err := errors.New("Invalid name")
		return nil, err
	}
	//Duplication check
	userList := userLists{}
	if err := s.db.Table("user_lists").Find(&userList, "name=?", newName).Error; err == nil {
		logger.Info("User name already exists", zap.Time("now", time.Now()), zap.String("name", newName))
		err := errors.New("Duplicated Name")
		return nil, err
	}
	// //confirm token
	// start the transaction
	tx := s.db.Begin()
	if err := s.db.Table("user_lists").Where("id=?", reqID).Update("name", newName).Error; err != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()), zap.String("name", newName), zap.Int("id", reqID), zap.Error(err))
		err := errors.New("Fail update DB")
		tx.Rollback()
		return nil, err
	}
	if err := s.db.Where("id=?", reqID).Take(&userList).Error; err != nil {
		logger.Info("Fail to confirm new DB", zap.Time("now", time.Now()), zap.String("name", newName), zap.Int("id", reqID))
		err := errors.New("Fail to confirm new DB")
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil
}
