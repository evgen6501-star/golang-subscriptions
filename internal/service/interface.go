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
	Update(ctx context.Context, id uuid.UUID, req *UpdateSubscriptionRequest) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int       `json:"price,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
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
