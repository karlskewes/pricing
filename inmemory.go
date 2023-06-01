package pricing

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Verify interface compliance at compile time
var _ Repository = (*InMemoryRepository)(nil)

type InMemoryRepository struct {
	brands map[string]int // brands[name]ID
	prices []Price        // suboptimal data structure
	mutex  sync.RWMutex
	// logger zerolog.Logger // log db queries, etc
}

// NewInMemoryRepository returns a memory backed Repository for persisting pricing data.
func NewInMemoryRepository(ctx context.Context) (*InMemoryRepository, error) {
	brands := make(map[string]int)
	prices := make([]Price, 0)
	return &InMemoryRepository{
		brands: brands,
		prices: prices,
	}, nil
}

func (imr *InMemoryRepository) Shutdown(ctx context.Context) error {
	// NO-OP
	return nil
}

func (imr *InMemoryRepository) AddBrand(ctx context.Context, name string) error {
	if id, ok := imr.brands[name]; ok {
		return fmt.Errorf("brand name already exists: %s with id: %d", name, id)
	}

	imr.mutex.Lock()
	defer imr.mutex.Unlock()

	imr.brands[name] = len(imr.brands) + 1 // start from 1 to match Postgres implementation

	return nil
}

func (imr *InMemoryRepository) GetBrand(ctx context.Context, name string) (Brand, error) {
	id, ok := imr.brands[name]
	if !ok {
		return Brand{}, fmt.Errorf("brand doesn't exist in repository: %s", name)
	}

	brand := Brand{
		Name: name,
		ID:   id,
	}

	return brand, nil
}

func (imr *InMemoryRepository) AddPrice(ctx context.Context, price Price) error {
	imr.mutex.Lock()
	defer imr.mutex.Unlock()
	imr.prices = append(imr.prices, price)

	return nil
}

func (imr *InMemoryRepository) GetPrice(ctx context.Context, brandID, productID int, date time.Time) (FinalPrice, error) {
	// initialize a slice of applicable rates which can filter later
	rates := make([]Price, 0)

	imr.mutex.RLock()
	defer imr.mutex.RUnlock()
	// O(n) walk slice to find suitable items
	for _, price := range imr.prices {
		if price.BrandID != brandID || price.ProductID != productID {
			continue
		}
		if price.StartDate.After(date) && !price.StartDate.Equal(date) {
			continue
		}
		if price.EndDate.Before(date) && !price.EndDate.Equal(date) {
			continue
		}

		rates = append(rates, price)
	}

	if len(rates) == 0 {
		return FinalPrice{}, errors.New("no matching price found")
	}

	// O(m) walk applicable prices and find highest priority
	pvp := rates[0]
	for _, price := range rates {
		if price.Priority > pvp.Priority {
			pvp = price
		}
	}

	return FinalPrice{
		BrandID:   pvp.BrandID,
		StartDate: pvp.StartDate,
		EndDate:   pvp.EndDate,
		ProductID: pvp.ProductID,
		Price:     pvp.Price,
		Curr:      pvp.Curr,
	}, nil
}
