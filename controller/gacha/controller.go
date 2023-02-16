package gachacontroller

import (
	"math/big"

	"crypto/rand"
	"time"

	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"go.uber.org/zap"
)

type GachaController struct {
	model *gachamodel.GachaModel
}

var (
	characterMasterDataMVC map[int]*gachamodel.CharacterLists
	srCharacters           []*gachamodel.CharacterLists
	rCharacters            []*gachamodel.CharacterLists
	nCharacters            []*gachamodel.CharacterLists
	rankWeightMVC          map[string]int
)

func NewGachaController(model *gachamodel.GachaModel) *GachaController {
	return &GachaController{
		model: model,
	}
}

// このServiceをControllerとServiceに分ける.
func (c *GachaController) Gacha(id, count int) ([]*gachamodel.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	// Gacha
	newCharacters := []*gachamodel.UserCharacterList{}
	// newResCharacters := []*model.GachaResponse{}
	var newResCharacters []*gachamodel.GachaResponse
	logger.Info("Gacha start", zap.Time("now", time.Now()))
	max := rankWeightMVC["SR"] + rankWeightMVC["R"] + rankWeightMVC["N"]
	for i := 0; i < count; i++ {
		// rand.Seed(time.Now().UnixNano())
		logger.Info("Gacha number", zap.Time("now", time.Now()),
			zap.Int("SR ratio", rankWeightMVC["SR"]), zap.Int("R ratio", rankWeightMVC["N"]),
			zap.Int("N ratio", rankWeightMVC["N"]))
		// result := rand.Intn(RankWeightMVC["SR"] + RankWeightMVC["R"] + RankWeightMVC["N"])
		result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			logger.Info("Fail to create rand", zap.Time("now", time.Now()), zap.Error(err))
			return nil, err
		}
		if int(result.Int64()) < rankWeightMVC["SR"] {
			newCharacters = characterGachaMVC(id, srCharacters, newCharacters)
		}
		if int(result.Int64()) >= rankWeightMVC["SR"] && int(result.Int64()) < rankWeightMVC["SR"]+rankWeightMVC["R"] {
			newCharacters = characterGachaMVC(id, rCharacters, newCharacters)
		}
		if int(result.Int64()) >= rankWeightMVC["SR"]+rankWeightMVC["R"] {
			newCharacters = characterGachaMVC(id, nCharacters, newCharacters)
		}
	}

	logger.Info("Gacha Finish, Start registration", zap.Time("now", time.Now()))
	// Register to user DB
	if err := c.model.RegisterCharacters(newCharacters); err != nil {
		logger.Info("Fail character registration", zap.Time("now", time.Now()))
		return nil, err
	}
	logger.Info("registration done", zap.Time("now", time.Now()))
	newResCharacters = resConverter(newCharacters)
	return newResCharacters, nil
}

func (c *GachaController) GetUserCharacterList(id int) ([]*gachamodel.GachaResponse, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	var req usermodel.UserLists
	req.ID = uint(id)

	list, err := c.model.GetCharaterList(&req)
	if err != nil {
		return nil, err
	}
	// res := []*model.GachaResponse{}
	res := resConverter(list)
	return res, nil
}

func characterGachaMVC(id int, characters []*gachamodel.CharacterLists,
	newCharacters []*gachamodel.UserCharacterList) []*gachamodel.UserCharacterList {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
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
		resCharacter.Name = characterMasterDataMVC[characterList[i].CharacterID].Name
		resCharacters = append(resCharacters, &resCharacter)
	}
	return resCharacters
}

func (c *GachaController) AddCharacter(name, rank, desc string, weight int) error {
	logger, _ := zap.NewProduction()
	// charList := formalCharacterList{}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
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
	if err := c.model.AddNewCharacter(&newCharacter); err != nil {
		return err
	}
	// Update master data
	if err := c.GetMasterData(); err != nil {
		return err
	}
	return nil
}

func (c *GachaController) GetMasterData() error {
	// masterCharacters := []*model.CharacterLists{}
	initializeMasterData()

	if err := c.getCharacterMasterData(); err != nil {
		return err
	}

	if err := c.getRankMasterData(); err != nil {
		return err
	}

	return nil
}

func initializeMasterData() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Initialize Master Data", zap.Time("now", time.Now()))
	if srCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		srCharacters = nil
	}
	if rCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		rCharacters = nil
	}
	if nCharacters != nil {
		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
		nCharacters = nil
	}
}

func (c *GachaController) getCharacterMasterData() error {
	masterCharacters, err := c.model.GetForamalCharacterList()
	if err != nil {
		return err
	}
	characterMasterDataMVC = map[int]*gachamodel.CharacterLists{}
	for i := 0; i < len(masterCharacters); i++ {
		characterMasterDataMVC[int(masterCharacters[i].ID)] = masterCharacters[i]
		if masterCharacters[i].Rank == "SR" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				srCharacters = append(srCharacters, masterCharacters[i])
			}
		}

		if masterCharacters[i].Rank == "R" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				rCharacters = append(rCharacters, masterCharacters[i])
			}
		}
		if masterCharacters[i].Rank == "N" {
			for j := 0; j < masterCharacters[i].Weight; j++ {
				nCharacters = append(nCharacters, masterCharacters[i])
			}
		}
	}
	return nil
}

func (c *GachaController) getRankMasterData() error {
	rankRatio, err := c.model.GetRankRatio()
	if err != nil {
		return err
	}
	rankWeightMVC = map[string]int{}
	for i := 0; i < 3; i++ {
		if rankRatio[i].Ranklevel == "SR" {
			rankWeightMVC["SR"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "R" {
			rankWeightMVC["R"] = rankRatio[i].Weight
		}
		if rankRatio[i].Ranklevel == "N" {
			rankWeightMVC["N"] = rankRatio[i].Weight
		}
	}
	return nil
}
