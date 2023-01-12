package service

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type GachaService struct {
	db *gorm.DB
}

type charLists struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	CharID int
	Name   string
	Rank   string
	//Level int
}

type rankratio struct {
	Ranklevel string
	Weight    int
}

type rowCharLists struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	CharID int
	Name   string
	Rank   string
	Level  int
}

type formalcharLists struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	Name   string
	Rank   string
	Weight int
	Desc   string
}

func NewGachaService(db *gorm.DB) *GachaService {
	return &GachaService{
		db: db,
	}
}

func (s *GachaService) Gacha(ctx context.Context, id, count int) ([]model.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	tableName := "charlists_userid_" + strconv.Itoa(id)

	// Confirm table existance of target user. If no, create table for that user
	if !s.db.HasTable(tableName) {
		logger.Info("No target table. Start to create table", zap.Time("now", time.Now()), zap.String("table name", "charlists_userid_"+strconv.Itoa(id)))
		if err := s.db.Table(tableName).AutoMigrate(&charLists{}).Error; err != nil {
			logger.Info("Error to create table", zap.Time("now", time.Now()), zap.Error(err))
			return nil, err
		}
		logger.Info("Table creation done", zap.Time("now", time.Now()))
	}

	// Monster Lists
	srChars := s.getChars("SR")
	rChars := s.getChars("R")
	nChars := s.getChars("N")
	// rank ratio
	rankratio := rankratio{}
	s.db.Table("rankratio").Where("ranklevel=?", "SR").Find(&rankratio)
	srRatio := rankratio.Weight
	s.db.Table("rankratio").Where("ranklevel=?", "R").Find(&rankratio)
	rRatio := rankratio.Weight
	s.db.Table("rankratio").Where("ranklevel=?", "N").Find(&rankratio)
	nRatio := rankratio.Weight

	// Gacha
	newChars := []charLists{}
	newResChars := []model.GachaResponse{}

	for i := 0; i < count; i++ {
		logger.Info("Gacha number", zap.Time("now", time.Now()), zap.Int("Gacha Counter", i))
		rand.Seed(time.Now().UnixNano())

		logger.Info("Gacha number", zap.Time("now", time.Now()), zap.Int("SR ratio", srRatio), zap.Int("R ratio", rRatio), zap.Int("N ratio", nRatio))

		result := rand.Intn(srRatio + rRatio + nRatio)

		if result < srRatio {
			rank := "SR"
			newChars = charGacha(rank, srChars, newChars)
		}
		if result >= srRatio && result < srRatio+rRatio {
			rank := "R"
			newChars = charGacha(rank, rChars, newChars)
		}
		if result >= srRatio+rRatio {
			rank := "N"
			newChars = charGacha(rank, nChars, newChars)
		}
	}

	// Register to user DB
	for i := 0; i < len(newChars); i++ {
		res := s.db.Table(tableName).Create(&newChars[i])
		if res.Error != nil {
			logger.Info("Failed to register new char to user db", zap.Time("now", time.Now()), zap.Error(res.Error))
			return nil, res.Error
		}
	}

	newResChars = resConverter(newChars)

	return newResChars, nil
}

func (s *GachaService) GetCharsList(ctx context.Context, id int) ([]model.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	tableName := "charlists_userid_" + strconv.Itoa(id)
	charList := []charLists{}
	if result := s.db.Table(tableName).Find(&charList); result.Error != nil {
		logger.Info("Fail to get char list from db", zap.Time("now", time.Now()), zap.Error(result.Error))
		return nil, result.Error
	}
	res := resConverter(charList)
	return res, nil
}

//func charGacha(rank string, chars, newChars []charLists, newResChars []model.GachaResponse) ([]charLists, []model.GachaResponse) {
func charGacha(rank string, chars, newChars []charLists) []charLists {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	newChar := charLists{}
	char := chars[rand.Intn(len(chars))]
	logger.Info("Char is", zap.Time("now", time.Now()), zap.String("Name", char.Name), zap.String("Rank", char.Rank), zap.Int("ID", char.CharID))
	newChar.CharID = int(char.ID)
	newChar.Name = char.Name
	newChar.Rank = rank
	newChars = append(newChars, newChar)

	return newChars //, newResChars
}

func resConverter(charList []charLists) []model.GachaResponse {
	resChar := model.GachaResponse{}
	resChars := []model.GachaResponse{}
	for i := 0; i < len(charList); i++ {
		resChar.CharID = int(charList[i].CharID)
		resChar.Name = charList[i].Name
		resChars = append(resChars, resChar)
	}
	return resChars
}

func (s *GachaService) AddCharacter(ctx context.Context, name, rank, desc string, weight int) error {
	logger, _ := zap.NewProduction()
	charList := formalcharLists{}
	defer logger.Sync()
	logger.Info("Start new char reg Process", zap.Time("now", time.Now()), zap.String("Name", name), zap.String("Rank", rank), zap.String("Desc", desc), zap.Int("Weight", weight))
	if res := s.db.Table("formalchar_lists").Find(&charList, "name=?", name); res.Error == nil {
		logger.Info("char name already exists", zap.Time("now", time.Now()), zap.String("Existing", name), zap.Uint("ID number:", charList.ID))
		err := errors.New("Duplicated Name")
		return err
	}
	newChar := formalcharLists{}
	newChar.Name = name
	newChar.Rank = rank
	newChar.Desc = desc
	newChar.Weight = weight
	res := s.db.Create(&newChar)
	if res.Error != nil {
		logger.Info("Fail to update DB", zap.Time("now", time.Now()))
		return res.Error
	}
	if res.RowsAffected != 1 {
		logger.Info("Fail to add new char to DB", zap.Time("now", time.Now()), zap.String("NewName is", name), zap.Int64("Affected row", res.RowsAffected))
		err := errors.New("Fail to create data to DB")
		return err
	}
	return nil
}

func (s *GachaService) getChars(rank string) []charLists {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("get character lists from db", zap.Time("now", time.Now()))

	charList := []charLists{}
	rankChars := []charLists{}

	s.db.Table("formalchar_lists").Where("`rank` = ? AND weight = ?", rank, 1).Find(&rankChars)
	charList = append(rankChars, charList...)
	s.db.Table("formalchar_lists").Where("`rank` = ? AND weight = ?", rank, 2).Find(&rankChars)
	charList = append(rankChars, charList...)
	charList = append(rankChars, charList...)

	logger.Info("Finish character lists from db", zap.Time("now", time.Now()), zap.String("rank", rank))
	return charList
}
