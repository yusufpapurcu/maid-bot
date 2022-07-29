package dto

import "time"

type RandomWaifuList struct {
	Images []struct {
		File          string    `json:"file"`
		Extension     string    `json:"extension"`
		ImageID       int       `json:"image_id"`
		Favourites    int       `json:"favourites"`
		DominantColor string    `json:"dominant_color"`
		Source        string    `json:"source"`
		UploadedAt    time.Time `json:"uploaded_at"`
		IsNsfw        bool      `json:"is_nsfw"`
		Width         int       `json:"width"`
		Height        int       `json:"height"`
		URL           string    `json:"url"`
		PreviewURL    string    `json:"preview_url"`
		Tags          []struct {
			TagID       int    `json:"tag_id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			IsNsfw      bool   `json:"is_nsfw"`
		} `json:"tags"`
	} `json:"images"`
}
