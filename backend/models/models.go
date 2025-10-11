package models

import "time"

type Note struct {
	ID         int       `json:"id"`
	FilePath   string    `json:"path"`
	Title      string    `json:"title"`
	UploadTime time.Time `json:"upload_time"`
}
