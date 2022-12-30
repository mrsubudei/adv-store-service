package entity

type Advert struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        int64    `json:"price"`
	MainPhotoUrl string   `json:"main_photo_url"`
	PhotosUrls   []string `json:"photo_urls"`
	CreatedAt    string   `json:"-"`
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
