package response

import (
	"fmt"
	"gocasts/ToDoApp/internal/model"
	"gocasts/ToDoApp/internal/pkg/formatter"
	"os"
	"text/tabwriter"
)

func TabWriter(arg ...any) {
	fmt.Println("The Game Begins.")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID\tName\tDescription\tStatus\tTags\tCreated_At")
	fmt.Fprintln(w, "--\t----\t-----------\t------\t----\t----")
	
	for _, argItem := range arg {
		switch entity := argItem.(type) {
		case []model.Item:
			if len(entity) == 0 {
				fmt.Println("No items found!")
				continue
			}
			for _, item := range entity {
				printItem(w, item)
			}
		case *model.Item:
			printItem(w, *entity)
		case []model.Tag:
			if len(entity) == 0 {
				fmt.Println("No tag found!")
				continue
			}
			for _, tag := range entity {
				printTag(w, tag)
			}
		default:
			fmt.Printf("TabWriter: unexpected type %T\n", argItem)
		}
	}
	w.Flush()

}

func printItem(w *tabwriter.Writer, item model.Item) {
	status, tags, createdAt, err := formatter.FormatItemRes(w, item)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		item.ID,
		item.Name,
		item.Description,
		status,
		tags,
		createdAt,
	)
}
func printTag(w *tabwriter.Writer, tag model.Tag) {
	fmt.Fprintln(w, "Name\tCreated_At")
	fmt.Fprintln(w, "--\t------")
	
	fmt.Fprintf(w, "%s\t%s\n",
		tag.Name,
		tag.CreatedAt,
	)
}
