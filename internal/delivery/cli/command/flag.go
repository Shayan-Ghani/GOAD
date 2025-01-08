package command

func isFlagDefined(flags ...string) bool {
	for _, f := range flags{
		if f == "" {
			return false
		}
	}
	return true
}

type Flags struct {
	Name        string
	All         string
	ID          string
	Description string
	Tags        string
	Short       string
	TagNames    string
	Done        bool
}
