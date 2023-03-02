package userservice

import (
	"errors"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	"github.com/form3tech-oss/jwt-go"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UserService struct {
	db        *gorm.DB
	ucmodel   *ucmodel.UcModel
	usermodel *usermodel.UserModel
	logger    *zap.Logger
}

func NewUserService(db *gorm.DB, ucmodel *ucmodel.UcModel, usermodel *usermodel.UserModel, logger *zap.Logger) *UserService {
	return &UserService{
		db:        db,
		ucmodel:   ucmodel,
		usermodel: usermodel,
		logger:    logger,
	}
}

func (s *UserService) CreateUser(info UserInfo) (*UserInfo, error) {
	tx := s.db.Begin()
	res, err := s.usermodel.CreateUser(tx, &usermodel.UserList{
		Name: info.Profile.Name,
	})
	if err != nil {
		s.logger.Info("Create User Service", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Info("Fail to create token", zap.Time("now", time.Now()))
		return nil, errors.New("fail to create token")
	}
	claims["id"] = res.ID
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))
	info.Profile.Token = tokenString
	return &info, nil
}

func (s *UserService) GetUserByID(info UserInfo) (*UserInfo, error) {
	tx := s.db.Begin()
	res, err := s.usermodel.GetUserByID(tx, &usermodel.UserList{
		Model: gorm.Model{
			ID: uint(info.Profile.ID),
		},
	}, false)
	if err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, errors.New("id not found")
	}
	tx.Commit()
	info.Profile.Name = res.Name
	return &info, nil
}

func (s *UserService) DuplicationCheck(info UserInfo) error {
	tx := s.db.Begin()
	res, err := s.usermodel.GetUserByName(tx, &usermodel.UserList{
		Name: info.Profile.Name,
	})
	if err != nil {
		tx.Rollback()
		s.logger.Info("Get User By Name Error", zap.Time("now", time.Now()), zap.Error(err))
		return errors.New("not found")
	}
	if len(res) != 0 {
		tx.Rollback()
		s.logger.Info("Duplicated", zap.Time("now", time.Now()), zap.Uint("ID", res[0].ID))
		return errors.New("not found")
	}
	tx.Commit()
	return nil
}

func (s *UserService) UpdateUser(info UserInfo) error {
	// transaction start
	tx := s.db.Begin()
	user, err := s.usermodel.GetUserByID(tx, &usermodel.UserList{
		Model: gorm.Model{
			ID: uint(info.Profile.ID),
		},
	}, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	user.Name = info.Profile.Name
	if err := s.usermodel.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		return errors.New("fail to confirm new db")
	}
	tx.Commit()
	return nil
}

func (s *UserService) GetUserCharacterList(info UserInfo) (*UserInfo, error) {
	tx := s.db.Begin()
	req := usermodel.UserList{
		Model: gorm.Model{
			ID: uint(info.Profile.ID),
		},
	}

	list, err := s.ucmodel.GetCharaterList(tx, &req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return convertToUserInfo(&req, list), nil
}
