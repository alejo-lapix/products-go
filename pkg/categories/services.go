package categories

import "github.com/alejo-lapix/multimedia-go/persistence"

type StoreCategoryService struct {
	repository CategoryRepository
}

func (service *StoreCategoryService) NewCategory(name, description, parentCategoryID *string, multimedia []*persistence.MultimediaItem) (*Category, error) {
	category, err := NewCategory(name, description, parentCategoryID, multimedia)

	if err != nil {
		return nil, err
	}

	err = service.repository.Store(category)

	if err != nil {
		return nil, err
	}

	return category, nil
}
