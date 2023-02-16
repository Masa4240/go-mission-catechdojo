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
	defer logger.Sync()
	logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newUser.Name))

	user := *newUser
	// newUser.Name = newName

	if res := s.db.Table("user_list").Create(&user); res.Error != nil {
		logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		return nil, res.Error
	}
	return &user, nil
}

func (s *UserModel) GetUserByID(user *UserLists) (*UserLists, error) {
	userList := UserLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.Int("Requested ID", int(user.ID)))
	if err := s.db.Table("user_list").Find(&userList, "id=?", user.ID).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return &userList, nil
}

func (s *UserModel) GetUserByName(user *UserLists) ([]*UserLists, error) {
	userList := []*UserLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.String("Requested ID", (user.Name)))
	if err := s.db.Table("user_list").Find(&userList, "name=?", user.Name).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return userList, nil
}

func (s *UserModel) UpdateUser(user *UserLists) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Update User model", zap.Time("now", time.Now()))

	// if err := s.db.Table("user_list").Where("id=?", user.ID).Update("name", user.Name).Error; err != nil {
	// 	logger.Info("Fail to update DB", zap.Time("now", time.Now()),
	// 		zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
	// 	err := errors.New("fail update db")
	// 	return err
	// }
	if err := s.db.Table("user_list").Updates(&user).Error; err != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
		return errors.New("fail update db")
		// return err
	}
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
