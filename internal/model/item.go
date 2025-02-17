package model

import (
	"time"
)

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDone      bool      `json:"is_done"`
	DueDate     time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ModifiedAt  time.Time `json:"modified_at,omitempty"`
	TagsNames   []string  `json:"tags"`
}

type ItemTag struct {
	ItemID string
	TagID  string
}
