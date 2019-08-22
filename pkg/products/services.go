package products

type ProductService struct {
	repository ProductRepository
}

func (service *ProductService) NewProduct(name, description, categoryID *string, price *float64, measurement *UnitOfMeasurement) (*Product, error) {
	product, err := NewProductEntity(name, description, categoryID, price, measurement)

	if err != nil {
		return nil, err
	}

	err = service.repository.Store(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}
