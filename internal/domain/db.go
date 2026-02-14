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

// arte
type Career struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// /carreras/resumen
type CareersResumen struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Materias    []SubjectResumen `json:"subjects"`
}

type SubjectResumen struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Profesores  []ProfessorResumen `json:"professors"`
}

type ProfessorResumen struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
