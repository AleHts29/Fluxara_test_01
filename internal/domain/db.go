package domain

import (
	"time"
)

type Db struct {
	Connection  string `mapstructure:"connection"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SslMode     string `mapstructure:"sslmode"`
	Retries     int    `mapstructure:"retries"`
	TimeRetries int    `mapstructure:"time_retries"`
}

type Product struct {
	ID         string
	SKU        string
	Name       string
	Category   string
	PriceCents int32
	Stock      int32
	IsActive   bool
	CreatedAt  time.Time
}

// // arte
// type Career struct {
// 	ID          int    `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// }

// // /carreras/resumen
// type CareersResumen struct {
// 	ID          int              `json:"id"`
// 	Name        string           `json:"name"`
// 	Description string           `json:"description"`
// 	Materias    []SubjectResumen `json:"subjects"`
// }

// type SubjectResumen struct {
// 	ID          int                `json:"id"`
// 	Name        string             `json:"name"`
// 	Description string             `json:"description"`
// 	Profesores  []ProfessorResumen `json:"professors"`
// }

// type ProfessorResumen struct {
// 	ID       int    `json:"id"`
// 	FullName string `json:"full_name"`
// 	Email    string `json:"email"`
// }

// arte
type CareerFull struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Plan        StudyPlan     `json:"plan"`
	Materias    []SubjectFull `json:"materias"`
}

type StudyPlan struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	DurationYears int     `json:"duration_years"`
	TotalCost     float64 `json:"total_cost"`
}

type SubjectFull struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	WeeklyHours int               `json:"weekly_hours"`
	Prices      SubjectPrice      `json:"prices"`
	Slots       SubjectSlots      `json:"slots"`
	Horarios    []SubjectSchedule `json:"horarios"`
	Profesores  []Professor       `json:"profesores"`
}

type SubjectSchedule struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Modality  string `json:"modality"`
}

type SubjectPrice struct {
	Monthly    float64 `json:"monthly"`
	Enrollment float64 `json:"enrollment"`
}

type SubjectSlots struct {
	Total     int `json:"total"`
	Available int `json:"available"`
}

type Professor struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
