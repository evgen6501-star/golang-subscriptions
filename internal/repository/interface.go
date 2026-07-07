package repository

import (
	"context"
	"time"

	"github.com/evgen6501-star/golang-subscriptions/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionReposytory interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Subscription, error)
	GetTotalPrice(ctx context.Context, userID *uuid.UUID, serviseName *string, startDate, endDate time.Time) (int, error)
}
