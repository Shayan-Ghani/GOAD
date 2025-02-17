package itemsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/internal/model"
	sqlrepository "github.com/Shayan-Ghani/GOAD/internal/repository/sql"
	"github.com/Shayan-Ghani/GOAD/pkg/formatter"
	itemrequest "github.com/Shayan-Ghani/GOAD/pkg/request/item"
	tagrequest "github.com/Shayan-Ghani/GOAD/pkg/request/tag"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

const (
	getItemTagsName = "/tags/item/"
	addItemTag      = "/tags/item"
	addTag          = "/tags"
)

type QueryTemplate struct {
	Condition string
	Args      []interface{}
}

type ServiceRepository interface {
	AddItem(name string, description string, dueDate time.Time) (insID int64, err error)

	DeleteItem(id string) error

	UpdateItem(id string, updates map[string]interface{}) error
	UpdateItemStatus(id string) error

	GetItem(id string) (*model.Item, error)
	GetItems(templates ...sqlrepository.QueryTemplate) ([]model.Item, error)
	GetItemByTag(tags []string) ([]model.Item, error)

	GetItemsDone() ([]model.Item, error)
}

type Service struct {
	repo      ServiceRepository
	TagSvcUrl string
}

func NewItemService(repo ServiceRepository, TagSvcUrl string) Service {
	return Service{
		repo:      repo,
		TagSvcUrl: TagSvcUrl,
	}
}

func (s Service) addTagToItem(id int, tags []string) error {

	var itPayload = tagrequest.BasePayload{
		ItemID: strconv.Itoa(id),
		Tags:   tags,
	}

	itemTag, err := json.Marshal(itPayload)

	if err != nil {
		return err
	}

	res, err := http.Post(s.TagSvcUrl+addItemTag, "application/json", bytes.NewBuffer(itemTag))

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s - got code %d", string(body), res.StatusCode)
	}


	var tPayload = tagrequest.Tag{
		Tags: tags,
	}

	t, err := json.Marshal(tPayload)
	if err != nil {
		return err
	}

	tagRes, err := http.Post(s.TagSvcUrl+addTag, "application/json", bytes.NewBuffer(t))

	if err != nil {
		return err
	}

	defer tagRes.Body.Close()

	tBody, err := io.ReadAll(tagRes.Body)

	if err != nil {
		return err
	}

	if tagRes.StatusCode != http.StatusCreated {
		return fmt.Errorf("%s - got code %d", string(tBody), tagRes.StatusCode)
	}

	return nil

}

func (s Service) Add(request itemrequest.Add) error {
	var err error

	if request.Name == "" || request.Description == "" {
		return &validation.ValidationError{
			Field:   "parameters",
			Message: "name or description not provided in the request payload!",
		}
	}

	var t = time.Time{}

	if cmdflag.IsFlagDefined(request.DueDate) {
		t, err = formatter.StringToTime(request.DueDate)

		if err != nil {
			return err
		}
	}

	insID, err := s.repo.AddItem(request.Name, request.Description, t)

	if err != nil {
		return err
	}

	if len(request.Tags) > 0 {
		if err := s.addTagToItem(int(insID), request.Tags); err != nil {
			return err
		}
	}

	return nil
}

type GetResponse struct {
	Err   string       `json:"error,omitempty"`
	Items []model.Item `json:"items"`
}

func (s Service) handleItems(result interface{}, err error) (*GetResponse, error) {

	if err != nil {
		return nil, err
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

	for i := range items {
		res, err := http.Get(s.TagSvcUrl + getItemTagsName + items[i].ID)
		if err != nil || res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("couldn't get tag name for item with id %s, got error: %v", items[i].ID, err)
		}

		body, err := io.ReadAll(res.Body)

		res.Body.Close()

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &items[i].TagsNames); err != nil {
			return nil, err
		}
	}

	fmt.Printf("%+v", items)
	return &GetResponse{
		Items: items,
	}, nil
}

func (s Service) Get() (*GetResponse, error) {
	return s.handleItems(s.repo.GetItems())
}

func (s Service) GetSingle(request itemrequest.Get) (*GetResponse, error) {
	if request.ID != "" {

		return s.handleItems(s.repo.GetItem(request.ID))
	}

	return nil, &validation.ValidationError{
		Field:   "field",
		Message: "id can't be empty",
	}
}

func (s Service) GetDone() (*GetResponse, error) {
	// if request.Done {
	// 	return s.handleItems(s.repo.GetItemsDone())
	// }

	return s.handleItems(s.repo.GetItemsDone())
}

func (s Service) GetByTag(request itemrequest.GetByTag) (*GetResponse, error) {

	// 	tags := formatter.SplitTags(request.Tags)
	// 	return s.handleItems(s.repo.GetItemByTag(tags))
	// }

	return nil, nil
}

func (s Service) Update(request itemrequest.Update) error {
	if len(request.Tags) > 0 {
		id , err := strconv.Atoi(request.ID)
		if err != nil {
			return err
		}

		if err := s.addTagToItem(id, request.Tags); err != nil {
			return err
		}
	}

	if request.Name == "" && request.Description == ""{
		return nil
	}

	updates := make(map[string]interface{}, 3)

	flagUpdates := map[string]interface{}{
		"name":        request.Name,
		"description": request.Description,
	}

	if cmdflag.IsFlagDefined(request.DueDate) {
		t, err := formatter.StringToTime(request.DueDate)

		if err != nil {
			return err
		}

		flagUpdates["due_date"] = t
	}

	fmt.Printf("flagUpdates: %v\n", flagUpdates)

	for key, value := range flagUpdates {
		if value != "" {
			updates[key] = value
		}
	}

	return s.repo.UpdateItem(request.ID, updates)
}

func (s Service) UpdateStatus(request itemrequest.UpdateStatus) error {
	return s.repo.UpdateItemStatus(request.ID)
}

func (s Service) Delete(request itemrequest.Delete) error {
	return s.repo.DeleteItem(request.ID)
}
