package domain

type MercadoPago struct {
	Token string
}

type CreateOrderRequest struct {
	CustomerID int `json:"customer_id"`
	AddressID  int `json:"address_id"`
	Items      []struct {
		ProductPresentationID int     `json:"product_presentation_id"`
		Quantity              float64 `json:"quantity"`
	} `json:"items"`
}

type OrderPreviewRequest struct {
	AddressID int `json:"address_id"`
	Items     []struct {
		ProductPresentationID int     `json:"product_presentation_id"`
		Quantity              float64 `json:"quantity"`
	} `json:"items"`
}

type CreateOrderResponse struct {
	Order   *Order       `json:"order"`
	Payment *PaymentLink `json:"payment"`
}

type PaymentLink struct {
	OrderID int    `json:"order_id"`
	URL     string `json:"url"`
}
