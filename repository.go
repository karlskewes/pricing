package pricing

import (
	"context"
	"time"
)

// Repository implements persisting and reading pricing data from a backend.
type Repository interface {
	AddPrice(ctx context.Context, price Price) error
	GetPrice(ctx context.Context, brandID, productID int, date time.Time) (FinalPrice, error)
	AddBrand(ctx context.Context, name string) error
	GetBrand(ctx context.Context, name string) (Brand, error)
	Shutdown(ctx context.Context) error
}

type MockRepository struct {
	// Could hold a few prices and brands so can test presence/absence
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (mr *MockRepository) AddPrice(ctx context.Context, price Price) error {
	// TODO not implemented
	return nil
}

func (mr *MockRepository) GetPrice(ctx context.Context, brandID, productID int, date time.Time) (FinalPrice, error) {
	return FinalPrice{
		BrandID:   brandID,
		ProductID: productID,
		StartDate: date,
		EndDate:   date.Add(24 * time.Hour),
		Price:     100,
		Curr:      "USD",
	}, nil
}

func (mr *MockRepository) AddBrand(ctx context.Context, name string) error {
	return nil
}

func (mr *MockRepository) GetBrand(ctx context.Context, name string) (Brand, error) {
	return Brand{
		ID:   1234,
		Name: name,
	}, nil
}

func (mr *MockRepository) Shutdown(ctx context.Context) error {
	return nil
}
