package response

import (
	"fmt"
	"gocasts/ToDoApp/internal/model"
	"os"
	"strings"
	"text/tabwriter"
)

func TabWriter(arg ...any) {
	fmt.Println("The Game Begins.")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID\tName\tDescription\tStatus\tTags")
	fmt.Fprintln(w, "--\t----\t-----------\t------\t----")

	for _, argItem := range arg {
		switch items := argItem.(type) {
		case []model.Item:
			if len(items) == 0 {
				fmt.Println("No items found!")
				continue
			}
			for _, item := range items {
				printItem(w, item)
			}
		case *model.Item:
			printItem(w, *items)
		default:
			fmt.Printf("TabWriter: unexpected type %T\n", argItem)
		}
	}
	w.Flush()

}

func printItem(w *tabwriter.Writer, item model.Item) {
	tags := strings.Join(item.TagsNames, ", ")
	if tags == "" {
		tags = "No tags"
	}

	status := "Pending"
	if item.IsDone {
		status = "Done"
	}

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
		item.ID,
		item.Name,
		item.Description,
		status,
		tags,
	)
}
