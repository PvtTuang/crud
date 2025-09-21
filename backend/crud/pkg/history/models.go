package history

import (
	"crud/pkg/product"
	"time"

	"github.com/google/uuid"
)

type PurchaseHistory struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	Items     []PurchaseItem `gorm:"foreignKey:PurchaseID" json:"items"`
}

type PurchaseItem struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	PurchaseID uuid.UUID       `gorm:"type:uuid;not null" json:"purchase_id"`
	ProductID  uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Quantity   int             `json:"quantity"`
	Product    product.Product `gorm:"foreignKey:ProductID" json:"product"`
}
