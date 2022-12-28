package entity

import "errors"

var (
	ErrNameAlreadyExist = errors.New("name already exists")
	ErrItemNotExists    = errors.New("item not exists")
	ErrUpdateFailed     = errors.New("update failed")
	ErrDeleteFailed     = errors.New("delete failed")
)
