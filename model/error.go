package model

type DuplicatedRequestError struct {
	Message string `json:"message"`
}

func (e DuplicatedRequestError) Error() string {
	return e.Message
}
