package v1

import "time"

type noteCreateRequest struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

type noteResponse struct {
	ID        string    `json:"uuid"`
	Title     string    `json:"title"`
	Desc      string    `json:"description"`
	IsDel     bool      `json:"is_delete"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
