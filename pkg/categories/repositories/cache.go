package repositories

import (
	"fmt"
	"github.com/alejo-lapix/products-go/pkg/categories"
	"time"
)

type cache interface {
	Put(string, interface{})
	Get(string) (interface{}, error)
	Has(string) bool
	Remember(string, int, func() (interface{}, error)) (interface{}, error)
}

type CacheCategoryRepository struct {
	categories.CategoryRepository
	cache cache
	ttl   int
}

func NewCacheCategoryRepository(repository categories.CategoryRepository, cache cache, ttl int) *CacheCategoryRepository {
	return &CacheCategoryRepository{
		CategoryRepository: repository,
		cache:              cache,
		ttl:                ttl,
	}
}

type inMemory struct {
	elements   map[string]interface{}
	timeStamps map[string]*time.Time
}

func NewInMemoryDriver() *inMemory {
	return &inMemory{
		elements:   map[string]interface{}{},
		timeStamps: map[string]*time.Time{},
	}
}

func (driver *inMemory) Put(key string, elements interface{}) {
	driver.elements[key] = elements
}

func (driver *inMemory) Get(key string) (interface{}, error) {
	element, ok := driver.elements[key]

	if !ok {
		return nil, fmt.Errorf("element \"%s\" not found", key)
	}

	return element, nil
}

func (driver *inMemory) Remember(key string, seconds int, callback func() (interface{}, error)) (interface{}, error) {
	showPullData := false
	timeStamp, ok := driver.timeStamps[key]

	if ok {
		// Time is saved with the requested duration, so it
		// needs to be compared with the current time
		showPullData = timeStamp.Sub(time.Now()) < 0
	}

	elements, err := driver.Get(key)

	if err == nil && !showPullData {
		return elements, nil
	}

	elements, err = callback()

	if err != nil {
		return nil, err
	}

	now := time.Now().Add(time.Second * time.Duration(seconds))
	driver.timeStamps[key] = &now

	driver.Put(key, elements)

	return elements, nil
}

func (driver *inMemory) Has(key string) bool {
	_, ok := driver.elements[key]

	return ok
}

func (repository *CacheCategoryRepository) MainCategories(limit, offset int) ([]*categories.Category, error) {
	signature := fmt.Sprintf("MainCategories %d-%d", limit, offset)
	elements, err := repository.cache.Remember(signature, repository.ttl, func() (interface{}, error) {
		return repository.CategoryRepository.MainCategories(limit, offset)
	})

	if err != nil {
		return nil, err
	}

	return elements.([]*categories.Category), nil
}

func (repository *CacheCategoryRepository) SubCategories(categoryID *string) ([]*categories.Category, error) {
	signature := fmt.Sprintf("SubCategories %s", *categoryID)
	elements, err := repository.cache.Remember(signature, repository.ttl, func() (interface{}, error) {
		return repository.CategoryRepository.SubCategories(categoryID)
	})

	if err != nil {
		return nil, err
	}

	return elements.([]*categories.Category), nil
}

func (repository *CacheCategoryRepository) Find(ID *string) (*categories.Category, error) {
	signature := fmt.Sprintf("Find %s", *ID)
	elements, err := repository.cache.Remember(signature, repository.ttl, func() (interface{}, error) {
		return repository.CategoryRepository.Find(ID)
	})

	if err != nil {
		return nil, err
	}

	return elements.(*categories.Category), nil
}

func (repository *CacheCategoryRepository) FindMany(ids []*string) ([]*Category, error) {
	return repository.CategoryRepository.FindMany(ids)
}

func (repository *CacheCategoryRepository) Store(category *categories.Category) error {
	return repository.CategoryRepository.Store(category)
}

func (repository *CacheCategoryRepository) Remove(ID *string) error {
	return repository.CategoryRepository.Remove(ID)
}

func (repository *CacheCategoryRepository) Update(ID *string, category *categories.Category) error {
	return repository.CategoryRepository.Update(ID, category)
}

func (repository *CacheCategoryRepository) All() ([]*categories.Category, error) {
	elements, err := repository.cache.Remember("All", repository.ttl, func() (interface{}, error) {
		return repository.CategoryRepository.All()
	})

	if err != nil {
		return nil, err
	}

	return elements.([]*categories.Category), nil
}

func (repository *CacheCategoryRepository) Total() (int64, error) {
	return repository.CategoryRepository.Total()
}
