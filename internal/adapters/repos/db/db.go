package db

import (
	"context"
	"database/sql"
	"fluxara/internal/config"
	"fluxara/internal/domain"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DbdAdapter struct {
	db *sql.DB
}

func NewDbAdapter(configs *config.Config) (*DbdAdapter, error) {
	conn, err := connectToDb(configs)
	if err != nil {
		log.Panic("Error al conectar a db desde adapter")
		return nil, err
	}

	db := newDbPQLDB(conn)

	return db, nil
}

func newDbPQLDB(db *sql.DB) *DbdAdapter {
	return &DbdAdapter{db: db}
}

func connectToDb(config *config.Config) (*sql.DB, error) {
	var conn *sql.DB
	var err error

	for i := 1; i <= config.Db.Retries; i++ {
		connDB := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Db.Host, config.Db.Port, config.Db.User,
			config.Db.Password, config.Db.Name, config.Db.SslMode,
		)

		fmt.Printf("[CONN] Esto es connDB %s \n", connDB)

		conn, err = sql.Open("postgres", connDB)
		if err == nil {
			err = conn.Ping()
		}

		if err == nil {
			break
		}

		log.Printf("retry %d/%d: error conectando a DB: %v", i, config.Db.Retries, err)
		time.Sleep(20 * time.Second)
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(30 * time.Minute) // Reciclá cada tanto para evitar conns zombis
	conn.SetConnMaxIdleTime(5 * time.Minute)

	return conn, nil
}

func (dPQLDB *DbdAdapter) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	var product domain.Product

	query := `
		SELECT
			id,
			sku,
			name,
			category,
			price_cents,
			stock,
			is_active,
			created_at
		FROM public.products
		WHERE id = $1
		AND is_active = true
	`

	ctx = context.Background()
	err := dPQLDB.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Category,
		&product.PriceCents,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, fmt.Errorf("producto no encontrado")
		}
		return product, err
	}

	return product, nil
}

func (dPQLDB *DbdAdapter) GetProductsAll(ctx context.Context, id string,
) ([]domain.Product, error) {

	products := []domain.Product{}

	query := `
		SELECT
			id,
			sku,
			name,
			category,
			price_cents,
			stock,
			is_active,
			created_at
		FROM public.products
		WHERE id = $1
		  AND is_active = true
	`

	rows, err := dPQLDB.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product

		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Category,
			&product.PriceCents,
			&product.Stock,
			&product.IsActive,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no se encontraron productos")
	}

	return products, nil
}
