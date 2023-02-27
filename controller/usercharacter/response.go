package uccontroller

type (
	UserCharacterReq struct {
		ID int
	}

	UserCharacterRes struct {
		CharacterID int    `json:"characterID"`
		Name        string `json:"name"`
	}
)
