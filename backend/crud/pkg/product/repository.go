package product

import (
	"crud/pkg/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *Product) error
	GetAll() ([]Product, error)
	GetByID(id uuid.UUID) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// ✅ Create with transaction
func (r *productRepository) Create(product *Product) (err error) {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	if err = t.DB().Create(product).Error; err != nil {
		return err
	}

	return t.Commit()
}

// ✅ GetAll (read-only ไม่ต้องใช้ transaction)
func (r *productRepository) GetAll() ([]Product, error) {
	var products []Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// ✅ GetByID (read-only)
func (r *productRepository) GetByID(id uuid.UUID) (*Product, error) {
	var product Product
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// ✅ Update with transaction
func (r *productRepository) Update(product *Product) (err error) {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	if err = t.DB().Save(product).Error; err != nil {
		return err
	}

	return t.Commit()
}

// ✅ Delete with transaction
func (r *productRepository) Delete(id uuid.UUID) (err error) {
	t, err := transaction.Begin(r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			t.Rollback()
		}
	}()

	if err = t.DB().Delete(&Product{}, "id = ?", id).Error; err != nil {
		return err
	}

	return t.Commit()
}
