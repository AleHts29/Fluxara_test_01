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

type DbdAdapterGergal struct {
	db *sql.DB
}

func NewDbAdapterGergal(configs *config.Config) (*DbdAdapter, error) {
	conn, err := connectToDb(configs)
	if err != nil {
		log.Panic("Error al conectar a db desde adapter")
		return nil, err
	}

	db := newDbPQLDB(conn)

	return db, nil
}

func newDbPQLDBGergal(db *sql.DB) *DbdAdapter {
	return &DbdAdapter{db: db}
}

func connectToDbGergal(config *config.Config) (*sql.DB, error) {
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

func (d *DbdAdapter) GetCatalog(ctx context.Context) ([]domain.Product, error) {
	query := `
	SELECT
		p.id,
		p.name,
		p.description,
		p.category,

		pp.id,
		pp.name,
		pp.unit_type,
		pp.unit_value,
		pp.price,

		s.total_quantity - s.reserved_quantity AS stock
	FROM products p
	JOIN product_presentations pp ON pp.product_id = p.id
	JOIN stock s ON s.product_presentation_id = pp.id
	WHERE p.active = true AND pp.active = true
	ORDER BY p.id, pp.unit_value;
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productsMap := make(map[int]*domain.Product)

	for rows.Next() {
		var (
			pID int
			p   domain.Product
			pp  domain.ProductPresentation
		)

		err := rows.Scan(
			&pID,
			&p.Name,
			&p.Description,
			&p.Category,
			&pp.ID,
			&pp.Name,
			&pp.UnitType,
			&pp.UnitValue,
			&pp.Price,
			&pp.Stock,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := productsMap[pID]; !ok {
			p.ID = pID
			p.Presentations = []domain.ProductPresentation{}
			productsMap[pID] = &p
		}

		productsMap[pID].Presentations = append(productsMap[pID].Presentations, pp)
	}

	var products []domain.Product
	for _, p := range productsMap {
		products = append(products, *p)
	}

	return products, nil
}
func (d *DbdAdapter) GetDeliveryZones(ctx context.Context) ([]domain.DeliveryZone, error) {
	query := `
		SELECT id, name, price, estimated_time
		FROM delivery_zones
		WHERE active = true
		ORDER BY id;
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []domain.DeliveryZone

	for rows.Next() {
		var z domain.DeliveryZone
		if err := rows.Scan(&z.ID, &z.Name, &z.Price, &z.EstimatedTime); err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}

	return zones, nil
}
