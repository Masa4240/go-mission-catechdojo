package userservice

type (
	UserProfile struct {
		ID    int
		Name  string
		Token string
	}

	UserCharacter struct {
		ID   int
		Name string
		Rank string
	}

	UserInfo struct {
		Profile    *UserProfile
		Characters []*UserCharacter
	}
)
