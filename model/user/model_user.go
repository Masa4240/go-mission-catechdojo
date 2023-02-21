package usermodel

import (
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

func (s *UserModel) CreateUser(newUser *UserLists) (*UserLists, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newUser.Name))

	user := *newUser
	tx := s.db.Begin()
	if res := tx.Table("user_list").Create(&user); res.Error != nil {
		logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		tx.Rollback()
		return nil, res.Error
	}
	tx.Commit()
	return &user, nil
}

func (s *UserModel) GetUserByID(user *UserLists) (*UserLists, error) {
	userList := UserLists{}
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.Int("Requested ID", int(user.ID)))
	tx := s.db.Begin()

	if err := tx.Table("user_list").Find(&userList, "id=?", user.ID).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &userList, nil
}

func (s *UserModel) GetUserByName(user *UserLists) ([]*UserLists, error) {
	userList := []*UserLists{}
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.String("Requested ID", (user.Name)))
	tx := s.db.Begin()

	if err := tx.Table("user_list").Find(&userList, "name=?", user.Name).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return userList, nil
}

func (s *UserModel) UpdateUser(user *UserLists) error {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Update User model", zap.Time("now", time.Now()))
	tx := s.db.Begin()
	tx.Lock()
	if err := tx.Table("user_list").Where("id=?", user.ID).Updates(&user).Error; err != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
		tx.Rollback()
		return errors.New("fail update db")
	}
	tx.Commit()
	return nil
}

// func (s *UserModel) DuplicationCheck(ctx context.Context, newName string) error {
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()
// 	logger.Info("Start Duplicationcheck model", zap.Time("now", time.Now()))
// 	var userList []UserLists
// 	err := s.db.Table("user_lists").Find(&userList, "name=?", newName).Error
// 	if err != nil {
// 		logger.Info("Error to find duplication", zap.Time("now", time.Now()), zap.String("name", newName))
// 		err := errors.New("duplicated name")
// 		return err
// 	}
// 	if len(userList) != 0 {
// 		logger.Info("Duplication found", zap.Time("now", time.Now()),
// 			zap.String("name", newName), zap.Int("", int(userList[0].ID)))
// 		err := errors.New("duplicated user name")
// 		return err
// 	}
// 	logger.Info("No Duplication", zap.Time("now", time.Now()))
// 	return nil
// }

// func (s *UserModel) TableConfirmation() error {
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()
// 	logger.Info("Start TableConfirmation model", zap.Time("now", time.Now()))
// 	if !s.db.HasTable("user_list") {
// 		logger.Info("No target table. Start to create table")
// 		if res := s.db.Table("user_list").AutoMigrate(&UserLists{}); res.Error != nil {
// 			logger.Info("Error to create table", zap.Error(res.Error))
// 			return res.Error
// 		}
// 	}
// 	return nil
// }
