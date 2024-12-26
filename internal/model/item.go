package model

import (
	"time"
)

type Item struct {
	ID          string
	Name        string
	Description string
	IsDone      bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
	TagsNames []string
}

type Tag struct {
	ID    string
	Name  string
	CreatedAt time.Time
}

type TagItems struct {
	ItemID string
	TagID  string
}