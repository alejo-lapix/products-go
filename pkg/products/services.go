package products

import "github.com/alejo-lapix/multimedia-go/persistence"

type ProductService struct {
	Repository ProductRepository
}

func (service *ProductService) NewProduct(name, description, categoryID *string, price *float64, measurement *UnitOfMeasurement, multimedia []*persistence.MultimediaItem) (*Product, error) {
	product, err := NewProductEntity(name, description, categoryID, price, measurement, multimedia)

	if err != nil {
		return nil, err
	}

	err = service.Repository.Store(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}
