package service

import (
	"context"
	"fmt"
	"time"
	db "user-api/db/sqlc"
	"user-api/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	queries db.Querier
}

func NewUserService(queries db.Querier) *UserService {
	return &UserService{queries: queries}
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return age
}

func parseDob(dobStr string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid dob format, use YYYY-MM-DD")
	}
	if t.After(time.Now()) {
		return pgtype.Date{}, fmt.Errorf("dob cannot be in the future")
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func dobToTime(d pgtype.Date) time.Time {
	return d.Time
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := parseDob(req.Dob)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  dobToTime(user.Dob).Format("2006-01-02"),
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	dob := dobToTime(user.Dob)
	return &models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  dob.Format("2006-01-02"),
		Age:  CalculateAge(dob),
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := parseDob(req.Dob)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		Name: req.Name,
		Dob:  dob,
		ID:   id,
	})
	if err != nil {
		return nil, fmt.Errorf("user not found or update failed")
	}

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  dobToTime(user.Dob).Format("2006-01-02"),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	_, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := s.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int) (*models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.queries.CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	data := make([]models.UserWithAgeResponse, len(users))
	for i, u := range users {
		dob := dobToTime(u.Dob)
		data[i] = models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			Dob:  dob.Format("2006-01-02"),
			Age:  CalculateAge(dob),
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &models.PaginatedUsersResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
