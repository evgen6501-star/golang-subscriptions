package service

import (
	"context"
	"errors"

	"github.com/evgen6501-star/golang-subscriptions/internal/domain"
	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/evgen6501-star/golang-subscriptions/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionImpl struct {
	repo   repository.SubscriptionReposytory
	logger *logger.Logger
}

func NewSubscriptionService(repo repository.SubscriptionReposytory, log *logger.Logger) *SubscriptionImpl {
	return &SubscriptionImpl{
		repo:   repo,
		logger: log,
	}
}

func (s *SubscriptionImpl) Create(ctx context.Context, req *CreateSubscriptionRequest) (*domain.Subscription, error) {
	s.logger.Info("Создание подписки")
	if req.ServiceName == "" {
		return nil, errors.New("service_name is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("error price")
	}
	if req.UserID == uuid.Nil {
		return nil, errors.New("user id is requered")
	}
	sub := domain.NewSubscription(req.ServiceName, req.Price, req.UserID, req.StartDate)

	if err := s.repo.Create(ctx, sub); err != nil {
		s.logger.Error("Ошибка создания подписки")
		return nil, err

	}
	return sub, nil

}

func (s *SubscriptionImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	s.logger.Info("Получение подписки по id", zap.Any("id", id))
	return s.repo.GetByID(ctx, id)

}
func (s *SubscriptionImpl) List(ctx context.Context, limit, offset int) ([]*domain.Subscription, error) {
	s.logger.Info("Получение списка подписок", zap.Any("limit", limit), zap.Any("offset", offset))
	return s.repo.List(ctx, limit, offset)
}
func (s *SubscriptionImpl) GetTotalPrice(ctx context.Context, req TotalPriceRequest) (int, error) {
	s.logger.Info("Подсчет суммы подписок")
	return s.repo.GetTotalPrice(ctx, req.UserID, req.ServiceName, req.StartDate, req.EndDate)
}
