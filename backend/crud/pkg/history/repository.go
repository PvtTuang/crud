package history

import (
	"crud/pkg/product"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type HistoryRepository interface {
	GetByUser(userID uuid.UUID) ([]PurchaseHistory, error)
	Create(history *PurchaseHistory) error
}

type historyRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

// Create PurchaseHistory + Items
func (h *historyRepository) Create(history *PurchaseHistory) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// insert history
	err = tx.QueryRow(`
		INSERT INTO purchase_histories (id, user_id, created_at)
		VALUES (gen_random_uuid(), $1, NOW())
		RETURNING id, created_at;
	`, history.UserID).Scan(&history.ID, &history.CreatedAt)
	if err != nil {
		return err
	}

	// insert items
	for i := range history.Items {
		item := &history.Items[i]
		_, err = tx.Exec(`
			INSERT INTO purchase_items (purchase_id, product_id, quantity)
			VALUES ($1, $2, $3);
		`, history.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
		item.PurchaseID = history.ID
	}

	return tx.Commit()
}

// ดึงประวัติการซื้อของ user
func (h *historyRepository) GetByUser(userID uuid.UUID) ([]PurchaseHistory, error) {
	rows, err := h.db.Query(`
		SELECT ph.id, ph.user_id, ph.created_at,
		       pi.id, pi.product_id, pi.quantity,
		       p.id, p.name, p.price, p.created_at, p.updated_at, p.deleted_at
		FROM purchase_histories ph
		LEFT JOIN purchase_items pi ON ph.id = pi.purchase_id
		LEFT JOIN products p ON pi.product_id = p.id
		WHERE ph.user_id = $1
		ORDER BY ph.created_at DESC;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	historiesMap := make(map[uuid.UUID]*PurchaseHistory)

	for rows.Next() {
		var (
			hID       uuid.UUID
			uID       uuid.UUID
			createdAt time.Time
			itemID    sql.NullInt64
			productID sql.NullString
			quantity  sql.NullInt32
			prod      product.Product
		)

		err = rows.Scan(&hID, &uID, &createdAt,
			&itemID, &productID, &quantity,
			&prod.ID, &prod.Name, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt, &prod.DeletedAt)
		if err != nil {
			return nil, err
		}

		// ถ้ายังไม่มี history ใน map → สร้างใหม่
		if _, ok := historiesMap[hID]; !ok {
			historiesMap[hID] = &PurchaseHistory{
				ID:        hID,
				UserID:    uID,
				CreatedAt: createdAt,
				Items:     []PurchaseItem{},
			}
		}

		// ถ้ามี item จริง → append
		if itemID.Valid {
			historiesMap[hID].Items = append(historiesMap[hID].Items, PurchaseItem{
				ID:         uint(itemID.Int64),
				PurchaseID: hID,
				ProductID:  uuid.MustParse(productID.String),
				Quantity:   int(quantity.Int32),
				Product:    prod,
			})
		}
	}

	// แปลง map → slice
	var histories []PurchaseHistory
	for _, h := range historiesMap {
		histories = append(histories, *h)
	}

	return histories, nil
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
