package products

import (
	"github.com/alejo-lapix/multimedia-go/persistence"
	"github.com/google/uuid"
	"time"
)

type UnitOfMeasurement struct {
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
}

type Product struct {
	ID                *string                       `json:"id"`
	Name              *string                       `json:"name"`
	Price             *float64                      `json:"price"`
	Description       *string                       `json:"description"`
	CategoryID        *string                       `json:"categoryId"`
	CreatedAt         *string                       `json:"createdAt"`
	Multimedia        []*persistence.MultimediaItem `json:"multimedia"`
	UnitOfMeasurement *UnitOfMeasurement            `json:"unitOfMeasurement"`
}

func NewProductEntity(name, description, categoryID *string, price *float64, measurement *UnitOfMeasurement) (*Product, error) {
	id := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)

	return &Product{
		ID:                &id,
		Name:              name,
		Price:             price,
		Description:       description,
		CategoryID:        categoryID,
		CreatedAt:         &createdAt,
		UnitOfMeasurement: measurement,
	}, nil
}

type ProductRepository interface {
	Store(*Product) error
	Update(*Product) error
	FindOne(id *string) (*Product, error)
	FindMany(ids []*string) (*Product, error)
	Delete(id *string) error
}
