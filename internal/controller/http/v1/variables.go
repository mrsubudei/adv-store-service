package v1

type ErrMessage struct {
	Code    int
	Message string
}

type Key string

const (
	EmptyFiledRequest  = "request has empty fields"
	WrongDataFormat    = "wrong data format"
	AdvertCreated      = "advert created"
	DescLengthExceeded = "description length exceeded"
	NameLengthExceeded = "name length exceeded"
	UrlsNumberExceeded = "photo_urls quantity exceeded"
)
