package repository

import (
	"context"

	"github.com/H0lyDiv3r/ticketing-grpc/user/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user domain.User) error {
	result := gorm.G[domain.User](r.db).Create(ctx, &user)
	return result
}
