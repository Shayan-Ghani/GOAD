package response

import (
	"fmt"
	"gocasts/ToDoApp/internal/model"
	"gocasts/ToDoApp/pkg/formatter"
	"time"
)

type ItemResponse struct {
	model.Item
	IsDone    string
	TagsNames string
	DueDate   string
	CreatedAt string
}

func NewItemRes(item *model.Item) (*ItemResponse, error) {
	tags := formatter.JoinTags(item.TagsNames)
	if tags == "" {
		tags = "No tags"
	}

	status := "Pending"
	if item.IsDone {
		status = "Done"
	}

	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}
	cLocal := item.CreatedAt.In(localLocation).Format("2006-01-02 15:04:05")

	dLocal := item.DueDate.In(localLocation).Format("2006-01-02 15:04:05")

	if item.DueDate.IsZero() {
		dLocal = "Not Set"
	}

	return &ItemResponse{
		Item:      *item,
		IsDone:    status,
		DueDate:   dLocal,
		TagsNames: tags,
		CreatedAt: cLocal,
	}, nil
}

func Respond(format string, args ...any) {
	if hasNoRecords(args) {
		fmt.Println("**Found 0 Records!**")
		return
	}

	if format == "table" {
		PrintTable(args...)
	} else {
		PrintJson(args...)
	}
}

func hasNoRecords(args []any) bool {
	for _, arg := range args {
		switch v := arg.(type) {
		case []model.Item:
			if len(v) == 0 {
				return true
			}
		case []model.Tag:
			if len(v) == 0 {
				return true
			}
		}
	}
	return false
}
