package response

import (
	"encoding/json"
	"fmt"
)

func PrintJson(entities ...any) error {
	for _, entity := range entities {
		ib, err := json.MarshalIndent(entity, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(ib))
	}
	return nil
}

func Json(entity ...any) error {
	i, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	fmt.Println(string(i))
	return nil
}
