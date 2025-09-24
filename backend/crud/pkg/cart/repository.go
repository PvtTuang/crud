package cart

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CartRepository interface {
	GetCart(userID uuid.UUID) (*Cart, error)
	AddProduct(userID, productID uuid.UUID, qty int) error
	RemoveProduct(userID, productID uuid.UUID) error
	ClearCart(userID uuid.UUID) error
	Checkout(userID uuid.UUID) (err error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCart(userID uuid.UUID) (*Cart, error) {
	var cart Cart

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRow(`
		SELECT id, created_at, updated_at, user_id
		FROM carts
		WHERE user_id = $1 
		LIMIT 1;`, userID)

	err = row.Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow(`
				INSERT INTO carts (id, user_id, created_at, updated_at)
				VALUES (gen_random_uuid(), $1, NOW(), NOW())
				RETURNING id, created_at, updated_at, user_id;
			`, userID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt, &cart.UserID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	rows, err := tx.Query(`
		SELECT cp.id, cp.cart_id, cp.product_id, cp.quantity,
		       p.id, p.name, p.price, p.created_at, p.updated_at, p.deleted_at
		FROM cart_products cp
		JOIN products p ON cp.product_id = p.id
		WHERE cp.cart_id = $1;`, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cp CartProduct
		err = rows.Scan(
			&cp.ID, &cp.CartID, &cp.ProductID, &cp.Quantity,
			&cp.Product.ID, &cp.Product.Name, &cp.Product.Price,
			&cp.Product.CreatedAt, &cp.Product.UpdatedAt, &cp.Product.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		cart.Products = append(cart.Products, cp)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) AddProduct(userID uuid.UUID, productID uuid.UUID, qty int) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var cartID uuid.UUID
	err = tx.QueryRow(`
        SELECT id FROM carts WHERE user_id = $1 LIMIT 1
    `, userID).Scan(&cartID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`
            INSERT INTO carts (id, user_id, created_at, updated_at)
            VALUES (gen_random_uuid(), $1, NOW(), NOW())
            RETURNING id
        `, userID).Scan(&cartID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var existingQty int
	err = tx.QueryRow(`
        SELECT quantity FROM cart_products WHERE cart_id = $1 AND product_id = $2
    `, cartID, productID).Scan(&existingQty)

	switch err {
	case sql.ErrNoRows:
		_, err = tx.Exec(`
            INSERT INTO cart_products (cart_id, product_id, quantity)
            VALUES ($1, $2, $3)
        `, cartID, productID, qty)
	case nil:
		_, err = tx.Exec(`
            UPDATE cart_products
            SET quantity = quantity + $3
            WHERE cart_id = $1 AND product_id = $2
        `, cartID, productID, qty)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *cartRepository) Checkout(userID uuid.UUID) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var cartID uuid.UUID
	err = tx.QueryRow(`
		SELECT id FROM carts WHERE user_id = $1 LIMIT 1;
	`, userID).Scan(&cartID)
	if err != nil {
		return fmt.Errorf("ไม่พบ cart ของ user %v: %w", userID, err)
	}

	var purchaseID uuid.UUID
	err = tx.QueryRow(`
		INSERT INTO purchase_histories (id, user_id, created_at)
		VALUES (gen_random_uuid(), $1, NOW())
		RETURNING id;
	`, userID).Scan(&purchaseID)
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
		INSERT INTO purchase_items (purchase_id, product_id, quantity)
		SELECT $1, product_id, quantity
		FROM cart_products
		WHERE cart_id = $2;
	`, purchaseID, cartID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cart ว่าง ไม่สามารถ checkout ได้")
	}

	_, err = tx.Exec(`DELETE FROM cart_products WHERE cart_id = $1;`, cartID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *cartRepository) ClearCart(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_products
		WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1 AND deleted_at IS NULL LIMIT 1);`, userID)
	return err
}

func (r *cartRepository) RemoveProduct(userID uuid.UUID, productID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		DELETE FROM cart_products
		WHERE product_id = $1
		  AND cart_id = (SELECT id FROM carts WHERE user_id = $2 LIMIT 1);`,
		productID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// func (r *cartRepository) GetCart(userID uuid.UUID) (*Cart, error) {
// 	var cart Cart
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	row := tx.QueryRow(``)
// }

// func (r *cartRepository) AddProduct(userID, productID uuid.UUID, qty int) error {
// 	t, err := transaction.Begin(r.db)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			t.Rollback()
// 		}
// 	}()

// 	var cart Cart
// 	if err = t.DB().Where("user_id = ?", userID).FirstOrCreate(&cart, Cart{UserID: userID}).Error; err != nil {
// 		return err
// 	}

// 	var cp CartProduct
// 	if err = t.DB().Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cp).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			cp = CartProduct{CartID: cart.ID, ProductID: productID, Quantity: qty}
// 			if err = t.DB().Create(&cp).Error; err != nil {
// 				return err
// 			}
// 		} else {
// 			return err
// 		}
// 	} else {
// 		cp.Quantity += qty
// 		if err = t.DB().Save(&cp).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return t.Commit()
// }

// func (r *cartRepository) RemoveProduct(userID, productID uuid.UUID) error {
// 	t, err := transaction.Begin(r.db)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			t.Rollback()
// 		}
// 	}()

// 	var cart Cart
// 	if err = t.DB().Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 		return err
// 	}

// 	if err = t.DB().Where("cart_id = ? AND product_id = ?", cart.ID, productID).
// 		Delete(&CartProduct{}).Error; err != nil {
// 		return err
// 	}

// 	return t.Commit()
// }

// func (r *cartRepository) ClearCart(userID uuid.UUID) (err error) {
// 	t, err := transaction.Begin(r.db)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			t.Rollback()
// 		}
// 	}()

// 	var cart Cart
// 	if err = t.DB().Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 		return err
// 	}

// 	if err = t.DB().Where("cart_id = ?", cart.ID).Delete(&CartProduct{}).Error; err != nil {
// 		return err
// 	}

// 	if err = t.DB().Delete(&cart).Error; err != nil {
// 		return err
// 	}

// 	return t.Commit()
// }

// func (r *cartRepository) Checkout(userID uuid.UUID) (err error) {
// 	t, err := transaction.Begin(r.db)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			t.Rollback()
// 		}
// 	}()

// 	var cart Cart
// 	if err = t.DB().Preload("Products").Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 		return err
// 	}

// 	// สร้าง PurchaseHistory
// 	purchase := history.PurchaseHistory{
// 		UserID: userID,
// 		Items:  []history.PurchaseItem{},
// 	}

// 	for _, cp := range cart.Products {
// 		purchase.Items = append(purchase.Items, history.PurchaseItem{
// 			ProductID: cp.ProductID,
// 			Quantity:  cp.Quantity,
// 		})
// 	}

// 	if err = t.DB().Create(&purchase).Error; err != nil {
// 		return err
// 	}

// 	// เคลียร์ cart
// 	if err = t.DB().Where("cart_id = ?", cart.ID).Delete(&CartProduct{}).Error; err != nil {
// 		return err
// 	}
// 	if err = t.DB().Delete(&cart).Error; err != nil {
// 		return err
// 	}

// 	return t.Commit()
// }
