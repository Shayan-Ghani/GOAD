package model

import "time"

type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CreatedAt time.Time `json:"omitempty"`
}