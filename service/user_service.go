package service

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

type UserServiceMVC struct {
	// db  *gorm.DB
	svc *model.UserModel
}

func NewUserServiceMVC(svc *model.UserModel) *UserServiceMVC {
	return &UserServiceMVC{
		svc: svc,
	}
}

func (s *UserServiceMVC) CreateUserService(ctx context.Context, newName string) (*string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()), zap.String("new name is", newName))

	if err := nameValidation(newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return nil, err
	}
	// tx := s.db.Begin()

	if err := s.svc.TableConfirmation(ctx); err != nil {
		//	tx.Rollback()
		return nil, err
	}

	if err := s.svc.DuplicationCheck(ctx, newName); err != nil {
		//	tx.Rollback()
		return nil, err
	}
	id, err := s.svc.CreateUser(ctx, newName)
	if err != nil {
		//	tx.Rollback()
		return nil, err
	}
	// tx.Commit()
	logger.Info("Finish Create User Service Registration", zap.Time("now", time.Now()), zap.Uint("New ID", *id))

	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = *id
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))

	return &tokenString, nil
}

func (s *UserServiceMVC) GetUserService(ctx context.Context, reqID int) (*string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User name", zap.Time("now", time.Now()), zap.Int("Requested ID", reqID))
	res, err := s.svc.GetUser(ctx, reqID)
	if err != nil {
		logger.Info("ID not found", zap.Time("now", time.Now()), zap.Int("req ID", reqID), zap.Error(err))
		err := errors.New("id not found")
		return nil, err
	}
	return &res, nil
}

func (s *UserServiceMVC) UpdateUserService(ctx context.Context, newName string, reqID int) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	if err := nameValidation(newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return err
	}

	// tx := s.db.Begin()

	// Duplication check
	if err := s.svc.DuplicationCheck(ctx, newName); err != nil {
		//	tx.Rollback()
		return err
	}
	// 	tx := s.db.Lock()をSelectするタイミングで呼ぶ必要がある
	if err := s.svc.UpdateUser(ctx, newName, reqID); err != nil {
		logger.Info("Fail to update new DB", zap.Time("now", time.Now()), zap.String("name", newName), zap.Int("id", reqID))
		err := errors.New("fail to confirm new db")
		//	tx.Rollback()
		return err
	}
	// tx.Commit()
	return nil
}

func nameValidation(newName string) error {
	if len(newName) == 0 {
		err := errors.New("null name")
		return err
	}
	maxLength := 10
	if utf8.RuneCountInString(newName) > maxLength {
		err := errors.New("too long name")
		return err
	}
	return nil
}
