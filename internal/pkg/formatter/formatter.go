package formatter

import (
	"gocasts/ToDoApp/internal/model"
	"strings"
	"text/tabwriter"
	"time"
)


func SplitTags(tags string) []string {
	return strings.Split(tags, ",")
}
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

func FormatItemRes(w *tabwriter.Writer, item model.Item) ( status string, tags string, createdAt string, err error) {
	tags = JoinTags(item.TagsNames)
	if tags == "" {
		tags = "No tags"
	}

	status = "Pending"
	if item.IsDone {
		status = "Done"
	}
	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		return "", "", "",err
	}
	createdAtLocal := item.CreatedAt.In(localLocation)

	createdAt = createdAtLocal.Format("2006-01-02 15:04:05")
	return
}