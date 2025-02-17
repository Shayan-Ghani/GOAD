package requesthandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/pkg/formatter"
	tagrequest "github.com/Shayan-Ghani/GOAD/pkg/request/tag"
)

func (h Handler) handleTag(request CliRequest) error {

	switch request.Command {
	case "view":
		return h.TagView()
	case "delete":
		return h.TagDelete(request)
	default:
		return nil
	}
}

const (
	base          = "/tags"
	addItemTag    = base + "/item"
	addTag        = base
	deleteTag     = base + "/"
	deleteItemTag = base + "/item/"
)

func (h Handler) TagView() error {
	return nil
}

func (h Handler) TagDelete(request CliRequest) error {
	if cmdflag.IsFlagDefined(request.Flags.Name) {
		return h.handleNewRequest(http.MethodDelete, h.TagSvcUrl+deleteTag+request.Flags.Name, nil, http.StatusOK)
	}
	
	if cmdflag.IsFlagDefined(request.Flags.ItemID){
		if !request.Flags.All{
			return h.DeleteItemTags(request)
		}
		
		return h.DeleteAllItemTags(request)
	}
	
	return fmt.Errorf("item id, can't be empty")
}

func (h Handler) DeleteItemTags(request CliRequest) error {

	if cmdflag.IsFlagDefined(request.Flags.ItemID) && cmdflag.IsFlagDefined(request.Flags.Tags) {

		var it = tagrequest.BasePayload{
			Tags: formatter.SplitTags(request.Flags.Tags),
		}

		payload, err := json.Marshal(it)
		if err != nil {
			return err
		}

		return h.handleNewRequest(http.MethodDelete, h.TagSvcUrl+deleteItemTag+request.Flags.ItemID,
			bytes.NewBuffer(payload), http.StatusOK)
	}

	return fmt.Errorf("tags or Item id, can't be empty")
}

func (h Handler) DeleteAllItemTags(request CliRequest) error {
	url := h.TagSvcUrl + deleteItemTag + request.Flags.ItemID + "/all"
	return h.handleNewRequest(http.MethodDelete, url, nil, http.StatusOK)

}
