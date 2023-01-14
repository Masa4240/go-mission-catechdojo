package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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
	UserID int
	Name   string
	Rank   string
	//Level int
}

type rankratio struct {
	Ranklevel string
	Weight    int
}

type formalcharLists struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	Name   string
	Rank   string
	Weight int
	Desc   string
}

type masterList struct {
	list []formalcharLists
}

var srChars, rChars, nChars *[]model.CharLists
var rankWeight [3]int

func NewGachaService(db *gorm.DB) *GachaService {
	return &GachaService{
		db: db,
	}
}

func (s *GachaService) Gacha(ctx context.Context, id, count int) ([]model.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	tableName := "charlists_users"

	// Monster Lists

	srChars := *srChars
	rChars := *rChars
	nChars := *nChars

	// Gacha
	newChars := []charLists{}
	newResChars := []model.GachaResponse{}
	logger.Info("Gacha start", zap.Time("now", time.Now()))
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		logger.Info("Gacha number", zap.Time("now", time.Now()), zap.Int("SR ratio", rankWeight[0]), zap.Int("R ratio", rankWeight[1]), zap.Int("N ratio", rankWeight[2]))
		result := rand.Intn(rankWeight[0] + rankWeight[1] + rankWeight[2])

		if result < rankWeight[0] {
			rank := "SR"
			newChars = charGacha(id, rank, srChars, newChars)
		}
		if result >= rankWeight[0] && result < rankWeight[0]+rankWeight[1] {
			rank := "R"
			newChars = charGacha(id, rank, rChars, newChars)
		}
		if result >= rankWeight[0]+rankWeight[1] {
			rank := "N"
			newChars = charGacha(id, rank, nChars, newChars)
		}
	}

	logger.Info("Gacha Finish, Start registration", zap.Time("now", time.Now()))
	// Register to user DB
	for i := 0; i < len(newChars); i++ {
		res := s.db.Table(tableName).Create(&newChars[i])
		if res.Error != nil {
			logger.Info("Failed to register new char to user db", zap.Time("now", time.Now()), zap.Error(res.Error))
			return nil, res.Error
		}
	}
	logger.Info("registration done", zap.Time("now", time.Now()))

	newResChars = resConverter(newChars)

	return newResChars, nil
}

func (s *GachaService) GetCharsList(ctx context.Context, id int) ([]model.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	tableName := "charlists_users"
	charList := []charLists{}
	if result := s.db.Table(tableName).Where("`user_id` = ?", id).Find(&charList); result.Error != nil {
		logger.Info("Fail to get char list from db", zap.Time("now", time.Now()), zap.Error(result.Error))
		return nil, result.Error
	}
	res := resConverter(charList)
	return res, nil
}

//func charGacha(rank string, chars, newChars []charLists, newResChars []model.GachaResponse) ([]charLists, []model.GachaResponse) {
func charGacha(id int, rank string, chars []model.CharLists, newChars []charLists) []charLists {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	newChar := charLists{}
	char := chars[rand.Intn(len(chars))]
	logger.Info("Char is", zap.Time("now", time.Now()), zap.String("Name", char.Name), zap.String("Rank", char.Rank), zap.Int("ID", char.CharID))
	newChar.CharID = int(char.ID)
	newChar.UserID = id
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
	s.GetChars()
	return nil
}

func (s *GachaService) GetChars() error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Get Master char lists", zap.Time("now", time.Now()))

	masterChars := []model.CharLists{}

	if srChars != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		srChars = nil
	}
	if rChars != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		rChars = nil
	}
	if nChars != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		nChars = nil
	}

	var srCharList, rCharList, nCharList []model.CharLists
	s.db.Table("formalchar_lists").Find(&masterChars)
	for i := 0; i < len(masterChars); i++ {
		if masterChars[i].Rank == "SR" {
			if masterChars[i].Weight == 1 {
				srCharList = append(srCharList, masterChars[i])
			}
			if masterChars[i].Weight == 2 {
				srCharList = append(srCharList, masterChars[i])
				srCharList = append(srCharList, masterChars[i])
			}
		}
		if masterChars[i].Rank == "R" {
			if masterChars[i].Weight == 1 {
				rCharList = append(rCharList, masterChars[i])
			}
			if masterChars[i].Weight == 2 {
				rCharList = append(rCharList, masterChars[i])
				rCharList = append(rCharList, masterChars[i])
			}
		}
		if masterChars[i].Rank == "N" {
			if masterChars[i].Weight == 1 {
				nCharList = append(nCharList, masterChars[i])
			}
			if masterChars[i].Weight == 2 {
				nCharList = append(nCharList, masterChars[i])
				nCharList = append(nCharList, masterChars[i])
			}
		}
	}

	srChars = &srCharList
	rChars = &rCharList
	nChars = &nCharList

	fmt.Println(srChars)
	// rank ratio
	rankRatio := []rankratio{}

	if err := s.db.Table("rankratio").Find(&rankRatio).Error; err != nil {
		logger.Info("Error to get Rank ratio", zap.Time("now", time.Now()), zap.Error(err))
		return err
	}
	for i := 0; i < 3; i++ {
		if rankRatio[i].Ranklevel == "SR" {
			rankWeight[0] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "R" {
			rankWeight[1] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "N" {
			rankWeight[2] = rankRatio[i].Weight
		}
	}

	return nil
}
