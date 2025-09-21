package cart

import "github.com/google/uuid"

type CartService struct {
	repo CartRepository
}

func NewCartService(repo CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCart(userID uuid.UUID) (*Cart, error) {
	return s.repo.GetCart(userID)
}

func (s *CartService) AddProduct(userID, productID uuid.UUID, qty int) error {
	return s.repo.AddProduct(userID, productID, qty)
}

func (s *CartService) RemoveProduct(userID, productID uuid.UUID) error {
	return s.repo.RemoveProduct(userID, productID)
}

func (s *CartService) ClearCart(userID uuid.UUID) error {
	return s.repo.ClearCart(userID)
}

func (s *CartService) Checkout(userID uuid.UUID) error {
	return s.repo.Checkout(userID)
}
