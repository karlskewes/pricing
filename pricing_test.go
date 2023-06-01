package pricing_test

import (
	"time"

	"github.com/karlskewes/pricing"
)

// TODO mock repository based tests if business logic added to storage.Service

func initialPrices() ([]pricing.Price, error) {
	t1, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-00.00.00")
	if err != nil {
		return nil, err
	}
	t2, err := time.Parse("2006-01-02-15.04.05", "2020-12-31-23.59.59")
	if err != nil {
		return nil, err
	}
	t3, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-15.00.00")
	if err != nil {
		return nil, err
	}
	t4, err := time.Parse("2006-01-02-15.04.05", "2020-06-14-18.30.00")
	if err != nil {
		return nil, err
	}
	t5, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-00.00.00")
	if err != nil {
		return nil, err
	}
	t6, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-11.00.00")
	if err != nil {
		return nil, err
	}
	t7, err := time.Parse("2006-01-02-15.04.05", "2020-06-15-16.00.00")
	if err != nil {
		return nil, err
	}
	t8, err := time.Parse("2006-01-02-15.04.05", "2020-12-31-23.59.59")
	if err != nil {
		return nil, err
	}

	return []pricing.Price{
		{BrandID: 1, StartDate: t1.UTC(), EndDate: t2.UTC(), ProductID: 35455, Priority: 0, Price: 3550, Curr: "EUR"},
		{BrandID: 1, StartDate: t3.UTC(), EndDate: t4.UTC(), ProductID: 35455, Priority: 1, Price: 2545, Curr: "EUR"},
		{BrandID: 1, StartDate: t5.UTC(), EndDate: t6.UTC(), ProductID: 35455, Priority: 1, Price: 3050, Curr: "EUR"},
		{BrandID: 1, StartDate: t7.UTC(), EndDate: t8.UTC(), ProductID: 35455, Priority: 1, Price: 3895, Curr: "EUR"},
	}, nil
}
