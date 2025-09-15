package pkg

import "errors"

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetAllProducts() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id uint) (*Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) UpdateProduct(product *Product) error {
	existing, err := s.repo.GetByID(product.ID)
	if err != nil || existing == nil {
		return errors.New("product not found")
	}
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	existing, err := s.repo.GetByID(id)
	if err != nil || existing == nil {
		return errors.New("product not found")
	}
	return s.repo.Delete(id)
}
