package history

import (
	"database/sql"

	"github.com/google/uuid"
)

type HistoryRepository interface {
	GetByUser(userID uuid.UUID) ([]PurchaseHistory, error)
	Create(history *PurchaseHistory) error
}

type historyRepository struct {
	db *sql.DB
}

// Create implements HistoryRepository.
func (h *historyRepository) Create(history *PurchaseHistory) error {
	panic("unimplemented")
}

// GetByUser implements HistoryRepository.
func (h *historyRepository) GetByUser(userID uuid.UUID) ([]PurchaseHistory, error) {
	panic("unimplemented")
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

// func (r *historyRepository) GetByUser(userID uuid.UUID) ([]PurchaseHistory, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			_ = tx.Rollback()
// 		}
// 	}()

// 	rows, err := tx.Query(`SELECT h.id, ph.user_id, ph.created_at,pi.id, pi.product_id, pi.quantity,p.name, p.price
// 	FROM purchase_histories ph
// 	LEFT JOIN purchase_item pi ON ph.id = pi.purchase_id
// 	LEFT JOIN products p ON pi.product_id = p.id
// 	WHERE ph.user_id = $ 1
// 	ORDER BY ph.created_at DECS;
// 	`)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// }

// func (r *historyRepository) Create(history *PurchaseHistory) error {
// 	return r.db.Create(history).Error
// }
