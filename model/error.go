package model

type ErrDuplicatedRequest struct {
	Message string `json:"message"`
}

func (e ErrDuplicatedRequest) Error() string {
	return e.Message
}
