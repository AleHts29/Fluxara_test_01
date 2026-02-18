package domain

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

// abm
type CareerFull struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	DurationYears int           `json:"duration_years"`
	Plan          StudyPlan     `json:"plan"`
	Materias      []SubjectFull `json:"materias"`
}

type StudyPlan struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	MonthlyPrice  float64 `json:"monthly_price"`
	EnrollmentFee float64 `json:"enrollment_fee"`
}

type SubjectFull struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	WeeklyHours int               `json:"weekly_hours,omitempty"`
	Slots       *SubjectSlots     `json:"slots,omitempty"`
	Horarios    []SubjectSchedule `json:"horarios,omitempty"`
	Profesores  []Professor       `json:"profesores,omitempty"`
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

// gergal
type Product struct {
	ID            int                   `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Category      string                `json:"category"`
	Presentations []ProductPresentation `json:"presentations,omitempty"`
}

type ProductPresentation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	UnitType  string  `json:"unit_type"`
	UnitValue float64 `json:"unit_value"`
	Price     float64 `json:"price"`
	Stock     float64 `json:"stock,omitempty"`
}

type DeliveryZone struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	EstimatedTime string  `json:"estimated_time"`
}

type Order struct {
	ID           int         `json:"id"`
	CustomerID   int         `json:"customer_id"`
	AddressID    int         `json:"address_id"`
	Status       string      `json:"status"`
	Subtotal     float64     `json:"subtotal"`
	DeliveryCost float64     `json:"delivery_cost"`
	Total        float64     `json:"total"`
	Items        []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductPresentationID int     `json:"product_presentation_id"`
	Name                  string  `json:"name"`
	Quantity              float64 `json:"quantity"`
	UnitPrice             float64 `json:"unit_price"`
	Subtotal              float64 `json:"subtotal"`
}
