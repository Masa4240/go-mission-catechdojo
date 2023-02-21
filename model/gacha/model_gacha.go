package gachamodel

import (
	"errors"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type GachaModel struct {
	db *gorm.DB
}

func NewGachaModel(db *gorm.DB) *GachaModel {
	return &GachaModel{
		db: db,
	}
}

// func (s *GachaModel) GachaTableCheck(ctx context.Context) error {
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()
// 	logger.Info("Start Create User Model", zap.Time("now", time.Now()))
// 	tableName := "characterlists_users"
// 	if !s.db.HasTable(tableName) {
// 		logger.Info("No target table. Start to create table",
// 			zap.Time("now", time.Now()), zap.String("table name", tableName))
// 		if res := s.db.Table(tableName).AutoMigrate(&UserCharacterList{}); res.Error != nil {
// 			logger.Info("Error to create table", zap.Time("now", time.Now()), zap.Error(res.Error))
// 			return res.Error
// 		}
// 		logger.Info("Table creation done", zap.Time("now", time.Now()))
// 	}
// 	return nil
// }

func (s *GachaModel) GetCharaterList(user *usermodel.UserLists) ([]*UserCharacterList, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	var res []*UserCharacterList

	tableName := "user_characterlist"
	// charList := []*CharacterLists{}
	tx := s.db.Begin()

	if err := tx.Table(tableName).Where("`user_id` = ?", user.ID).Find(&res).Error; err != nil {
		logger.Info("Fail to get char list from db", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return res, nil
}

func (s *GachaModel) RegisterCharacters(characters []*UserCharacterList) error {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	tx := s.db.Begin()

	for i := 0; i < len(characters); i++ {
		res := tx.Table("user_characterlist").Create(characters[i])
		if res.Error != nil {
			logger.Info("Failed to register new char to user db", zap.Time("now", time.Now()), zap.Error(res.Error))
			tx.Rollback()
			return res.Error
		}
	}
	tx.Commit()
	return nil
}

func (s *GachaModel) AddNewCharacter(character *CharacterLists) error {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	tx := s.db.Begin()

	res := tx.Table("formal_character_list").Create(character)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected != 1 {
		logger.Info("Fail to add new char to DB", zap.Time("now", time.Now()), zap.Int64("Affected row", res.RowsAffected))
		tx.Rollback()
		err := errors.New("fail to create data to db")
		return err
	}
	tx.Commit()
	return nil
}

func (s *GachaModel) GetForamalCharacterList() ([]*CharacterLists, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	characters := []*CharacterLists{}
	tx := s.db.Begin()

	if err := tx.Table("formal_character_list").Find(&characters).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return characters, nil
}

func (s *GachaModel) GetRankRatio() ([]*RankRatio, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Get Rarelity Process", zap.Time("now", time.Now()))
	rankRatio := []*RankRatio{}
	tx := s.db.Begin()
	if err := tx.Table("rankratio").Find(&rankRatio).Error; err != nil {
		logger.Info("Error to get Rank ratio", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return rankRatio, nil
}

///////////////////////////////////////////////////
func (s *GachaModel) CharacterTableCheck() error {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Create User Model", zap.Time("now", time.Now()))
	tableName := "formal_character_list"
	if !s.db.HasTable(tableName) {
		logger.Info("No target table. Start to create table",
			zap.Time("now", time.Now()), zap.String("table name", tableName))
		if res := s.db.Table(tableName).AutoMigrate(&CharacterLists{}); res.Error != nil {
			logger.Info("Error to create table", zap.Time("now", time.Now()), zap.Error(res.Error))
			return res.Error
		}
		logger.Info("Table creation done", zap.Time("now", time.Now()))
	}
	tableName = "rankratio"
	if !s.db.HasTable(tableName) {
		logger.Info("No target table. Start to create table",
			zap.Time("now", time.Now()), zap.String("table name", tableName))
		if res := s.db.Table(tableName).AutoMigrate(&RankRatio{}); res.Error != nil {
			logger.Info("Error to create table", zap.Time("now", time.Now()), zap.Error(res.Error))
			return res.Error
		}
		logger.Info("Table creation done", zap.Time("now", time.Now()))
	}
	tableName = "user_characterlist"
	if !s.db.HasTable(tableName) {
		logger.Info("No target table. Start to create table",
			zap.Time("now", time.Now()), zap.String("table name", tableName))
		if res := s.db.Table(tableName).AutoMigrate(&UserCharacterList{}); res.Error != nil {
			logger.Info("Error to create table", zap.Time("now", time.Now()), zap.Error(res.Error))
			return res.Error
		}
		logger.Info("Table creation done", zap.Time("now", time.Now()))
	}
	return nil
}
