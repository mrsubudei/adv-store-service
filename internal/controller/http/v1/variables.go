package v1

import "github.com/mrsubudei/adv-store-service/internal/entity"

type Answer interface {
	getCode() int
}

type ErrMessage struct {
	Error  string `json:"error,omitempty"`
	Detail string `json:"detail,omitempty"`
	code   int
}

type Response struct {
	Meta *MetaData       `json:"meta_data,omitempty"`
	Data []entity.Advert `json:"data,omitempty"`
	code int
}

type MetaData struct {
	MaxPage int64 `json:"max_page,omitempty"`
}

func (r Response) getCode() int {
	return r.code
}

func (e ErrMessage) getCode() int {
	return e.code
}

const (
	ItemNameExists     = "item with name '%v' already exists"
	JsonNotCorrect     = "json format is not correct"
	NoContentFound     = "no content found with id: "
	WrongQueryRequest  = "queries have wrong value"
	EmptyFiledRequest  = "request has empty fields"
	WrongDataFormat    = "wrong data format"
	AdvertCreated      = "advert created"
	DescLengthExceeded = "'description:' field's length exceeded"
	NameLengthExceeded = "'name:' field's length exceeded"
	UrlsNumberExceeded = "'photo_urls:' field's quantity exceeded"
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
