package categories

import (
	"github.com/alejo-lapix/multimedia-go/persistence"
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID               *string                       `json:"id"`
	Name             *string                       `json:"name" validate:"required"`
	Multimedia       []*persistence.MultimediaItem `json:"multimedia"`
	ParentCategoryID *string                       `json:"parentCategoryId"`
	CreatedAt        *string                       `json:"createdAt"`
	subCategories    []*Category
}

func createdAt() *string {
	createdAt := time.Now().Format(time.RFC3339)

	return &createdAt
}

func NewMainCategoryEntity(name *string) (*Category, error) {
	id := uuid.New().String()

	return &Category{
		ID:        &id,
		Name:      name,
		CreatedAt: createdAt(),
	}, nil
}

func NewSubCategoryEntity(name, parentCategoryID *string) (*Category, error) {
	category, err := NewMainCategoryEntity(name)

	if err != nil {
		return nil, err
	}

	category.ParentCategoryID = parentCategoryID

	return category, nil
}

func (category *Category) AddMultimediaItem(item *persistence.MultimediaItem) {
	category.Multimedia = append(category.Multimedia, item)
}

func (category *Category) RemoveMultimediaItem(id *string) bool {
	newList := make([]*persistence.MultimediaItem, len(category.Multimedia))
	itemExists := false

	for _, item := range category.Multimedia {
		if *item.Key() != *id {
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
	MainCategories(limit, offset int) ([]*Category, error)
	SubCategories(categoryID int) ([]*Category, error)
	Find(ID *string) (*Category, error)
	Store(*Category) error
	Remove(ID *string) error
	Update(ID *string, category *Category) error
}
