package formatter

import (
	"gocasts/ToDoApp/internal/model"
	"strings"
	"text/tabwriter"
)


func SplitTags(tags string) []string {
	return strings.Split(tags, ",")
}
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

func FormatItemRes(w *tabwriter.Writer, item model.Item) ( status string, tags string) {
	tags = JoinTags(item.TagsNames)
	if tags == "" {
		tags = "No tags"
	}

	status = "Pending"
	if item.IsDone {
		status = "Done"
	}
	return
}