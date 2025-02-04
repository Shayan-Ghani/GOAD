package formatter

import (
	"strings"
	"time"
)

func SplitTags(tags string) []string {
	return strings.Split(tags, ",")
}
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

func StringToTime(t string) (time.Time, error) {
	var tt time.Time
	tt, err := time.Parse("2006-01-02 15:04:05", t)
	
	return tt, err
}
