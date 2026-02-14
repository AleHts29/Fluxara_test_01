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
func (dPQLDB *DbdAdapter) GetFullData(ctx context.Context) ([]domain.CareerFull, error) {

	query := `
	SELECT
	    c.id, c.name, c.description, c.duration_years,
	    sp.id, sp.name,
	    spp.monthly_price, spp.enrollment_fee,
	    s.id, s.name, s.description, s.weekly_hours,
	    sch.day_of_week, sch.start_time, sch.end_time, sch.modality,
	    sl.total_slots, sl.available_slots,
	    p.id, p.full_name, p.email
	FROM careers c
	JOIN study_plans sp ON sp.career_id = c.id
	LEFT JOIN study_plan_prices spp ON spp.study_plan_id = sp.id
	JOIN career_subjects cs ON cs.career_id = c.id
	JOIN subjects s ON s.id = cs.subject_id
	LEFT JOIN subject_schedules sch ON sch.subject_id = s.id
	LEFT JOIN subject_slots sl ON sl.subject_id = s.id
	LEFT JOIN subject_professors spf ON spf.subject_id = s.id
	LEFT JOIN professors p ON p.id = spf.professor_id
	ORDER BY c.id, s.id;
	`

	rows, err := dPQLDB.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	careerMap := make(map[int]*domain.CareerFull)

	for rows.Next() {

		var (
			careerID, duration     int
			careerName, careerDesc string

			planID              int
			planName            string
			monthly, enrollment sql.NullFloat64

			subID, weeklyHours int
			subName, subDesc   string

			day, modality sql.NullString
			start, end    sql.NullTime

			totalSlots, availSlots sql.NullInt64

			profID              sql.NullInt64
			profName, profEmail sql.NullString
		)

		err := rows.Scan(
			&careerID, &careerName, &careerDesc, &duration,
			&planID, &planName, &monthly, &enrollment,
			&subID, &subName, &subDesc, &weeklyHours,
			&day, &start, &end, &modality,
			&totalSlots, &availSlots,
			&profID, &profName, &profEmail,
		)
		if err != nil {
			return nil, err
		}

		carrera, ok := careerMap[careerID]
		if !ok {
			carrera = &domain.CareerFull{
				ID:            careerID,
				Name:          careerName,
				Description:   careerDesc,
				DurationYears: duration,
				Plan: domain.StudyPlan{
					ID:            planID,
					Name:          planName,
					MonthlyPrice:  monthly.Float64,
					EnrollmentFee: enrollment.Float64,
				},
			}
			careerMap[careerID] = carrera
		}

		var materia *domain.SubjectFull
		for i := range carrera.Materias {
			if carrera.Materias[i].ID == subID {
				materia = &carrera.Materias[i]
				break
			}
		}

		if materia == nil {
			materia = &domain.SubjectFull{
				ID:          subID,
				Name:        subName,
				Description: subDesc,
				WeeklyHours: weeklyHours,
				Slots: domain.SubjectSlots{
					Total:     int(totalSlots.Int64),
					Available: int(availSlots.Int64),
				},
			}
			carrera.Materias = append(carrera.Materias, *materia)
			materia = &carrera.Materias[len(carrera.Materias)-1]
		}

		if day.Valid {
			materia.Horarios = append(materia.Horarios, domain.SubjectSchedule{
				DayOfWeek: day.String,
				StartTime: start.Time.Format("15:04"),
				EndTime:   end.Time.Format("15:04"),
				Modality:  modality.String,
			})
		}

		if profID.Valid {
			materia.Profesores = append(materia.Profesores, domain.Professor{
				ID:       int(profID.Int64),
				FullName: profName.String,
				Email:    profEmail.String,
			})
		}
	}

	result := make([]domain.CareerFull, 0, len(careerMap))
	for _, c := range careerMap {
		result = append(result, *c)
	}

	return result, nil
}
func (dPQLDB *DbdAdapter) GetCarrerasAll(ctx context.Context) ([]domain.CareerFull, error) {
	var carrers []domain.CareerFull

	query := `
		SELECT * FROM public.careers
	`

	rows, err := dPQLDB.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var career domain.CareerFull

		err := rows.Scan(
			&career.ID,
			&career.Name,
			&career.Description,
			&career.Plan,
			&career.Materias,
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

// func (dPQLDB *DbdAdapter) GetCarrerasResumen(ctx context.Context) ([]domain.CareerFull, error) {
// 	query := `
// 		SELECT
// 			c.id, c.name, c.description, c.duration_years,

// 			sp.id, sp.name,
// 			spp.monthly_price, spp.enrollment_fee,

// 			s.id, s.name, s.description, s.weekly_hours,

// 			sch.day_of_week, sch.start_time, sch.end_time, sch.modality,

// 			sl.total_slots, sl.available_slots,

// 			p.id, p.full_name, p.email

// 		FROM careers c
// 		JOIN study_plans sp ON sp.career_id = c.id
// 		LEFT JOIN study_plan_prices spp ON spp.study_plan_id = sp.id

// 		JOIN career_subjects cs ON cs.career_id = c.id
// 		JOIN subjects s ON s.id = cs.subject_id

// 		LEFT JOIN subject_schedules sch ON sch.subject_id = s.id
// 		LEFT JOIN subject_slots sl ON sl.subject_id = s.id
// 		LEFT JOIN subject_professors spf ON spf.subject_id = s.id
// 		LEFT JOIN professors p ON p.id = spf.professor_id

// 		ORDER BY c.id, s.id;

// 	`

// 	rows, err := dPQLDB.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	carrerasMap := make(map[int]*domain.CareerFull)

// 	for rows.Next() {

// 		var (
// 			careerID               int
// 			careerName, careerDesc string

// 			planID, planDuration int
// 			planName             string
// 			planCost             float64

// 			subID, weeklyHours int
// 			subName, subDesc   string

// 			day, modality sql.NullString
// 			start, end    sql.NullTime

// 			monthly, enrollment        sql.NullFloat64
// 			totalSlots, availableSlots sql.NullInt64

// 			profID              sql.NullInt64
// 			profName, profEmail sql.NullString
// 		)

// 		err := rows.Scan(
// 			&careerID, &careerName, &careerDesc,
// 			&planID, &planName, &planDuration, &planCost,
// 			&subID, &subName, &subDesc, &weeklyHours,
// 			&day, &start, &end, &modality,
// 			&monthly, &enrollment,
// 			&totalSlots, &availableSlots,
// 			&profID, &profName, &profEmail,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// -------------------------------
// 		// Carrera
// 		// -------------------------------
// 		carrera, ok := carrerasMap[careerID]
// 		if !ok {
// 			carrera = &domain.CareerFull{
// 				ID:          careerID,
// 				Name:        careerName,
// 				Description: careerDesc,
// 				Plan: domain.StudyPlan{
// 					ID:            planID,
// 					Name:          planName,
// 					DurationYears: planDuration,
// 					TotalCost:     planCost,
// 				},
// 			}
// 			carrerasMap[careerID] = carrera
// 		}

// 		// -------------------------------
// 		// Materia
// 		// -------------------------------
// 		var materia *domain.SubjectFull
// 		for i := range carrera.Materias {
// 			if carrera.Materias[i].ID == subID {
// 				materia = &carrera.Materias[i]
// 				break
// 			}
// 		}

// 		if materia == nil {
// 			materia = &domain.SubjectFull{
// 				ID:          subID,
// 				Name:        subName,
// 				Description: subDesc,
// 				WeeklyHours: weeklyHours,
// 				Prices: domain.SubjectPrice{
// 					Monthly:    monthly.Float64,
// 					Enrollment: enrollment.Float64,
// 				},
// 				Slots: domain.SubjectSlots{
// 					Total:     int(totalSlots.Int64),
// 					Available: int(availableSlots.Int64),
// 				},
// 			}
// 			carrera.Materias = append(carrera.Materias, *materia)
// 			materia = &carrera.Materias[len(carrera.Materias)-1]
// 		}

// 		// -------------------------------
// 		// Horarios
// 		// -------------------------------
// 		if day.Valid {
// 			materia.Horarios = append(materia.Horarios, domain.SubjectSchedule{
// 				DayOfWeek: day.String,
// 				StartTime: start.Time.Format("15:04"),
// 				EndTime:   end.Time.Format("15:04"),
// 				Modality:  modality.String,
// 			})
// 		}

// 		// -------------------------------
// 		// Profesores
// 		// -------------------------------
// 		if profID.Valid {
// 			materia.Profesores = append(materia.Profesores, domain.Professor{
// 				ID:       int(profID.Int64),
// 				FullName: profName.String,
// 				Email:    profEmail.String,
// 			})
// 		}
// 	}

// 	result := make([]domain.CareerFull, 0, len(carrerasMap))
// 	for _, c := range carrerasMap {
// 		result = append(result, *c)
// 	}

// 	return result, nil
// }
