package categories

import (
	"github.com/alejo-lapix/multimedia-go/persistence"
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID               *string                       `json:"id"`
	Name             *string                       `json:"name" validate:"required"`
	Description      *string                       `json:"description"`
	Multimedia       []*persistence.MultimediaItem `json:"multimedia"`
	ParentCategoryID *string                       `json:"parentCategoryId,omitempty"`
	IsMainCategory   *string                       `json:"isMainCategory"`
	Visible          *bool                         `json:"visible"`
	CreatedAt        *string                       `json:"createdAt"`
}

func createdAt() *string {
	createdAt := time.Now().Format(time.RFC3339)

	return &createdAt
}

func NewCategory(name, description, parentCategoryID *string, visible *bool, multimedia []*persistence.MultimediaItem) (*Category, error) {
	id := uuid.New().String()
	isMainCategory := "y"

	if parentCategoryID != nil && *parentCategoryID != "" {
		isMainCategory = "n"
	}

	category := &Category{
		ID:               &id,
		Name:             name,
		Description:      description,
		ParentCategoryID: parentCategoryID,
		Multimedia:       multimedia,
		Visible:          visible,
		IsMainCategory:   &isMainCategory,
		CreatedAt:        createdAt(),
	}

	return category, nil
}

func (category *Category) AddMultimediaItem(item *persistence.MultimediaItem) {
	category.Multimedia = append(category.Multimedia, item)
}

func (category *Category) RemoveMultimediaItem(id *string) bool {
	newList := make([]*persistence.MultimediaItem, len(category.Multimedia))
	itemExists := false

	for _, item := range category.Multimedia {
		if *item.ID != *id {
			newList = append(newList, item)
		} else {
			itemExists = true
		}
	}

	return itemExists
}

type Commitable interface {
	Commit() error
}

type CategoryRepository interface {
	// MainCategories shows only the visible categories that does
	// not have a parent category, it's useful for the end user
	MainCategories(limit, offset int) ([]*Category, error)

	// SubCategories shows the categories related to other one
	SubCategories(categoryID *string) ([]*Category, error)
	Find(ID *string) (*Category, error)
	Store(*Category) error
	Remove(ID *string) error
	Update(ID *string, category *Category) error
	All() ([]*Category, error)
	Total() (int64, error)
}
