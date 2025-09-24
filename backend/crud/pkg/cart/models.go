package cart

import (
	"crud/pkg/product"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	UserID    uuid.UUID     `json:"user_id"`
	Products  []CartProduct `json:"products"`
}

type CartProduct struct {
	ID        uint            `json:"id"`
	CartID    uuid.UUID       `json:"cart_id"`
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Product   product.Product `json:"product"`
}
