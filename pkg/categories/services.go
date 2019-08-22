package categories

type StoreCategoryService struct {
	repository CategoryRepository
}

func (service *StoreCategoryService) NewMainCategory(name *string) (*Category, error) {
	category, err := NewMainCategoryEntity(name)

	if err != nil {
		return nil, err
	}

	err = service.repository.Store(category)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (service *StoreCategoryService) NewSubCategory(name, parentCategoryID *string) (*Category, error) {
	childCategory, err := NewSubCategoryEntity(name, parentCategoryID)

	if err != nil {
		return nil, err
	}

	err = service.repository.Store(childCategory)

	if err != nil {
		return nil, err
	}

	return childCategory, nil
}
