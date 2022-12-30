package v1

import "github.com/mrsubudei/adv-store-service/internal/entity"

type Answer interface {
	getCode() int
}

type ErrMessage struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
	code   int
}

type SingleResponse struct {
	Data entity.Advert `json:"data"`
	code int
}

type MultiResponse struct {
	Data []entity.Advert `json:"data"`
	code int
}

func (r SingleResponse) getCode() int {
	return r.code
}

func (r MultiResponse) getCode() int {
	return r.code
}

func (e ErrMessage) getCode() int {
	return e.code
}

type Pagination struct {
	Limit   int
	Offset  int
	SortBy  string
	OrderBy string
}

const (
	WrongQueryRequest  = "queries have wrong value"
	EmptyFiledRequest  = "request has empty fields"
	WrongDataFormat    = "wrong data format"
	AdvertCreated      = "advert created"
	DescLengthExceeded = "description length exceeded"
	NameLengthExceeded = "name length exceeded"
	UrlsNumberExceeded = "photo_urls quantity exceeded"
)

const (
	QueryFields         = "fields"
	QueryLimit          = "limit"
	QueryOffset         = "offset"
	QuerySortBy         = "sort_by"
	QueryOrderBy        = "order_by"
	QueryValueTrue      = "true"
	QueryValueAsc       = "asc"
	QueryValueDesc      = "desc"
	QueryValueCreatedAt = "create_at"
	QueryValuePrice     = "price"
)
