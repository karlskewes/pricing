package pricing_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/karlskewes/pricing"
)

func TestNewInMemoryRepository(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestInMemory_AddPrice(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddBrand(ctx, "EXAMPLE")
	if err != nil {
		t.Fatal(err)
	}

	startDate := time.Now().Round(time.Microsecond).UTC()
	endDate := startDate.Add(1 * time.Hour)

	price := pricing.Price{
		BrandID:   1,
		StartDate: startDate,
		EndDate:   endDate,
		ProductID: 3,
		Priority:  1,
		Price:     100,
		Curr:      "EUR",
	}

	err = db.AddPrice(ctx, price)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestInMemory_GetPrice(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddBrand(ctx, "EXAMPLE")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Table tests for no match, before, inside, after, priority, etc
	// Round time to Microseconds as that is the precision that Postgres
	// supports.
	startDate := time.Now().Round(time.Microsecond).UTC()
	date := startDate.Add(30 * time.Minute)
	endDate := startDate.Add(1 * time.Hour)
	date2 := startDate.Add(2 * time.Hour)

	price := pricing.Price{
		BrandID:   1,
		StartDate: startDate,
		EndDate:   endDate,
		ProductID: 3,
		Priority:  1,
		Price:     100,
		Curr:      "EUR",
	}

	err = db.AddPrice(ctx, price)
	if err != nil {
		t.Fatal(err)
	}

	want := pricing.FinalPrice{
		BrandID:   price.BrandID,
		StartDate: price.StartDate,
		EndDate:   price.EndDate,
		ProductID: price.ProductID,
		Price:     price.Price,
		Curr:      price.Curr,
	}

	// test inside of start & end dates, should be matching price
	got, err := db.GetPrice(ctx, 1, 3, date)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("db.GetPrice(...) mismatch (-want +got):\n%s", diff)
	}

	// test outside of start & end dates, should be no matching prices
	_, err = db.GetPrice(ctx, 1, 3, date2)
	if err == nil {
		t.Errorf("unexpected lack of error")
	}

	if err := db.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestInMemory_AddBrand(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	brandName := "EXAMPLE"

	err = db.AddBrand(ctx, brandName)
	if err != nil {
		t.Error(err)
	}

	// Add a second brand with the same name
	err = db.AddBrand(ctx, brandName)
	if err == nil {
		t.Error("expected duplicate brand error")
	}
}

func TestInMemory_GetBrand(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	brandName := "EXAMPLE"
	// TODO, add collision

	err = db.AddBrand(ctx, brandName)
	if err != nil {
		t.Fatal(err)
	}

	got, err := db.GetBrand(ctx, brandName)
	if err != nil {
		t.Errorf("failed to get brand that should be in repository: %v", err)
	}
	if got.Name != "EXAMPLE" {
		t.Errorf("want: %s - got: %s", brandName, got.Name)
	}
}

func TestInMemory_Shutdown(t *testing.T) {
	imr := &pricing.InMemoryRepository{}
	if err := imr.Shutdown(context.Background()); err != nil {
		t.Errorf("InMemory.Shutdown() error = %v", err)
	}
}

func TestGetPrices(t *testing.T) {
	ctx := context.Background()

	db, err := pricing.NewInMemoryRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddBrand(ctx, "ZARA")
	if err != nil {
		t.Fatal(err)
	}

	prices, err := initialPrices()
	if err != nil {
		t.Fatal(err)
	}

	for _, price := range prices {
		err := db.AddPrice(ctx, price)
		if err != nil {
			t.Fatal(err)
		}
	}

	type test struct {
		date      time.Time
		productID int
		brandID   int
	}

	testCases := map[string]struct {
		test test
		want pricing.FinalPrice
	}{
		"Test 1": {
			test: test{
				date:      time.Date(2020, 06, 14, 10, 0, 0, 0, time.UTC),
				productID: 35455,
				brandID:   1,
			},
			want: pricing.FinalPrice{
				BrandID:   1,
				StartDate: prices[0].StartDate,
				EndDate:   prices[0].EndDate,
				ProductID: 35455,
				Price:     3550,
				Curr:      "EUR",
			},
		},
	}

	for name, tc := range testCases {
		tt := tc
		t.Run(name, func(t *testing.T) {
			got, err := db.GetPrice(context.Background(), tt.test.brandID, tt.test.productID, tt.test.date)
			if err != nil {
				t.Errorf("failed to get price: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("db.GetPrice(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
