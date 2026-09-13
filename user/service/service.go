package service

import (
	"context"

	upb "github.com/H0lyDiv3r/ticketing-grpc/gen/user"
	"github.com/H0lyDiv3r/ticketing-grpc/user/domain"
	"github.com/H0lyDiv3r/ticketing-grpc/user/repository"
	"github.com/H0lyDiv3r/ticketing-grpc/user/utils"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Signup(ctx context.Context, req *upb.RegisterRequest) (*upb.RegisterResponse, error) {
	user := domain.User{
		Username: req.Username,
		Email:    req.Email,
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = *hashedPassword
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &upb.RegisterResponse{
		Token: token,
		User: &upb.User{
			Id:       int64(user.ID),
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}
