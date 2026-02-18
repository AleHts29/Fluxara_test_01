package db

import (
	"context"
	"database/sql"
	"errors"
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

// get
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

// post
func (d *DbdAdapter) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	fmt.Printf("Create Order -.------------- \n")
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var subtotal float64
	var items []domain.OrderItem

	for _, it := range req.Items {
		var name string
		var price float64
		fmt.Printf("For Item -.------------- \n")
		err := tx.QueryRowContext(ctx, `
			SELECT pp.name, pp.price
			FROM product_presentations pp
			WHERE pp.id = $1 AND pp.active = true
		`, it.ProductPresentationID).Scan(&name, &price)
		if err != nil {
			return nil, err
		}

		if err := d.reserveStock(ctx, tx, it.ProductPresentationID, it.Quantity); err != nil {
			return nil, err
		}

		sub := price * it.Quantity
		subtotal += sub

		items = append(items, domain.OrderItem{
			ProductPresentationID: it.ProductPresentationID,
			Name:                  name,
			Quantity:              it.Quantity,
			UnitPrice:             price,
			Subtotal:              sub,
		})
	}

	delivery, err := d.getDeliveryCost(ctx, tx, req.AddressID)
	if err != nil {
		return nil, err
	}

	total := subtotal + delivery

	fmt.Printf("subtotal + deliverys -.------------- \n")

	var orderID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (customer_id, address_id, status, subtotal, delivery_cost, total)
		VALUES ($1,$2,'payment_pending',$3,$4,$5)
		RETURNING id
	`, req.CustomerID, req.AddressID, subtotal, delivery, total).
		Scan(&orderID)
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_items
			(order_id, product_presentation_id, quantity, unit_price, subtotal)
			VALUES ($1,$2,$3,$4,$5)
		`, orderID, it.ProductPresentationID, it.Quantity, it.UnitPrice, it.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:           orderID,
		CustomerID:   req.CustomerID,
		AddressID:    req.AddressID,
		Status:       "payment_pending",
		Subtotal:     subtotal,
		DeliveryCost: delivery,
		Total:        total,
		Items:        items,
	}, nil
}

// aux
func (d *DbdAdapter) reserveStock(ctx context.Context, tx *sql.Tx, presID int, qty float64) error {
	res, err := tx.ExecContext(ctx, `
		UPDATE stock
		SET reserved_quantity = reserved_quantity + $1,
		    updated_at = now()
		WHERE product_presentation_id = $2
		  AND total_quantity - reserved_quantity >= $1
	`, qty, presID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("stock insuficiente")
	}
	return nil
}

func (d *DbdAdapter) getDeliveryCost(ctx context.Context, tx *sql.Tx, addressID int) (float64, error) {
	var cost float64
	fmt.Printf("getDeliveryCost -.------------- \n")
	err := tx.QueryRowContext(ctx, `
		SELECT dz.price
		FROM customer_addresses ca
		INNER JOIN delivery_zones dz ON dz.id = ca.zone_id
		WHERE ca.id = $1`, addressID).Scan(&cost)

	return cost, err
}

func (d *DbdAdapter) MarkOrderPaid(ctx context.Context, orderID int) error {
	_, err := d.db.ExecContext(ctx, `
		UPDATE orders
		SET status = 'paid'
		WHERE id = $1
	`, orderID)
	return err
}
