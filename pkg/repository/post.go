package repository

import "time"

type PostData struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ShortDesc string    `json:"short_desc"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
