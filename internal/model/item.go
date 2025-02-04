package model

import (
	"time"
)

type Item struct {
	ID          string
	Name        string
	Description string
	IsDone      bool
	DueDate  	time.Time
	CreatedAt   time.Time
	ModifiedAt  time.Time
	TagsNames []string
}

type TagItems struct {
	ItemID string
	TagID  string
}