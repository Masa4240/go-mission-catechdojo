package rankmodel

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type RankModel struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRankModel(db *gorm.DB, logger *zap.Logger) *RankModel {
	return &RankModel{
		db:     db,
		logger: logger,
	}
}

var (
	rankWeight map[string]int
)

func (s *RankModel) GetRankRatio() error {
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
		return err
	}
	tx.Commit()

	rankWeight = map[string]int{}
	for i := 0; i < 3; i++ {
		if rankRatio[i].Ranklevel == "SR" {
			rankWeight["SR"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "R" {
			rankWeight["R"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "N" {
			rankWeight["N"] = rankRatio[i].Weight
		}
	}
	return nil
}

func GetData() map[string]int {
	return rankWeight
}
