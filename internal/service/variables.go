package service

import "errors"

var (
	ErrUniqueName = errors.New("UNIQUE constraint failed: adverts.name")
)
