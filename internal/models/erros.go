package models

import "errors"

var (
	ErrUserExists           = errors.New("user exists with this login")
	ErrUserPasswordNotValid = errors.New("user password not match")
)

var (
	ErrOrderExist                        = errors.New("order exists")
	ErrOrderNotFound                     = errors.New("order not found")
	ErrInsufficientBalance               = errors.New("insufficient balance")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrOrderInProcessed                  = errors.New("order in processed")
)
