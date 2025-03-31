package service

import (
	"github.com/ombima56/insights-edge/backend/models"
	"github.com/ombima56/insights-edge/backend/repository"
)

type InsightService struct {
	insightRepository repository.InsightRepository
	userRepository    repository.UserRepository
}

func NewInsightService(insightRepository repository.InsightRepository, userRepository repository.UserRepository) *InsightService {
	return &InsightService{
		insightRepository: insightRepository,
		userRepository:    userRepository,
	}
}

func (s *InsightService) CreateInsight(insight *models.Insight) error {
	return s.insightRepository.CreateInsight(insight)
}

func (s *InsightService) GetInsight(id int64) (*models.Insight, error) {
	return s.insightRepository.GetInsight(id)
}

func (s *InsightService) GetInsights() ([]*models.Insight, error) {
	return s.insightRepository.GetInsights()
}

func (s *InsightService) PurchaseInsight(insightID int64, walletAddr string) error {
	// Verify user exists
	_, err := s.userRepository.GetUserByWallet(walletAddr)
	if err != nil {
		return err
	}

	// Create purchase
	purchase := &models.Purchase{
		InsightID: insightID,
		Buyer:     walletAddr,
	}

	return s.insightRepository.CreatePurchase(purchase)
}

func (s *InsightService) GetMyInsights(walletAddr string) ([]*models.Insight, error) {
	purchases, err := s.insightRepository.GetPurchasesByWallet(walletAddr)
	if err != nil {
		return nil, err
	}

	var insights []*models.Insight
	for _, purchase := range purchases {
		insight, err := s.insightRepository.GetInsight(purchase.InsightID)
		if err != nil {
			return nil, err
		}
		insights = append(insights, insight)
	}

	return insights, nil
}

func (s *InsightService) GetMyPurchases(walletAddr string) ([]*models.Purchase, error) {
	return s.insightRepository.GetPurchasesByWallet(walletAddr)
}
