package requesthandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	
	"github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/internal/model"
	"github.com/Shayan-Ghani/GOAD/pkg/formatter"
	itemrequest "github.com/Shayan-Ghani/GOAD/pkg/request/item"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

const (
	getItemTagsName = "/tags/"
	baseEndpoint    = "/items"
	
	addItem = baseEndpoint
	addDone = getDone
	
	getItems      = baseEndpoint
	getItemSingle = getItems + "/"
	getDone       = getItems + "/done"
	
	updateItem   = baseEndpoint
	updateStatus = updateItem + "/done"
	
	deleteItem  = baseEndpoint + "/"
	contentType = "application/json"
)




func (h Handler) handleItem(request CliRequest) (*CliResponse, error) {
	
	switch request.Command {
	case "add":
		return nil, h.ItemAdd(request)
	case "view":
		return h.ItemGet(request)
	case "delete":
		return nil, h.ItemDelete(request)
	case "update":
		return nil, h.ItemUpdate(request)
	case "done":
		return nil, h.ItemDone(request)
	default:
		return nil, fmt.Errorf("unknown command : %v", request.Command)
	}
}

func (h Handler) handleNewRequest(method string, url string, body io.Reader, statusCode int) error {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response: %v", err)
	}

	if resp.StatusCode != statusCode {

		return fmt.Errorf("faild with Status: %s", string(data))
	}
	return nil
}

func (h Handler) ItemAdd(request CliRequest) error {
	var err error

	var tags []string

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		tags = formatter.SplitTags(request.Flags.Tags)
	}

	var ir = itemrequest.Add{
		BasePayload: itemrequest.BasePayload{
			Name:        request.Flags.Name,
			Description: request.Flags.Description,
			Tags:        tags,
			DueDate:     request.Flags.DueDate,
		},
	}

	payload, err := json.Marshal(ir)

	if err != nil {
		return err
	}

	resp, err := h.client.Post(h.TagSvcUrl+addItem, contentType, bytes.NewBuffer(payload))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return &validation.ValidationError{
			Field:   "request",
			Message: fmt.Sprintf("could add item - got err : %s", resp.Status),
		}
	}

	return nil
}

func (h Handler) ItemGet(request CliRequest) (*CliResponse, error) {

	handleResponse := func(endpoint string) (*CliResponse, error) {
		var url = h.ItemSvcUrl + endpoint

		res, err := h.client.Get(url)
		if err != nil || res.StatusCode != http.StatusOK {
			return nil, &validation.ValidationError{
				Field:   "requests",
				Message: fmt.Sprintf("couldn't get items, got error: %v, url : %s", err, url),
			}
		}

		body, err := io.ReadAll(res.Body)

		defer res.Body.Close()

		if err != nil {
			return nil, err
		}

		var i []model.Item

		if err := json.Unmarshal(body, &i); err != nil {
			return nil, err
		}

		return &CliResponse{
			Items: i,
		}, nil

	}

	if request.Flags.Done {
		return handleResponse(getDone)
	}

	if request.Flags.All {
		return handleResponse(getItems)
	}

	if err := validation.ValidateFlagsDefinedStr([]string{"-i"}, request.Flags.ID); err != nil {
		return nil, err
	}

	return handleResponse(getItemSingle + request.Flags.ID)
}

func (h Handler) ItemUpdate(request CliRequest) error {
	var err error

	var tags []string

	if cmdflag.IsFlagDefined(request.Flags.Tags) {
		tags = formatter.SplitTags(request.Flags.Tags)
	}

	var iU = itemrequest.Update{
		BasePayload: itemrequest.BasePayload{
			Name:        request.Flags.Name,
			Description: request.Flags.Description,
			DueDate:     request.Flags.DueDate,
			Tags:        tags,
		},
		ID: request.Flags.ID,
	}

	payload, err := json.Marshal(iU)

	if err != nil {
		return err
	}

	return h.handleNewRequest(http.MethodPut, h.ItemSvcUrl+updateItem, bytes.NewBuffer(payload), http.StatusNoContent)
}

func (h Handler) ItemDone(request CliRequest) error {

	if cmdflag.IsFlagDefined(request.Flags.ID) {
		var i = itemrequest.UpdateStatus{
			ID: request.Flags.ID,
		}

		body, err := json.Marshal(i)
		if err != nil {
			return err
		}

		resp, err := h.client.Post(h.ItemSvcUrl+addDone, contentType, bytes.NewBuffer(body))

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("updating item status failed with %v, got response:  ", resp.Status)
		}

	}

	return fmt.Errorf("id can't be empty")
}

func (h Handler) ItemDelete(request CliRequest) error {

	if cmdflag.IsFlagDefined(request.Flags.ID) {
		return h.handleNewRequest(http.MethodDelete, h.ItemSvcUrl+deleteItem+request.Flags.ID, nil, http.StatusOK )
	}

	return fmt.Errorf("id can't be empty")
}

