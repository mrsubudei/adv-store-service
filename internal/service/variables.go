package service

import "errors"

var (
	ErrUniqueName = errors.New("UNIQUE constraint failed: adverts.name")
	DateFormat    = "2006-01-02 15:04:05"
)
