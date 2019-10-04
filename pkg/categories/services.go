package categories

import (
	"github.com/alejo-lapix/multimedia-go/banners"
	"github.com/alejo-lapix/multimedia-go/persistence"
)

type StoreCategoryService struct {
	repository CategoryRepository
}

func (service *StoreCategoryService) NewCategory(name, description, parentCategoryID *string, visible *bool, multimedia []*persistence.MultimediaItem, banner *banners.Banner) (*Category, error) {
	category, err := NewCategory(name, description, parentCategoryID, visible, multimedia, banner)

	if err != nil {
		return nil, err
	}

	err = service.repository.Store(category)

	if err != nil {
		return nil, err
	}

	return category, nil
}
