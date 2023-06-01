package pricing

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	httpShutdownPreStopDelaySeconds = 1
	httpShutdownTimeoutSeconds      = 1
	defaultListenAddr               = "localhost:8080"
)

type App struct {
	srv *http.Server
	// logger, etc
}

func NewApp(args []string) (*App, error) {
	fs := flag.NewFlagSet("pricing", flag.ContinueOnError)

	enablePostgres := fs.Bool("enable-postgres", false, "use postgres pricing repository")
	dbConnStr := fs.String("db-conn-str", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable", "database connection string")
	dbPoolSettings := fs.String("db-pool-settings", "", "database pool settings")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, fmt.Errorf("unable to parse flags: %w", err)
	}

	ctx := context.Background()

	var repo Repository
	if *enablePostgres {
		postgres, err := NewPostgresRepository(ctx, *dbConnStr, *dbPoolSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to create new postgres database pool: %w", err)
		}

		repo = postgres
	} else {
		imr, err := NewInMemoryRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create new in-memory repository: %w", err)
		}

		repo = imr
	}
	svc := NewService(repo)

	seedPrices, err := initialPrices()
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial pricing data: %w", err)
	}

	if err := svc.repo.AddBrand(ctx, "EXAMPLE"); err != nil {
		return nil, fmt.Errorf("failed to add brand: %w", err)
	}

	for _, price := range seedPrices {
		if err := svc.repo.AddPrice(ctx, price); err != nil {
			return nil, fmt.Errorf("failed to add an initial price to repository: %w", err)
		}
	}

	handler, err := NewHandler(svc)
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing handler: %w", err)
	}

	mux := http.NewServeMux()
	// add middleware for Prometheus metrics, logging, OTEL, etc
	// TODO: replace NotFoundHandler with API Handler
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/api/v1/brands", handler.GetBrand)
	mux.HandleFunc("/api/v1/prices", handler.GetPrice)

	app := &App{
		srv: &http.Server{
			Addr:              defaultListenAddr,
			Handler:           mux,
			IdleTimeout:       30 * time.Second,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			// etc
		},
	}

	return app, nil
}

// Run starts an HTTP server and gracefully shuts down when the provided
// context is marked done.
func (app *App) Run(ctx context.Context) error {
	var group errgroup.Group

	group.Go(func() error {
		<-ctx.Done()

		// Before shutting down the HTTP server wait for any HTTP requests that are
		// in transit on the network. Common in Kubernetes and other distributed
		// systems.
		time.Sleep(httpShutdownPreStopDelaySeconds * time.Second)

		// Give active connections time to complete or disconnect before closing.
		ctx2, cancel := context.WithTimeout(ctx, httpShutdownTimeoutSeconds*time.Second)
		defer cancel()

		return app.srv.Shutdown(ctx2)
	})

	group.Go(func() error {
		err := app.srv.ListenAndServe()
		// http.ErrServerClosed is expected at shutdown.
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	return group.Wait()
}

func initialPrices() ([]Price, error) {
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

	return []Price{
		{BrandID: 1, StartDate: t1.UTC(), EndDate: t2.UTC(), ProductID: 35455, Priority: 0, Price: 3550, Curr: "EUR"},
		{BrandID: 1, StartDate: t3.UTC(), EndDate: t4.UTC(), ProductID: 35455, Priority: 1, Price: 2545, Curr: "EUR"},
		{BrandID: 1, StartDate: t5.UTC(), EndDate: t6.UTC(), ProductID: 35455, Priority: 1, Price: 3050, Curr: "EUR"},
		{BrandID: 1, StartDate: t7.UTC(), EndDate: t8.UTC(), ProductID: 35455, Priority: 1, Price: 3895, Curr: "EUR"},
	}, nil
}
