package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/karlskewes/pricing/storage"
	"github.com/karlskewes/pricing/storage/inmemory"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	db, err := inmemory.New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestInMemory_AddPrice(t *testing.T) {
	ctx := context.Background()

	db, err := inmemory.New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddBrand(ctx, "EXAMPLE")
	if err != nil {
		t.Fatal(err)
	}

	startDate := time.Now().Round(time.Microsecond).UTC()
	endDate := startDate.Add(1 * time.Hour)

	price := storage.Price{
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

	db, err := inmemory.New(ctx)
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

	price := storage.Price{
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

	want := storage.FinalPrice{
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

	db, err := inmemory.New(ctx)
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

	db, err := inmemory.New(ctx)
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
	imr := &inmemory.InMemory{}
	if err := imr.Shutdown(context.Background()); err != nil {
		t.Errorf("InMemory.Shutdown() error = %v", err)
	}
}

func initialPrices() ([]storage.Price, error) {
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

	return []storage.Price{
		{BrandID: 1, StartDate: t1.UTC(), EndDate: t2.UTC(), ProductID: 35455, Priority: 0, Price: 3550, Curr: "EUR"},
		{BrandID: 1, StartDate: t3.UTC(), EndDate: t4.UTC(), ProductID: 35455, Priority: 1, Price: 2545, Curr: "EUR"},
		{BrandID: 1, StartDate: t5.UTC(), EndDate: t6.UTC(), ProductID: 35455, Priority: 1, Price: 3050, Curr: "EUR"},
		{BrandID: 1, StartDate: t7.UTC(), EndDate: t8.UTC(), ProductID: 35455, Priority: 1, Price: 3895, Curr: "EUR"},
	}, nil
}
