package pricing

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Verify interface compliance at compile time
var _ Repository = (*Postgres)(nil)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Postgres is an instance of the database handler and contains a connection
// pool for concurrent use by methods.
type Postgres struct {
	pool *pgxpool.Pool
	// timeoutSeconds int // query timeout for ctx, potentially set at Storage.Service{}
	// logger zerolog.Logger // log SQL queries, etc
}

// NewPostgresRepository returns a Postgres backed Repository for persisting pricing data.
// New also runs any migrations in the ./migrations directory and it does this
// over a single new connection before closing the connection and providing
// a Postgres connection pool for the application main use.
// It might be better to split this functionality and still avoid a race
// condition with connections.
// urlExample := "postgres://username:password@localhost:5432/database_name"
// poolSettingsExample := "?sslmode=verify-ca&pool_max_conns=10"
func NewPostgresRepository(ctx context.Context, URL, poolSettings string) (*Postgres, error) {
	connConfig, err := pgx.ParseConfig(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL config: %w", err)
	}

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection handler: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = runMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close migration db connection: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("%s%s", URL, poolSettings))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL & pool settings: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database via connection pool: %w", err)
	}

	return &Postgres{pool}, nil
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to prepare migrations: %w", err)
	}

	// Set as an config param
	// goose.SetVerbose(true)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Shutdown closes the connection pool to new acquires and waits gracefully for
// existing connections to close.
func (pg *Postgres) Shutdown(ctx context.Context) error {
	// TODO: return early if context cancelled before pool Closed
	pg.pool.Close()

	return nil
}

func (pg *Postgres) AddBrand(ctx context.Context, name string) error {
	sql := `INSERT INTO brand (name) VALUES ($1)`

	_, err := pg.pool.Exec(ctx, sql, name)
	if err != nil {
		return fmt.Errorf("failed to insert brand into database: %w", err)
	}

	return nil
}

func (pg *Postgres) GetBrand(ctx context.Context, name string) (Brand, error) {
	sql := `SELECT (id, name) FROM brand WHERE name=$1`

	var brand Brand
	err := pg.pool.QueryRow(ctx, sql, name).Scan(&brand)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Brand{}, errors.New("no matching brand found")
		}

		return Brand{}, fmt.Errorf("failed to query database: %w", err)
	}

	return brand, nil
}

func (pg *Postgres) AddPrice(ctx context.Context, price Price) error {
	sql := `INSERT INTO price (brand_id, start_date, end_date, product_id, priority, price, curr) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := pg.pool.Exec(ctx, sql, price.BrandID, price.StartDate, price.EndDate, price.ProductID, price.Priority, price.Price, price.Curr)
	if err != nil {
		return fmt.Errorf("failed to insert price into database: %w", err)
	}

	return nil
}

func (pg *Postgres) GetPrice(ctx context.Context, brandID, productID int, date time.Time) (FinalPrice, error) {
	sql := `SELECT start_date, end_date, price, curr FROM price WHERE brand_id=$1 AND product_id=$2 AND start_date<=$3 AND end_date>=$3 ORDER BY priority DESC LIMIT 1`

	fp := FinalPrice{
		BrandID:   brandID,
		ProductID: productID,
	}
	err := pg.pool.QueryRow(ctx, sql, brandID, productID, date).Scan(&fp.StartDate, &fp.EndDate, &fp.Price, &fp.Curr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FinalPrice{}, errors.New("no matching price found")
		}

		return FinalPrice{}, fmt.Errorf("failed to query database: %w", err)
	}

	return fp, nil
}
