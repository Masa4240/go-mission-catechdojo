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
		//UserInfo `json:"userinfo"`
		Token string `json:"token"`
	}

	UserGetRequest struct {
		//Name string `json:"name"`
		ID int64 `json:"id"`
	}

	UserGetResponse struct {
		Name string `json:"name"`
	}

	UserUpdateRequest struct {
		Newname string `json:"name"`
		//Tokenname string `json:"tokenname"`
		ID int64 `json:"id"`
	}

	UserUpdateReponse struct {
		Name string `json:"name"`
	}
)
