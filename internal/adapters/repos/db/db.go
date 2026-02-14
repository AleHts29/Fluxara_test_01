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

// func (dPQLDB *DbdAdapter) GetProductsAll(ctx context.Context) ([]domain.Product, error) {
// 	var products []domain.Product

// 	query := `
// 		SELECT
// 			id,
// 			sku,
// 			name,
// 			category,
// 			price_cents,
// 			stock,
// 			is_active,
// 			created_at
// 		FROM public.products
// 		WHERE is_active = true
// 	`

// 	rows, err := dPQLDB.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var product domain.Product

// 		err := rows.Scan(
// 			&product.ID,
// 			&product.SKU,
// 			&product.Name,
// 			&product.Category,
// 			&product.PriceCents,
// 			&product.Stock,
// 			&product.IsActive,
// 			&product.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		products = append(products, product)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return products, nil
// }
// func (dPQLDB *DbdAdapter) GetProduct(ctx context.Context, id string) (domain.Product, error) {
// 	var product domain.Product

// 	query := `
// 		SELECT
// 			id,
// 			sku,
// 			name,
// 			category,
// 			price_cents,
// 			stock,
// 			is_active,
// 			created_at
// 		FROM public.products
// 		WHERE id = $1
// 		AND is_active = true
// 	`

// 	ctx = context.Background()
// 	err := dPQLDB.db.QueryRowContext(ctx, query, id).Scan(
// 		&product.ID,
// 		&product.SKU,
// 		&product.Name,
// 		&product.Category,
// 		&product.PriceCents,
// 		&product.Stock,
// 		&product.IsActive,
// 		&product.CreatedAt,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return product, fmt.Errorf("producto no encontrado")
// 		}
// 		return product, err
// 	}

// 	return product, nil
// }

// arte
func (dPQLDB *DbdAdapter) GetCarrerasAll(ctx context.Context) ([]domain.Career, error) {
	var carrers []domain.Career

	query := `
		SELECT * FROM public.careers
	`

	rows, err := dPQLDB.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var career domain.Career

		err := rows.Scan(
			&career.ID,
			&career.Name,
			&career.Description,
		)
		if err != nil {
			return nil, err
		}

		carrers = append(carrers, career)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carrers, nil
}
func (dPQLDB *DbdAdapter) GetCarrerasResumen(ctx context.Context) ([]domain.CareersResumen, error) {
	query := `
		SELECT
            c.id,
            c.name,
            c.description,
            s.id,
            s.name,
            s.description,
            p.id,
            p.full_name,
            p.email
        FROM careers c
        JOIN career_subjects cs ON cs.career_id = c.id
        JOIN subjects s ON s.id = cs.subject_id
        LEFT JOIN subject_professors sp ON sp.subject_id = s.id
        LEFT JOIN professors p ON p.id = sp.professor_id
        ORDER BY c.id, s.id, p.id
	`

	rows, err := dPQLDB.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	carrerasMap := make(map[int]*domain.CareersResumen)

	for rows.Next() {
		var (
			careerID                      int
			careerName, careerDesc        string
			subjectID                     int
			subjectName, subjectDesc      string
			professorID                   sql.NullInt64
			professorName, professorEmail sql.NullString
		)

		if err := rows.Scan(
			&careerID, &careerName, &careerDesc,
			&subjectID, &subjectName, &subjectDesc,
			&professorID, &professorName, &professorEmail,
		); err != nil {
			return nil, err
		}

		// Carrera
		carrera, ok := carrerasMap[careerID]
		if !ok {
			carrera = &domain.CareersResumen{
				ID:          careerID,
				Name:        careerName,
				Description: careerDesc,
			}
			carrerasMap[careerID] = carrera
		}

		// Materia
		var materia *domain.SubjectResumen
		for i := range carrera.Materias {
			if carrera.Materias[i].ID == subjectID {
				materia = &carrera.Materias[i]
				break
			}
		}

		if materia == nil {
			materia = &domain.SubjectResumen{
				ID:          subjectID,
				Name:        subjectName,
				Description: subjectDesc,
			}
			carrera.Materias = append(carrera.Materias, *materia)
			materia = &carrera.Materias[len(carrera.Materias)-1]
		}

		// Profesor
		if professorID.Valid {
			materia.Profesores = append(materia.Profesores, domain.ProfessorResumen{
				ID:       int(professorID.Int64),
				FullName: professorName.String,
				Email:    professorEmail.String,
			})
		}
	}

	result := make([]domain.CareersResumen, 0, len(carrerasMap))
	for _, c := range carrerasMap {
		result = append(result, *c)
	}

	return result, nil
}
func (dPQLDB *DbdAdapter) GetCarrerasByName(ctx context.Context, name string) (domain.Career, error) {
	var career domain.Career

	query := `
		SELECT *
		FROM public.careers
		WHERE name = $1
	`

	ctx = context.Background()
	err := dPQLDB.db.QueryRowContext(ctx, query, name).Scan(
		&career.ID,
		&career.Name,
		&career.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return career, fmt.Errorf("carrera no encontrado")
		}
		return career, err
	}

	return career, nil
}
