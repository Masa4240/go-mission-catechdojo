package gachaservice

import (
	"context"
	"math/big"

	"crypto/rand"
	"time"

	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"go.uber.org/zap"
)

type GachaService struct {
	svc *gachamodel.GachaModel
}

var (
	CharacterMasterDataMVC map[int]*gachamodel.CharacterLists
	SrCharacters           []*gachamodel.CharacterLists
	RCharacters            []*gachamodel.CharacterLists
	NCharacters            []*gachamodel.CharacterLists
	RankWeightMVC          map[string]int
)

func NewGachaService(svc *gachamodel.GachaModel) *GachaService {
	return &GachaService{
		svc: svc,
	}
}

// このServiceをControllerとServiceに分ける
func (s *GachaService) Gacha(ctx context.Context, id, count int) ([]*gachamodel.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	// if err := s.svc.GachaTableCheck(ctx); err != nil {
	// 	logger.Info("Error table handling", zap.Time("now", time.Now()))
	// 	return nil, err
	// }

	// Gacha
	newCharacters := []*gachamodel.UserCharacterList{}
	// newResCharacters := []*model.GachaResponse{}
	var newResCharacters []*gachamodel.GachaResponse
	logger.Info("Gacha start", zap.Time("now", time.Now()))
	max := RankWeightMVC["SR"] + RankWeightMVC["R"] + RankWeightMVC["N"]
	for i := 0; i < count; i++ {
		// rand.Seed(time.Now().UnixNano())
		logger.Info("Gacha number", zap.Time("now", time.Now()),
			zap.Int("SR ratio", RankWeightMVC["SR"]), zap.Int("R ratio", RankWeightMVC["N"]),
			zap.Int("N ratio", RankWeightMVC["N"]))
		// result := rand.Intn(RankWeightMVC["SR"] + RankWeightMVC["R"] + RankWeightMVC["N"])
		result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			logger.Info("Fail to create rand", zap.Time("now", time.Now()), zap.Error(err))
			return nil, err
		}
		if int(result.Int64()) < RankWeightMVC["SR"] {
			newCharacters = characterGachaMVC(id, SrCharacters, newCharacters)
		}
		if int(result.Int64()) >= RankWeightMVC["SR"] && int(result.Int64()) < RankWeightMVC["SR"]+RankWeightMVC["R"] {
			newCharacters = characterGachaMVC(id, RCharacters, newCharacters)
		}
		if int(result.Int64()) >= RankWeightMVC["SR"]+RankWeightMVC["R"] {
			newCharacters = characterGachaMVC(id, NCharacters, newCharacters)
		}
	}

	logger.Info("Gacha Finish, Start registration", zap.Time("now", time.Now()))
	// Register to user DB
	if err := s.svc.RegisterCharacters(newCharacters); err != nil {
		logger.Info("Fail character registration", zap.Time("now", time.Now()))
		return nil, err
	}
	logger.Info("registration done", zap.Time("now", time.Now()))
	newResCharacters = resConverter(newCharacters)
	return newResCharacters, nil
}

func (s *GachaService) GetUserCharacterList(ctx context.Context, id int) ([]*gachamodel.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	var req usermodel.UserLists
	req.ID = uint(id)

	list, err := s.svc.GetCharaterList(ctx, &req)
	if err != nil {
		return nil, err
	}
	// res := []*model.GachaResponse{}
	res := resConverter(list)
	return res, nil
}

func characterGachaMVC(id int, characters []*gachamodel.CharacterLists, newCharacters []*gachamodel.UserCharacterList) []*gachamodel.UserCharacterList {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	newCharacter := gachamodel.UserCharacterList{}
	// no := rand.Intn(len(characters))
	max := len(characters)
	val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return nil
	}
	tempCharacter := *characters[val.Int64()]
	logger.Info("Char is", zap.Time("now", time.Now()), zap.Int("ID", newCharacter.CharacterID))
	newCharacter.CharacterID = int(tempCharacter.ID)
	newCharacter.UserID = id

	newCharacters = append(newCharacters, &newCharacter)

	return newCharacters
}

func resConverter(characterList []*gachamodel.UserCharacterList) []*gachamodel.GachaResponse {
	resCharacters := []*gachamodel.GachaResponse{}
	for i := 0; i < len(characterList); i++ {
		resCharacter := gachamodel.GachaResponse{}
		resCharacter.CharacterID = characterList[i].CharacterID
		// resCharacter.Name = CharacterMasterDataMVC[int(characterList[i].CharacterID)].Name
		resCharacter.Name = CharacterMasterDataMVC[characterList[i].CharacterID].Name
		resCharacters = append(resCharacters, &resCharacter)
	}
	return resCharacters
}

func (s *GachaService) AddCharacter(ctx context.Context, name, rank, desc string, weight int) error {
	logger, _ := zap.NewProduction()
	// charList := formalCharacterList{}
	defer logger.Sync()
	logger.Info("Start new char reg Process", zap.Time("now", time.Now()),
		zap.String("Name", name), zap.String("Rank", rank), zap.String("Desc", desc), zap.Int("Weight", weight))
	// if err := s.svc.CharacterTableCheck(ctx); err != nil {
	// 	return err
	// }

	newCharacter := gachamodel.CharacterLists{}
	newCharacter.Name = name
	newCharacter.Rank = rank
	newCharacter.Desc = desc
	newCharacter.Weight = weight
	if err := s.svc.AddNewCharacter(&newCharacter); err != nil {
		return err
	}
	// Update master data
	if err := s.GetMasterData(); err != nil {
		return err
	}
	return nil
}

func (s *GachaService) GetMasterData() error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Get Master character lists", zap.Time("now", time.Now()))

	// masterCharacters := []*model.CharacterLists{}

	if SrCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		SrCharacters = nil
	}
	if RCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		RCharacters = nil
	}
	if NCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		NCharacters = nil
	}

	masterCharacters, err := s.svc.GetForamalCharacterList()
	if err != nil {
		return err
	}

	CharacterMasterDataMVC = map[int]*gachamodel.CharacterLists{}
	for i := 0; i < len(masterCharacters); i++ {
		CharacterMasterDataMVC[int(masterCharacters[i].ID)] = masterCharacters[i]
		if masterCharacters[i].Rank == "SR" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				SrCharacters = append(SrCharacters, masterCharacters[i])
			}
		}

		if masterCharacters[i].Rank == "R" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				RCharacters = append(RCharacters, masterCharacters[i])
			}
		}
		if masterCharacters[i].Rank == "N" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				NCharacters = append(NCharacters, masterCharacters[i])
			}
		}
	}
	rankRatio, err := s.svc.GetRankRatio()
	if err != nil {
		return err
	}
	RankWeightMVC = map[string]int{}
	for i := 0; i < 3; i++ {
		if rankRatio[i].Ranklevel == "SR" {
			RankWeightMVC["SR"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "R" {
			RankWeightMVC["R"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "N" {
			RankWeightMVC["N"] = rankRatio[i].Weight
		}
	}
	return nil
}
