package gachaservice

import (
	"crypto/rand"
	"math/big"
	"time"

	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	rankmodel "github.com/Masa4240/go-mission-catechdojo/model/rankratio"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	"go.uber.org/zap"
)

type GachaService struct {
	cmodel  *charactermodel.CharacterModel
	ucmodel *ucmodel.UcModel
	umodel  *usermodel.UserModel
	rmodel  *rankmodel.RankModel
	logger  *zap.Logger
}

func NewGachaService(
	cmodel *charactermodel.CharacterModel,
	ucmodel *ucmodel.UcModel,
	umodel *usermodel.UserModel,
	rmodel *rankmodel.RankModel,
	logger *zap.Logger,
) *GachaService {
	return &GachaService{
		cmodel:  cmodel,
		ucmodel: ucmodel,
		umodel:  umodel,
		rmodel:  rmodel,
		logger:  logger,
	}
}

func (s *GachaService) Gacha(req GachaRequest) ([]*charactermodel.CharacterData, error) {
	count := req.Times
	id := req.ID

	// 1. Gacha and make character list
	rankRatio := rankmodel.GetData()
	max := rankRatio["SR"] + rankRatio["R"] + rankRatio["N"]
	newCharacters := []*ucmodel.UserCharacterList{}
	resCharacters := []*charactermodel.CharacterData{}

	for i := 0; i < count; i++ {
		s.logger.Info("Gacha number", zap.Time("now", time.Now()),
			zap.Int("SR ratio", rankRatio["SR"]), zap.Int("R ratio", rankRatio["N"]),
			zap.Int("N ratio", rankRatio["N"]))
		result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			s.logger.Info("Fail to create rand", zap.Time("now", time.Now()), zap.Error(err))
			return nil, err
		}
		if int(result.Int64()) < rankRatio["SR"] {
			list := charactermodel.GetSRCharacters()
			max := len(list)
			val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
			if err != nil {
				return nil, err
			}
			tempCharacter := *list[val.Int64()]
			newCharacter := ucmodel.UserCharacterList{}
			newCharacter.UserID = id
			newCharacter.CharacterID = tempCharacter.CharacterID
			newCharacters = append(newCharacters, &newCharacter)
			resCharacters = append(resCharacters, &tempCharacter)
		}
		if int(result.Int64()) >= rankRatio["SR"] && int(result.Int64()) < rankRatio["SR"]+rankRatio["R"] {
			list := charactermodel.GetRCharacters()
			max := len(list)
			val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
			if err != nil {
				return nil, err
			}
			tempCharacter := *list[val.Int64()]
			newCharacter := ucmodel.UserCharacterList{}
			newCharacter.UserID = id
			newCharacter.CharacterID = tempCharacter.CharacterID
			newCharacters = append(newCharacters, &newCharacter)
			resCharacters = append(resCharacters, &tempCharacter)
		}
		if int(result.Int64()) >= rankRatio["SR"]+rankRatio["R"] {
			list := charactermodel.GetNCharacters()
			max := len(list)
			val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
			if err != nil {
				return nil, err
			}
			tempCharacter := *list[val.Int64()]
			newCharacter := ucmodel.UserCharacterList{}
			newCharacter.UserID = id
			newCharacter.CharacterID = tempCharacter.CharacterID
			newCharacters = append(newCharacters, &newCharacter)
			resCharacters = append(resCharacters, &tempCharacter)
		}
	}

	// 2. Register Gacha result
	if err := s.ucmodel.RegisterCharacters(newCharacters); err != nil {
		return nil, err
	}

	// if err := s.model.RegisterCharacters(newCharacters); err != nil {
	// 	return nil, err
	// }
	// var newResCharacters []*gachamodel.GachaResponse
	// newResCharacters = resConverter(newCharacters)
	// return newResCharacters, nil
	return resCharacters, nil
}

// func resConverter(characterList []*gachamodel.UserCharacterList) []*gachamodel.GachaResponse {
// 	resCharacters := []*gachamodel.GachaResponse{}
// 	for i := 0; i < len(characterList); i++ {
// 		resCharacter := gachamodel.GachaResponse{}
// 		resCharacter.CharacterID = characterList[i].CharacterID
// 		// resCharacter.Name = gachacontroller.CharacterMasterDataMVC[characterList[i].CharacterID].Name
// 		resCharacters = append(resCharacters, &resCharacter)
// 	}
// 	return resCharacters
// }

func (s *GachaService) InitMasterData() error {
	if err := s.cmodel.GetForamalCharacterList(); err != nil {
		return err
	}
	if err := s.rmodel.GetRankRatio(); err != nil {
		return err
	}
	return nil
}
