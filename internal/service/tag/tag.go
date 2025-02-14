package tagsvc

import (
	"github.com/Shayan-Ghani/GOAD/internal/model"
	request "github.com/Shayan-Ghani/GOAD/pkg/request/tag"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

type ServiceRepository interface {
	DeleteTag(name string) error
	GetTags() ([]model.Tag, error)

	AddTagInto(name string) error
	AddTag(tags []string) error

	AddItemTag(id string, tags []string) error
	DeleteItemTags(id string, tags []string) error
	GetItemTagsName(id string) ([]string, error)
	DeleteAllItemTags(id string) error
}

type Service struct {
	repo ServiceRepository
}

func NewTagService(repo ServiceRepository) Service {
	return Service{
		repo: repo,
	}
}

type BaseRequest struct {
	Name        string
}


func (s Service) Get() ([]model.Tag, error) {
	_, err := s.repo.GetTags()
	return nil, err
}

func (s Service) Delete(request request.Delete) error {
	if request.Name == "" {
		return &validation.ValidationError{
			Field: "Json Field",
			Message: "name either not defined or empty",
		}
	}

	return s.repo.DeleteTag(request.Name)
}


func (s Service) Add(request request.Tag) error {
	
	if len(request.Tags) < 1 {
		return &validation.ValidationError{
			Field: "Json Field",
			Message: "itemID or Tags not defined or empty",
		}
	}
	
	return s.repo.AddTag(request.Tags)
}

func (s Service) AddToItem(request request.BasePayload) error {

	if request.ItemID == "" || len(request.Tags) < 1 {
		return &validation.ValidationError{
			Field: "Json Field",
			Message: "itemID or Tags not defined or empty",
		}
	}

	return s.repo.AddItemTag(request.ItemID, request.Tags)
}


func (s Service) GetFromItems(request request.Base) ([]string , error) {
	if request.ItemID == "" {
		return nil, &validation.ValidationError{
			Field: "Json Field",
			Message: "itemID not defined or empty",
		}
	}

	return s.repo.GetItemTagsName(request.ItemID)
	
}


func (s Service) DeleteFromItem(request request.BasePayload) error {
	
	if request.ItemID == "" || len(request.Tags) < 1 {
		return &validation.ValidationError{
			Field: "Json Field",
			Message: "itemID or Tags not defined or empty",
		}
	}	
	return s.repo.DeleteItemTags(request.ItemID,request.Tags)
}



func (s Service) DeleteAllFromItem(request request.Base) error {

	if request.ItemID == "" {
		return &validation.ValidationError{
			Field: "Json Field",
			Message: "itemID not defined or empty",
		}
	}


	return s.repo.DeleteAllItemTags(request.ItemID)
}