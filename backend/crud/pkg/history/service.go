package history

import "github.com/google/uuid"

type HistoryService struct {
	repo HistoryRepository
}

func NewHistoryService(repo HistoryRepository) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) GetByUser(userID uuid.UUID) ([]PurchaseHistory, error) {
	return s.repo.GetByUser(userID)
}

func (s *HistoryService) Create(history *PurchaseHistory) error {
	return s.repo.Create(history)
}
