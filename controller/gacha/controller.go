package gachacontroller

import (
	"time"

	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"
	"github.com/jinzhu/copier"

	"go.uber.org/zap"
)

type GachaController struct {
	svc    *gachaservice.GachaService
	logger *zap.Logger
}

func NewGachaController(svc *gachaservice.GachaService, logger *zap.Logger) *GachaController {
	return &GachaController{
		svc:    svc,
		logger: logger,
	}
}

func (c *GachaController) Gacha(req GachaReq) ([]*GachaResponse, error) {
	c.logger.Info("Start Gacha Process in controller", zap.Time("now", time.Now()))
	// // Gacha
	// newCharacters := []*
	// max := rankWeightMVC["SR"] + rankWeightMVC["R"] + rankWeightMVC["N"]
	// for i := 0; i < count; i++ {
	// 	c.logger.Info("Gacha number", zap.Time("now", time.Now()),
	// 		zap.Int("SR ratio", rankWeightMVC["SR"]), zap.Int("R ratio", rankWeightMVC["N"]),
	// 		zap.Int("N ratio", rankWeightMVC["N"]))
	// 	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	// 	if err != nil {
	// 		c.logger.Info("Fail to create rand", zap.Time("now", time.Now()), zap.Error(err))
	// 		return nil, err
	// 	}
	// 	if int(result.Int64()) < rankWeightMVC["SR"] {
	// 		newCharacters = characterGachaMVC(id, srCharacters, newCharacters)
	// 	}
	// 	if int(result.Int64()) >= rankWeightMVC["SR"] && int(result.Int64()) < rankWeightMVC["SR"]+rankWeightMVC["R"] {
	// 		newCharacters = characterGachaMVC(id, rCharacters, newCharacters)
	// 	}
	// 	if int(result.Int64()) >= rankWeightMVC["SR"]+rankWeightMVC["R"] {
	// 		newCharacters = characterGachaMVC(id, nCharacters, newCharacters)
	// 	}
	// }

	// c.logger.Info("Gacha Finish, Start registration", zap.Time("now", time.Now()))
	greq := gachaservice.GachaRequest{}
	if err := copier.Copy(&greq, &req); err != nil {
		return nil, err
	}

	result, err := c.svc.Gacha(greq)
	if err != nil {
		c.logger.Info("Error in service", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	res := []*GachaResponse{}
	if err := copier.Copy(&res, &result); err != nil {
		return nil, err
	}
	return res, nil
}

// func characterGachaMVC(id int, characters []*gachamodel.CharacterLists,
// 	newCharacters []*gachamodel.UserCharacterList) []*gachamodel.UserCharacterList {
// 	logger, _ := zap.NewProduction()
// 	defer func(logger *zap.Logger) {
// 		err := logger.Sync()
// 		if err != nil {
// 			return
// 		}
// 	}(logger)
// 	newCharacter := gachamodel.UserCharacterList{}
// 	max := len(characters)
// 	val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
// 	if err != nil {
// 		return nil
// 	}
// 	tempCharacter := *characters[val.Int64()]
// 	logger.Info("Char is", zap.Time("now", time.Now()), zap.Int("ID", newCharacter.CharacterID))
// 	newCharacter.CharacterID = int(tempCharacter.ID)
// 	newCharacter.UserID = id
// 	newCharacters = append(newCharacters, &newCharacter)
// 	return newCharacters
// }

// func (c *GachaController) GetMasterData() error {
// 	logger, _ := zap.NewProduction()
// 	defer func(logger *zap.Logger) {
// 		err := logger.Sync()
// 		if err != nil {
// 			//return err
// 		}
// 	}(logger)
// 	logger.Info("Master Data Initialization", zap.Time("now", time.Now()))
// 	initializeMasterData()
// 	logger.Info("Restore default settings", zap.Time("now", time.Now()))
// 	if err := c.getCharacterMasterData(); err != nil {
// 		return err
// 	}
// 	if err := c.getRankMasterData(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func initializeMasterData() {
// 	logger, _ := zap.NewProduction()
// 	defer func(logger *zap.Logger) {
// 		err := logger.Sync()
// 		if err != nil {
// 			fmt.Println("Err in logger sync")
// 		}
// 	}(logger)
// 	logger.Info("Initialize Master Data", zap.Time("now", time.Now()))
// 	if srCharacters != nil {
// 		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
// 		srCharacters = nil
// 	}
// 	if rCharacters != nil {
// 		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
// 		rCharacters = nil
// 	}
// 	if nCharacters != nil {
// 		logger.Info("Initialize sr char list", zap.Time("now", time.Now()))
// 		nCharacters = nil
// 	}
// 	logger.Info("Initialize Master Data Done", zap.Time("now", time.Now()))
// }

// func (c *GachaController) getCharacterMasterData() error {
// 	masterCharacters, err := c.svc.GetForamalCharacterList()
// 	if err != nil {
// 		return err
// 	}
// 	CharacterMasterDataMVC = map[int]*gachamodel.CharacterLists{}
// 	for i := 0; i < len(masterCharacters); i++ {
// 		CharacterMasterDataMVC[int(masterCharacters[i].ID)] = masterCharacters[i]
// 		if masterCharacters[i].Rank == "SR" {
// 			for j := 0; j < masterCharacters[i].Weight; j++ {
// 				srCharacters = append(srCharacters, masterCharacters[i])
// 			}
// 		}
// 		if masterCharacters[i].Rank == "R" {
// 			for j := 0; j < masterCharacters[i].Weight; j++ {
// 				rCharacters = append(rCharacters, masterCharacters[i])
// 			}
// 		}
// 		if masterCharacters[i].Rank == "N" {
// 			for j := 0; j < masterCharacters[i].Weight; j++ {
// 				nCharacters = append(nCharacters, masterCharacters[i])
// 			}
// 		}
// 	}
// 	return nil
// }

// func (c *GachaController) getRankMasterData() error {
// 	rankRatio, err := c.svc.GetRankMasterData()
// 	if err != nil {
// 		return err
// 	}
// 	rankWeightMVC = map[string]int{}
// 	for i := 0; i < 3; i++ {
// 		if rankRatio[i].Ranklevel == "SR" {
// 			rankWeightMVC["SR"] = rankRatio[i].Weight
// 		}
// 		if rankRatio[i].Ranklevel == "R" {
// 			rankWeightMVC["R"] = rankRatio[i].Weight
// 		}
// 		if rankRatio[i].Ranklevel == "N" {
// 			rankWeightMVC["N"] = rankRatio[i].Weight
// 		}
// 	}
// 	return nil
// }
