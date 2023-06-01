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
	svc *Service
	// logger, etc
}

func NewApp(args []string) (*App, error) {
	// parse flags for configuration
	fs := flag.NewFlagSet("pricing", flag.ContinueOnError)

	enablePostgres := fs.Bool("enable-postgres", false, "use postgres pricing repository")
	dbConnStr := fs.String("db-conn-str", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable", "database connection string")
	dbPoolSettings := fs.String("db-pool-settings", "", "database pool settings")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, fmt.Errorf("unable to parse flags: %w", err)
	}

	var repo Repository
	if *enablePostgres {
		postgres, err := NewPostgresRepository(context.Background(), *dbConnStr, *dbPoolSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to create new postgres database pool: %w", err)
		}

		repo = postgres
	} else {
		imr, err := NewInMemoryRepository(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to create new in-memory repository: %w", err)
		}

		repo = imr
	}

	svc := NewService(repo)

	mux := http.NewServeMux()
	// add middleware for Prometheus metrics, logging, OTEL, etc
	// TODO: replace NotFoundHandler with API Handler
	mux.Handle("/", http.NotFoundHandler())

	app := &App{
		srv: &http.Server{
			Addr:    defaultListenAddr,
			Handler: mux,
		},
		svc: svc,
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
