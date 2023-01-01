package entity

import "errors"

var (
	ErrNameAlreadyExist = errors.New("name already exists")
	ErrItemNotExists    = errors.New("item does not exist")
	ErrNoItems          = errors.New("there are no items")
)
