package service

const (
	UniqueNameConstraint = "UNIQUE constraint failed: adverts.name"
	DateFormat           = "2006-01-02 15:04:05"
)

type ContextKey string

const (
	KeyId      ContextKey = "id"
	KeyLimit   ContextKey = "limit"
	KeyOffset  ContextKey = "offset"
	KeySortBy  ContextKey = "sort"
	KeyOrderBy ContextKey = "order"
)
