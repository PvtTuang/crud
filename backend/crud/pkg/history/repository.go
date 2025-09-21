package history

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HistoryRepository interface {
	GetByUser(userID uuid.UUID) ([]PurchaseHistory, error)
	Create(history *PurchaseHistory) error
}

type historyRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) GetByUser(userID uuid.UUID) ([]PurchaseHistory, error) {
	var histories []PurchaseHistory
	if err := r.db.
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *historyRepository) Create(history *PurchaseHistory) error {
	return r.db.Create(history).Error
}
