package usermodel

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	// "github.com/jinzhu/gorm"
	// "go.uber.org/zap"
)

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (s *UserModel) CreateUser(ctx context.Context, newUser *UserLists) (*UserLists, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Model", zap.Time("now", time.Now()), zap.String("new name is", newUser.Name))

	user := *newUser
	// newUser.Name = newName

	if res := s.db.Create(&user); res.Error != nil {
		logger.Info("Start Create User Process", zap.Time("now", time.Now()), zap.Error(res.Error))
		return nil, res.Error
	}
	return &user, nil
}

func (s *UserModel) GetUserById(ctx context.Context, user *UserLists) (*UserLists, error) {
	userList := UserLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.Int("Requested ID", int(user.ID)))
	if err := s.db.Table("user_lists").Find(&userList, "id=?", user.ID).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return &userList, nil
}

func (s *UserModel) GetUserByName(ctx context.Context, user *UserLists) ([]*UserLists, error) {
	userList := []*UserLists{}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User Model", zap.Time("now", time.Now()), zap.String("Requested ID", (user.Name)))
	if err := s.db.Table("user_lists").Find(&userList, "name=?", user.Name).Error; err != nil {
		logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	return userList, nil
}

// Modelを受け取る、例えばUserListで受け取って
// GORMならModelを渡すと変更がある部分だけ変更してくれる
func (s *UserModel) UpdateUser(ctx context.Context, user *UserLists) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Update User model", zap.Time("now", time.Now()))

	if err := s.db.Table("user_lists").Where("id=?", user.ID).Update("name", user.Name).Error; err != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()),
			zap.String("name", user.Name), zap.Int("id", int(user.ID)), zap.Error(err))
		err := errors.New("fail update db")
		return err
	}
	return nil
}

// GetUserByName でNameで検索、Service側で読んでNilなら重複なし、Nilじゃなかったら重複あり
// Business LogicはService, Modeは基本的にSQLとやり取りするだけ

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

// なくていい、テーブルの更新、作成はMigrationという仕組みで行う
// func (s *UserModel) TableConfirmation(ctx context.Context) error {
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()
// 	logger.Info("Start TableConfirmation model", zap.Time("now", time.Now()))
// 	if !s.db.HasTable("user_lists") {
// 		logger.Info("No target table. Start to create table")
// 		if res := s.db.Table("user_lists").AutoMigrate(&UserLists{}); res.Error != nil {
// 			logger.Info("Error to create table", zap.Error(res.Error))
// 			return res.Error
// 		}
// 	}
// 	return nil
// }
