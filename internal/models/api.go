package models

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type JWTToken string
