package command

func isFlagDefined(flags ...string) bool {
	for _, f := range flags {
		if f == "" {
			return false
		}
	}
	return true
}

type Flags struct {
	Name        string
	ID          string
	Description string
	Tags        string
	Short       string
	Format      string
	All         bool
	Done        bool
	DelTags     bool
}
