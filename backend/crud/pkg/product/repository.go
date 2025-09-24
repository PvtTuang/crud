package product

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *Product) error
	GetAll() ([]Product, error)
	GetByID(id uuid.UUID) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *Product) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	result, errrs := tx.Exec(`INSERT INTO products (name, price) VALUES ($1 , $2 )`,
		product.Name, product.Price,
	)

	if errrs != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err = fmt.Errorf("no product found ")
		return err
	}
	return tx.Commit()

}

func (r *productRepository) GetAll() ([]Product, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	rows, err := tx.Query(`SELECT id , name  , price , created_at , updated_at FROM products;`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var product []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		product = append(product, p)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*Product, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var p Product
	row := tx.QueryRow(`SELECT id , name  , price , created_at , updated_at FROM products WHERE id = $1;`, id)
	err = row.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return &Product{}, err
	}
	return &p, nil
}

func (r *productRepository) Update(product *Product) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	result, errrs := tx.Exec(`UPDATE products SET name = $1 , price = $2 , created_at = NOW() , updated_at = NOW() WHERE id = $3;`,
		product.Name, product.Price, product.ID,
	)
	if errrs != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err = fmt.Errorf("no product found ")
		return err
	}

	return tx.Commit()
}

func (r *productRepository) Delete(id uuid.UUID) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	result, errrs := tx.Exec(`UPDATE products SET deleted_at = NOW() WHERE id = $1;`, id)

	if errrs != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err = fmt.Errorf("no product found ")
		return err
	}

	return tx.Commit()
}
