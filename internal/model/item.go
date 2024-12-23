package model

import (
	"time"
	// "fmt"
)

type Item struct {
	ID          string
	Name        string
	Description string
	IsDone      bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
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