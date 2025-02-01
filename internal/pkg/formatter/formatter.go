package formatter

import (
	"strings"
)

func SplitTags(tags string) []string {
	return strings.Split(tags, ",")
}
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

