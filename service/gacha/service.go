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

func (s *GachaService) Gacha(req GachaRequest) (*GachaContent, error) {
	res := GachaContent{
		Request:    &req,
		Characters: nil,
	}

	// 1. Gacha and make character list
	rankRatio := rankmodel.GetData()
	max := rankRatio["SR"] + rankRatio["R"] + rankRatio["N"]

	for i := 0; i < req.Times; i++ {
		s.logger.Info("Gacha number", zap.Time("now", time.Now()),
			zap.Int("SR ratio", rankRatio["SR"]), zap.Int("R ratio", rankRatio["N"]),
			zap.Int("N ratio", rankRatio["N"]))
		result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			s.logger.Info("Fail to create rand", zap.Time("now", time.Now()), zap.Error(err))
			return nil, err
		}
		if int(result.Int64()) < rankRatio["SR"] {
			res, err = getCharacter("SR", res)
		}
		if int(result.Int64()) >= rankRatio["SR"] && int(result.Int64()) < rankRatio["SR"]+rankRatio["R"] {
			res, err = getCharacter("R", res)
		}
		if int(result.Int64()) >= rankRatio["SR"]+rankRatio["R"] {
			res, err = getCharacter("N", res)
		}
	}
	newCharacters := convertToUserCharacter(res)

	// 2. Register Gacha result
	if err := s.ucmodel.RegisterCharacters(newCharacters); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *GachaService) InitMasterData() error {
	if err := s.cmodel.GetForamalCharacterList(); err != nil {
		return err
	}
	if err := s.rmodel.GetRankRatio(); err != nil {
		return err
	}
	return nil
}

func getCharacter(rank string, gacha GachaContent) (GachaContent, error) {
	list := charactermodel.GetCharacters(rank)
	max := len(list)
	val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return gacha, err
	}
	tempCharacter := GachaResponse{
		CharacterID: list[val.Int64()].CharacterID,
		Name:        list[val.Int64()].Name,
	}
	gacha.Characters = append(gacha.Characters, &tempCharacter)
	return gacha, nil
}
