package model

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (s *UserModel) CreateUser(ctx context.Context, newName string) (*uint, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newName))

	newUser := UserLists{}
	newUser.Name = newName

	if res := s.db.Create(&newUser); res.Error != nil {
		logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		return nil, res.Error
	}
	return &newUser.ID, nil
}

func (s *UserModel) GetUser(ctx context.Context, reqID int) (string, error) {
	userList := UserLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.Int("Requested ID", reqID))
	if err := s.db.Table("user_lists").Find(&userList, "id=?", reqID).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return "", err
	}
	return userList.Name, nil
}

func (s *UserModel) UpdateUser(ctx context.Context, newName string, reqID int) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Update User model", zap.Time("now", time.Now()))

	if err := s.db.Table("user_lists").Where("id=?", reqID).Update("name", newName).Error; err != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", newName), zap.Int("id", reqID), zap.Error(err))
		err := errors.New("fail update db")
		return err
	}
	return nil
}

func (s *UserModel) DuplicationCheck(ctx context.Context, newName string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Duplicationcheck model", zap.Time("now", time.Now()))
	var userList []UserLists
	err := s.db.Table("user_lists").Find(&userList, "name=?", newName).Error
	if err != nil {
		logger.Info("Error to find duplication", zap.Time("now", time.Now()), zap.String("name", newName))
		err := errors.New("duplicated name")
		return err
	}
	if len(userList) != 0 {
		logger.Info("Duplication found", zap.Time("now", time.Now()),
			zap.String("name", newName), zap.Int("", int(userList[0].ID)))
		err := errors.New("duplicated user name")
		return err
	}
	logger.Info("No Duplication", zap.Time("now", time.Now()))
	return nil
}

func (s *UserModel) TableConfirmation(ctx context.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start TableConfirmation model", zap.Time("now", time.Now()))
	if !s.db.HasTable("user_lists") {
		logger.Info("No target table. Start to create table")
		if res := s.db.Table("session_events").AutoMigrate(&UserLists{}); res.Error != nil {
			logger.Info("Error to create table", zap.Error(res.Error))
			return res.Error
		}
	}
	return nil
}
