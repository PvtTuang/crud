package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
