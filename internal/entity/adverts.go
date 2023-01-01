package entity

type Advert struct {
	Id           int64    `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Price        int64    `json:"price,omitempty"`
	MainPhotoUrl string   `json:"main_photo_url,omitempty"`
	PhotosUrls   []string `json:"photo_urls,omitempty"`
	CreatedAt    string   `json:"-"`
	MaxCount     int64    `json:"-"`
}

type ContextKey string

const (
	KeyId      ContextKey = "id"
	KeyLimit   ContextKey = "limit"
	KeyOffset  ContextKey = "offset"
	KeySortBy  ContextKey = "sort_by"
	KeyOrderBy ContextKey = "order_by"
	KeyFields  ContextKey = "fields"
)
