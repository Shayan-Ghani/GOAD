package repository

import (
	"time"

	"github.com/Shayan-Ghani/GOAD/internal/model"
)

type QueryTemplate struct {
	Condition string
	Args      []interface{}
}

type Repository interface {
	ItemRepository
	TagRepository
	ItemTagRepository
}

type TagRepository interface {
	DeleteTag(name string) error
	GetTags() ([]model.Tag, error)

	AddTagInto(name string) error
	AddTag(tags []string) error
}

type ItemTagRepository interface {
	AddItemTag(id string, tags []string) error
	DeleteItemTags(id string, tags []string) error
	GetItemTagsName(id string) ([]string, error)
	DeleteAllItemTags(id string) error
}

type ItemRepository interface {
	AddItem(name string, description string, dueDate time.Time, tags ...string) error

	DeleteItem(id string) error

	UpdateItem(id string, updates map[string]interface{}) error
	UpdateItemStatus(id string) error

	GetItem(id string) (*model.Item, error)
	GetItems(templates ...QueryTemplate) ([]model.Item, error)
	GetItemByTag(tags []string) ([]model.Item, error)

	GetItemsDone() ([]model.Item, error)
}
