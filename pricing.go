// package pricing implements implements storing and retrieving pricing data.
package pricing

import (
	"context"
	"time"
)

type Price struct {
	BrandID   int       // BRAND_ID: foreign key of the group chain (1 = EXAMPLE).
	StartDate time.Time // START_DATE: date range in which the indicated price applies.
	EndDate   time.Time // END_DATE: date range in which the indicated price applies.
	ProductID int       // PRODUCT_ID: Product code identifier.
	Priority  int       // PRIORITY: Price application disambiguator. If two prices coincide in a date range, the one with higher priority (higher numerical value) is applied.
	Price     int       // PRICE: final selling price. Lowest unit for currency, e.g: cents // could be money.Money
	Curr      string    // CURR: currency iso.
}

type FinalPrice struct {
	BrandID   int       // BRAND_ID: foreign key of the group chain (1 = EXAMPLE).
	StartDate time.Time // START_DATE: date range in which the indicated price applies.
	EndDate   time.Time // END_DATE: date range in which the indicated price applies.
	ProductID int       // PRODUCT_ID: Product code identifier.
	Price     int       // PRICE: final selling price. Lowest unit for currency, e.g: cents
	Curr      string    // CURR: currency iso.
}

type Brand struct {
	ID   int
	Name string
}

// Service contains a Repository and actions any business logic before/after
// interacting with the Repository.
type Service struct {
	repo Repository
	// timeoutSeconds int // query timeout for ctx, potentially per endpoint
	// logger zerolog.Logger // log business logic, etc
}

// NewService creates a new Service with the supplied repository ready for
// reading and writing Price data.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (srv *Service) AddBrand(ctx context.Context, name string) error {
	// TODO: Any business logic common to Repositories
	// TODO: Add any timeout to ctx
	return srv.repo.AddBrand(ctx, name)
}

func (srv *Service) GetBrand(ctx context.Context, name string) (Brand, error) {
	// TODO: Any business logic common to Repositories
	// TODO: Add any timeout to ctx
	return srv.repo.GetBrand(ctx, name)
}

// AddPrice inserts a new Price into the backing storage repository.
func (srv *Service) AddPrice(ctx context.Context, price Price) error {
	// TODO: Any business logic common to Repositories
	// TODO: Add any timeout to ctx
	return srv.repo.AddPrice(ctx, price)
}

// GetPrice returns the final price to apply given the provided brand, product
// and date. Price is an integer in the currencies lowest common demoninator,
// For example, cents in USD, yen in JPY.
func (srv *Service) GetPrice(ctx context.Context, brandID, productID int, date time.Time) (FinalPrice, error) {
	// TODO: Any business logic common to Repositories
	// TODO: Add any timeout to ctx
	return srv.repo.GetPrice(ctx, brandID, productID, date)
}

// Consider
// func (srv *Service) DeleteBrand(...)
// func (srv *Service) DeletePrice(...)
