package cart

import (
	"crud/pkg/history"
	"crud/pkg/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(userID uuid.UUID) (*Cart, error)
	AddProduct(userID, productID uuid.UUID, qty int) error
	RemoveProduct(userID, productID uuid.UUID) error
	ClearCart(userID uuid.UUID) error
	Checkout(userID uuid.UUID) (err error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCart(userID uuid.UUID) (*Cart, error) {
	var cart Cart
	if err := r.db.
		Preload("Products.Product").
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) AddProduct(userID, productID uuid.UUID, qty int) error {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	var cart Cart
	if err = t.DB().Where("user_id = ?", userID).FirstOrCreate(&cart, Cart{UserID: userID}).Error; err != nil {
		return err
	}

	var cp CartProduct
	if err = t.DB().Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cp = CartProduct{CartID: cart.ID, ProductID: productID, Quantity: qty}
			if err = t.DB().Create(&cp).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		cp.Quantity += qty
		if err = t.DB().Save(&cp).Error; err != nil {
			return err
		}
	}
	return t.Commit()
}

func (r *cartRepository) RemoveProduct(userID, productID uuid.UUID) error {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	var cart Cart
	if err = t.DB().Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	if err = t.DB().Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		Delete(&CartProduct{}).Error; err != nil {
		return err
	}

	return t.Commit()
}

func (r *cartRepository) ClearCart(userID uuid.UUID) (err error) {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	var cart Cart
	if err = t.DB().Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	if err = t.DB().Where("cart_id = ?", cart.ID).Delete(&CartProduct{}).Error; err != nil {
		return err
	}

	if err = t.DB().Delete(&cart).Error; err != nil {
		return err
	}

	return t.Commit()
}

func (r *cartRepository) Checkout(userID uuid.UUID) (err error) {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	var cart Cart
	if err = t.DB().Preload("Products").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	// สร้าง PurchaseHistory
	purchase := history.PurchaseHistory{
		UserID: userID,
		Items:  []history.PurchaseItem{},
	}

	for _, cp := range cart.Products {
		purchase.Items = append(purchase.Items, history.PurchaseItem{
			ProductID: cp.ProductID,
			Quantity:  cp.Quantity,
		})
	}

	if err = t.DB().Create(&purchase).Error; err != nil {
		return err
	}

	// เคลียร์ cart
	if err = t.DB().Where("cart_id = ?", cart.ID).Delete(&CartProduct{}).Error; err != nil {
		return err
	}
	if err = t.DB().Delete(&cart).Error; err != nil {
		return err
	}

	return t.Commit()
}
