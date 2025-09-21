package cart

import (
	"crud/pkg/product"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Products  []CartProduct  `gorm:"foreignKey:CartID" json:"products"`
}

type CartProduct struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	CartID    uuid.UUID       `gorm:"type:uuid;not null" json:"cart_id"`
	ProductID uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	Product   product.Product `gorm:"foreignKey:ProductID" json:"product"`
}
