package entity

import "errors"

var (
	ErrNameAlreadyExist = errors.New("name already exists")
	ErrItemNotExists    = errors.New("item not exists")
	ErrItemsNotExist    = errors.New("there are no items")
	ErrUpdateFailed     = errors.New("update failed")
	ErrDeleteFailed     = errors.New("delete failed")
)
