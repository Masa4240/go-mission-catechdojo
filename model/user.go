package model

type (
	UserInfo struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}
	// A CreateTODORequest expresses ...
	UserResistrationRequest struct {
		Name string `json:"name"`
	}
	UserResistrationResponse struct {
		Token string `json:"token"`
	}
)
