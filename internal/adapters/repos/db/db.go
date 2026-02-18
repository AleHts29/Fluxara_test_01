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

func (d *DbdAdapter) GetCarrerasAll(ctx context.Context) ([]domain.CareerFull, error) {
	careers, err := d.getCareersBase(ctx)
	if err != nil {
		return nil, err
	}

	for i := range careers {
		subjects, err := d.getSubjectsBasicByCareer(ctx, careers[i].ID)
		if err != nil {
			return nil, err
		}

		careers[i].Materias = subjects
	}

	return careers, nil
}

// aux
func (d *DbdAdapter) getSubjectsBasicByCareer(ctx context.Context, careerID int) ([]domain.SubjectFull, error) {
	query := `
		SELECT
			s.id,
			s.name,
			s.description
		FROM subjects s
		JOIN career_subjects cs
			ON cs.subject_id = s.id
		WHERE cs.career_id = $1
		ORDER BY s.id;
	`

	rows, err := d.db.QueryContext(ctx, query, careerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []domain.SubjectFull

	for rows.Next() {
		var s domain.SubjectFull

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
		)
		if err != nil {
			return nil, err
		}

		// WeeklyHours, Slots, Horarios, Profesores quedan vacíos
		subjects = append(subjects, s)
	}

	return subjects, nil
}

func (d *DbdAdapter) getCareersBase(ctx context.Context) ([]domain.CareerFull, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.duration_years,

			sp.id,
			sp.name,

			spp.monthly_price,
			spp.enrollment_fee
		FROM careers c
		LEFT JOIN study_plans sp
			ON sp.career_id = c.id
		LEFT JOIN study_plan_prices spp
			ON spp.study_plan_id = sp.id
		ORDER BY c.id;
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var careers []domain.CareerFull

	for rows.Next() {
		var c domain.CareerFull
		var plan domain.StudyPlan

		var planID sql.NullInt64
		var planName sql.NullString
		var monthly sql.NullFloat64
		var enrollment sql.NullFloat64

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.DurationYears,

			&planID,
			&planName,

			&monthly,
			&enrollment,
		)
		if err != nil {
			return nil, err
		}

		if planID.Valid {
			plan.ID = int(planID.Int64)
			plan.Name = planName.String

			if monthly.Valid {
				plan.MonthlyPrice = monthly.Float64
			}
			if enrollment.Valid {
				plan.EnrollmentFee = enrollment.Float64
			}

			c.Plan = plan
		}

		careers = append(careers, c)
	}

	return careers, nil
}

func (d *DbdAdapter) getSubjectsByCareer(ctx context.Context, careerID int) ([]domain.SubjectFull, error) {
	query := `
		SELECT
			s.id,
			s.name,
			s.description,
			s.weekly_hours
		FROM subjects s
		JOIN career_subjects cs
			ON cs.subject_id = s.id
		WHERE cs.career_id = $1
		ORDER BY s.id;
	`

	rows, err := d.db.QueryContext(ctx, query, careerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []domain.SubjectFull

	for rows.Next() {
		var s domain.SubjectFull

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.WeeklyHours,
		)
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, s)
	}

	return subjects, nil
}

func (d *DbdAdapter) getSubjectSlots(ctx context.Context, subjectID int) (domain.SubjectSlots, error) {
	query := `
		SELECT
			total_slots,
			available_slots
		FROM subject_slots
		WHERE subject_id = $1;
	`

	var slots domain.SubjectSlots

	err := d.db.QueryRowContext(ctx, query, subjectID).Scan(
		&slots.Total,
		&slots.Available,
	)

	if err == sql.ErrNoRows {
		return slots, nil
	}
	if err != nil {
		return slots, err
	}

	return slots, nil
}

func (d *DbdAdapter) getSubjectSchedules(ctx context.Context, subjectID int) ([]domain.SubjectSchedule, error) {
	query := `
		SELECT
			day_of_week,
			start_time,
			end_time,
			modality
		FROM subject_schedules
		WHERE subject_id = $1
		ORDER BY day_of_week, start_time;
	`

	rows, err := d.db.QueryContext(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.SubjectSchedule

	for rows.Next() {
		var s domain.SubjectSchedule

		err := rows.Scan(
			&s.DayOfWeek,
			&s.StartTime,
			&s.EndTime,
			&s.Modality,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, s)
	}

	return schedules, nil
}

func (d *DbdAdapter) getSubjectProfesors(ctx context.Context, subjectID int) ([]domain.Professor, error) {
	query := `
		SELECT
			p.id,
			p.full_name,
			p.email
		FROM professors p
		JOIN subject_professors sp
			ON sp.professor_id = p.id
		WHERE sp.subject_id = $1
		ORDER BY full_name;
	`

	rows, err := d.db.QueryContext(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profs []domain.Professor

	for rows.Next() {
		var p domain.Professor

		err := rows.Scan(
			&p.ID,
			&p.FullName,
			&p.Email,
		)
		if err != nil {
			return nil, err
		}

		profs = append(profs, p)
	}

	return profs, nil
}
