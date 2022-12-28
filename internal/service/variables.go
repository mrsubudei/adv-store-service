package service

import "errors"

const (
	ErrUniqueName = errors.New("UNIQUE constraint failed: adverts.name")
)
