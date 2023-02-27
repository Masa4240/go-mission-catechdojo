package gachacontroller

type (
	GachaReq struct {
		ID    int
		Times int `json:"times"`
	}

	GachaResponse struct {
		CharacterID int    `json:"characterID"`
		Name        string `json:"name"`
	}
)
