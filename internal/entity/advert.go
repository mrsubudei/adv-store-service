package entity

type Advert struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        int64    `json:"price"`
	MainPhotoUrl string   `json:"main_photo"`
	PhotosUrls   []string `json:"photos"`
}
