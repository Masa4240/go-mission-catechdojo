package ucservice

import (
	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"go.uber.org/zap"
)

type UcService struct {
	model  *ucmodel.UcModel
	logger *zap.Logger
}

func NewUcService(model *ucmodel.UcModel, logger *zap.Logger) *UcService {
	return &UcService{
		model:  model,
		logger: logger,
	}
}

func (s *UcService) GetUserCharacterList(info userservice.UserInfo) ([]*UserCharacterList, error) {
	var req usermodel.UserList
	req.ID = uint(info.ID)

	list, err := s.model.GetCharaterList(&req)
	if err != nil {
		return nil, err
	}
	var res []*UserCharacterList
	for i := 0; i < len(list); i++ {
		resCharacter := charactermodel.GetCharacterByID(*list[i])
		var tempres UserCharacterList
		tempres.CharacterID = list[i].CharacterID
		tempres.Name = resCharacter.Name
		res = append(res, &tempres)
	}
	return res, nil
}
