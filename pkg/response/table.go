package response

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/Shayan-Ghani/GOAD/internal/model"
)

func PrintTable(arg ...any) {
	fmt.Println("The Game Begins.")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID\tName\tDescription\tStatus\tDue_Date\tTags\tCreated_At")
	fmt.Fprintln(w, "--\t----\t-----------\t------\t----\t------\t------")

	for _, argItem := range arg {
		switch entity := argItem.(type) {
		case []model.Item:
			for _, item := range entity {
				i, err := NewItemRes(&item)
				if err != nil {
					log.Fatalln(err)
				}
				printItem(w, *i)
			}

		case *model.Item:
			i, err := NewItemRes(entity)
			if err != nil {
				log.Fatalln(err)
			}
			printItem(w, *i)

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

func printItem(w *tabwriter.Writer, item ItemResponse) {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		item.ID,
		item.Name,
		item.Description,
		item.IsDone,
		item.DueDate,
		item.TagsNames,
		item.CreatedAt,
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
