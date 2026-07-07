package service

import (
	"context"
	"time"

	"github.com/evgen6501-star/golang-subscriptions/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, req *CreateSubscriptionRequest) (*domain.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Subscription, error)
	GetTotalPrice(ctx context.Context, req TotalPriceRequest) (int, error)
}

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
type TotalPriceRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
}
