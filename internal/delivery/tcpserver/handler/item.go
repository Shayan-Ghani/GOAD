package tcphandler

import (
	"fmt"
	"log"
	"time"

	"github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/internal/model"
	"github.com/Shayan-Ghani/GOAD/pkg/formatter"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

func (h *RequestHandler) handleItem(request CliRequest) CliResponse {
	handleResult := func(err error) CliResponse {
		if err == nil {
			return CliResponse{}
		}

		if cliResp, ok := err.(*CliResponse); ok {
			return *cliResp
		}

		return CliResponse{
			Err: err.Error(),
		}
	}

	switch request.Command {
	case "add":
		return handleResult(h.handleItemAdd(request))
	case "view":
		return handleResult(h.handleItemGet(request))
	case "delete":
		return handleResult(h.handleItemDelete(request))
	case "update":
		return handleResult(h.handleItemUpdate(request))
	case "done":
		return handleResult(h.handleItemDone(request))
	default:
		return CliResponse{
			Err: "unknown command",
		}
	}
}

func (h *RequestHandler) handleItemAdd(request CliRequest) error {
	var err error

	var tags []string

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		tags = formatter.SplitTags(request.Flags.Tags)
	}

	var t = time.Time{}

	if cmdflag.IsFlagDefined(request.Flags.DueDate) {
		t, err = formatter.StringToTime(request.Flags.DueDate)

		if err != nil {
			return err
		}
	}

	err = h.repo.AddItem(request.Flags.Name, request.Flags.Description, t, tags...)

	return err
}

func (h *RequestHandler) handleItemGet(request CliRequest) error {
	handleItems := func(result interface{}, err error) error {
		if err != nil {
			return err
		}

		var items []model.Item
		switch v := result.(type) {
		case []model.Item:
			items = v
		case *model.Item:
			if v != nil {
				items = []model.Item{*v}
			}
		default:
			log.Println("got Unhandled Item type!")
		}

		return &CliResponse{
			Items: items,
		}
	}

	if request.Flags.Done {
		return handleItems(h.repo.GetItemsDone())
	}
	if request.Flags.All {
		return handleItems(h.repo.GetItems())
	}

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		tags := formatter.SplitTags(request.Flags.Tags)
		return handleItems(h.repo.GetItemByTag(tags))
	}

	if err := validation.ValidateFlagsDefinedStr([]string{"-i"}, request.Flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	return handleItems(h.repo.GetItem(request.Flags.ID))
}

func (h *RequestHandler) handleItemUpdate(request CliRequest) error {
	var err error

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		h.repo.AddItemTag(request.Flags.ID, formatter.SplitTags(request.Flags.Tags))
	}

	updates := make(map[string]interface{}, 3)

	flagUpdates := map[string]interface{}{
		"name":        request.Flags.Name,
		"description": request.Flags.Description,
	}

	if cmdflag.IsFlagDefined(request.Flags.DueDate) {
		t, err := formatter.StringToTime(request.Flags.DueDate)

		if err != nil {
			return err
		}

		flagUpdates["due_date"] = t
	}

	for key, value := range flagUpdates {
		updates[key] = value
	}

	err = h.repo.UpdateItem(request.Flags.ID, updates)

	return err
}

func (h *RequestHandler) handleItemDone(request CliRequest) error {
	return h.repo.UpdateItemStatus(request.Flags.ID)
}

func (h *RequestHandler) handleItemDelete(request CliRequest) error {
	var err error

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		err = h.repo.DeleteItemTags(request.Flags.ID, formatter.SplitTags(request.Flags.Tags))
		return err
	}

	if request.Flags.DelTags {
		err = h.repo.DeleteAllItemTags(request.Flags.ID)
		return err
	}

	err = h.repo.DeleteItem(request.Flags.ID)

	return err
}
