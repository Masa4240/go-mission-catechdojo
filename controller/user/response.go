package usercontroller

type (
	UserResistrationRequest struct {
		Name string `json:"name"`
	}

	UserResistrationResponse struct {
		Token string `json:"token"`
	}

	UserGetRequest struct {
		ID int64 `json:"id"`
	}

	UserGetResponse struct {
		Name string `json:"name"`
	}

	UserUpdateRequest struct {
		Newname string `json:"name"`
		ID      int64  `json:"id"`
	}
)
