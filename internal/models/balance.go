package models

type Balance struct {
	ID             int     `json:"-"`
	UserID         int     `json:"-"`
	CurrentBalance float32 `json:"current"`
	Withdrawn      float32 `json:"withdrawn"`
}
