package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/evgen6501-star/golang-subscriptions/internal/config"
	"github.com/evgen6501-star/golang-subscriptions/internal/domain"
	"github.com/evgen6501-star/golang-subscriptions/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepoPostgres struct {
	db     *pgx.Conn
	logger *logger.Logger
}

func NewSubscriptionRepoPostgres(db *pgx.Conn, log *logger.Logger) *SubscriptionRepoPostgres {
	return &SubscriptionRepoPostgres{
		db:     db,
		logger: log,
	}
}

func (r *SubscriptionRepoPostgres) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `
	INSERT INTO subscriptions(id, service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	
	`
	var endDate sql.NullTime
	if sub.EndDate != nil {
		endDate = sql.NullTime{
			Time:  *sub.EndDate,
			Valid: true,
		}
	}
	err := r.db.QueryRow(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		endDate,
	).Scan(&sub.ID)
	if err != nil {
		r.logger.Error("failed to create subscription")
		return fmt.Errorf("failed to create subscription : %w", err)
	}
	r.logger.Info("Подписка создана")
	return nil
}
func (r *SubscriptionRepoPostgres) GetTotalPrice(ctx context.Context, userID *uuid.UUID, serviseName *string, startDate, endDate time.Time) (int, error) {
	query := `
SELECT COALESCE(SUM(price), 0)
FROM subscriptions
WHERE start_date >= $1 AND start_date <= $2
`
	args := []interface{}{startDate, endDate}
	argIDX := 3
	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIDX)
		args = append(args, *userID)
		argIDX++
	}
	if serviseName != nil && *serviseName != "" {
		query += fmt.Sprintf(" AND service_name =$%d", argIDX)
		args = append(args, *serviseName)
		argIDX++
	}
	var total int
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		r.logger.Error("failed to calculate total")
		return 0, fmt.Errorf("failed to calculate total^ %w", err)
	}
	return total, nil

}
func (r *SubscriptionRepoPostgres) List(ctx context.Context, limit, offset int) ([]*domain.Subscription, error) {
	query := `
SELECT id, service_name, price, user_id, start_date, end_date
FROM subscriptions
ORDER BY  user_id, service_name
LIMIT $1 OFFSET $2
`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list subscriptions")
		return nil, fmt.Errorf("failed to list subscriptions %w", err)

	}
	defer rows.Close()
	var subscribtions []*domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var endDate sql.NullTime
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to Scan subscription %w ", err)
		}
		if endDate.Valid {
			sub.EndDate = &endDate.Time
		}
		subscribtions = append(subscribtions, &sub)

	}
	return subscribtions, nil
}
func (r *SubscriptionRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
SELECT id, service_name, price, user_id, start_date, end_date
FROM subscriptions
WHERE id = $1
`
	var sub domain.Subscription
	var endDate sql.NullTime
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&endDate,
	)
	if err != nil {
		r.logger.Error("failed to get subscription")
		return nil, fmt.Errorf("failed to get subscription : %w", err)
	}
	if endDate.Valid {
		sub.EndDate = &endDate.Time
	}
	return &sub, nil
}
func NewPostgresConnection(cfg *config.Config) (*pgx.Conn, error) {
	dsn := cfg.DatabaseDSN()

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return conn, nil
}
