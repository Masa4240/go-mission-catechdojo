package usermodel

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserModel struct {
	// db     *gorm.DB
	logger *zap.Logger
}

func NewUserModel(logger *zap.Logger) *UserModel {
	return &UserModel{
		// db:     db,
		logger: logger,
	}
}

func (s *UserModel) CreateUser(db *gorm.DB, newUser *UserList) (*UserList, error) {
	s.logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newUser.Name))

	user := *newUser
	if res := db.Table("user_list").Create(&user); res.Error != nil {
		s.logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		return nil, res.Error
	}
	return &user, nil
}

func (s *UserModel) GetUserByID(db *gorm.DB, user *UserList, lock bool) (*UserList, error) {
	userList := UserList{}
	if lock {
		db.Lock()
	}

	s.logger.Info("Start to get GetUserByID Model", zap.Time("now", time.Now()), zap.Int("Requested ID", int(user.ID)))

	if err := db.Table("user_list").Find(&userList, "id=?", user.ID).Error; err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return &userList, nil
}

func (s *UserModel) GetUserByName(db *gorm.DB, user *UserList) ([]*UserList, error) {
	userList := []*UserList{}
	s.logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.String("Requested ID", (user.Name)))

	if err := db.Table("user_list").Find(&userList, "name=?", user.Name).Error; err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return userList, nil
}

func (s *UserModel) UpdateUser(db *gorm.DB, user *UserList) error {
	s.logger.Info("Update User model", zap.Time("now", time.Now()))

	if err := db.Table("user_list").Where("id=?", user.ID).Updates(&user).Error; err != nil {
		s.logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
		return errors.New("fail update db")
	}
	return nil
}
