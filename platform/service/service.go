package service

import (
	"context"

	upb "github.com/H0lyDiv3r/ticketing-grpc/gen/user"
)

type Service struct {
	userClient upb.UserServiceClient
}

func NewService(userClient upb.UserServiceClient) *Service {
	return &Service{userClient: userClient}
}

func (s *Service) Signup(ctx context.Context, username, email, password string) (*upb.RegisterResponse, error) {
	return s.userClient.Register(ctx, &upb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
}
