package gachaservice

type (
	GachaRequest struct {
		ID    int
		Times int
	}
	GachaResponse struct {
		CharacterID int
		Name        string
	}
	GachaContent struct {
		Request    *GachaRequest
		Characters []*GachaResponse
	}
)
