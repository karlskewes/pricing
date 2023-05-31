package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/karlskewes/pricing/storage"
	"github.com/karlskewes/pricing/storage/postgres"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultPostgresImage = "docker.io/postgres:15.3-alpine"
	dbname               = "test-db"
	user                 = "postgres"
	password             = "password"
)

type dbContainer struct {
	testcontainers.Container
	connStr string
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

func setupDB(ctx context.Context) (*dbContainer, error) {
	container, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage(defaultPostgresImage),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	// explicitly set sslmode=disable because the container is not configured to use TLS
	connStr, err := container.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		return nil, err
	}

	return &dbContainer{Container: container, connStr: connStr}, nil
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	dbContainer, err := setupDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	db, err := postgres.New(ctx, dbContainer.connStr, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestAddPrice(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	dbContainer, err := setupDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	db, err := postgres.New(ctx, dbContainer.connStr, "")
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

func TestGetPrice(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	dbContainer, err := setupDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	db, err := postgres.New(ctx, dbContainer.connStr, "")
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

func TestPrices(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	dbContainer, err := setupDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	db, err := postgres.New(ctx, dbContainer.connStr, "")
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
		want storage.FinalPrice
	}{
		"Test 1": {
			test: test{
				date:      time.Date(2020, 06, 14, 10, 0, 0, 0, time.UTC),
				productID: 35455,
				brandID:   1,
			},
			want: storage.FinalPrice{
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

// TODO, AddBrand, GetBrand
