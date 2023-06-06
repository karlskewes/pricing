package pricing

import (
	"context"
	"fmt"
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

func SeedExampleData(ctx context.Context, repo Repository) error {
	t1, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-00.00.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t2, err := time.Parse("2006-01-02-15.04.05", "2020-12-31-23.59.59")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t3, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-15.00.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t4, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-18.30.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t5, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-00.00.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t6, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-11.00.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t7, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-16.00.00")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t8, err := time.Parse("2006-01-02-15.04.05", "2020-12-31-23.59.59")
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	prices := []Price{
		{BrandID: 1, StartDate: t1.UTC(), EndDate: t2.UTC(), ProductID: 35455, Priority: 0, Price: 3550, Curr: "EUR"},
		{BrandID: 1, StartDate: t3.UTC(), EndDate: t4.UTC(), ProductID: 35455, Priority: 1, Price: 2545, Curr: "EUR"},
		{BrandID: 1, StartDate: t5.UTC(), EndDate: t6.UTC(), ProductID: 35455, Priority: 1, Price: 3050, Curr: "EUR"},
		{BrandID: 1, StartDate: t7.UTC(), EndDate: t8.UTC(), ProductID: 35455, Priority: 1, Price: 3895, Curr: "EUR"},
	}

	if err := repo.AddBrand(ctx, "EXAMPLE"); err != nil {
		return fmt.Errorf("failed to add brand: %w", err)
	}

	for _, price := range prices {
		if err := repo.AddPrice(ctx, price); err != nil {
			return fmt.Errorf("failed to add an initial price to repository: %w", err)
		}
	}

	return nil
}
