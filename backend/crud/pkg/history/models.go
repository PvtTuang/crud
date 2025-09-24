package history

import (
	"crud/pkg/product"
	"time"

	"github.com/google/uuid"
)

type PurchaseHistory struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	Items     []PurchaseItem `json:"items"`
}

type PurchaseItem struct {
	ID         uint            `json:"id"`
	PurchaseID uuid.UUID       `json:"purchase_id"`
	ProductID  uuid.UUID       `json:"product_id"`
	Quantity   int             `json:"quantity"`
	Product    product.Product `json:"product"`
}
