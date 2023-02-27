package usermodel

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserModel struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserModel(db *gorm.DB, logger *zap.Logger) *UserModel {
	return &UserModel{
		db:     db,
		logger: logger,
	}
}

func (s *UserModel) CreateUser(newUser *UserList) (*UserList, error) {
	s.logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newUser.Name))

	user := *newUser
	tx := s.db.Begin()
	if res := tx.Table("user_list").Create(&user); res.Error != nil {
		s.logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		tx.Rollback()
		return nil, res.Error
	}
	tx.Commit()
	return &user, nil
}

func (s *UserModel) GetUserByID(user *UserList) (*UserList, error) {
	userList := UserList{}

	s.logger.Info("Start to get GetUserByID Model", zap.Time("now", time.Now()), zap.Int("Requested ID", int(user.ID)))
	tx := s.db.Begin()

	if err := tx.Table("user_list").Find(&userList, "id=?", user.ID).Error; err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &userList, nil
}

func (s *UserModel) GetUserByName(user *UserList) ([]*UserList, error) {
	userList := []*UserList{}
	s.logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.String("Requested ID", (user.Name)))
	tx := s.db.Begin()

	if err := tx.Table("user_list").Find(&userList, "name=?", user.Name).Error; err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return userList, nil
}

func (s *UserModel) UpdateUser(user *UserList) error {
	s.logger.Info("Update User model", zap.Time("now", time.Now()))

	tx := s.db.Begin()
	tx.Lock()
	if err := tx.Table("user_list").Where("id=?", user.ID).Updates(&user).Error; err != nil {
		s.logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
		tx.Rollback()
		return errors.New("fail update db")
	}
	tx.Commit()
	return nil
}
