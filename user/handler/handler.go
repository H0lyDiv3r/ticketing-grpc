package handler

import (
	"context"

	upb "github.com/H0lyDiv3r/ticketing-grpc/gen/user"
	"github.com/H0lyDiv3r/ticketing-grpc/user/service"
)

type Handler struct {
	service *service.Service
	upb.UnimplementedUserServiceServer
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Signup(ctx context.Context, req *upb.RegisterRequest) (*upb.RegisterResponse, error) {
	res, err := h.service.Signup(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// func (h *Handler) Login(ctx context.Context, req *upb.LoginRequest) (*upb.LoginResponse, error) {

// }

// func (h *Handler) GetUser(ctx context.Context, req *upb.GetUserRequest) (*upb.GetUserResponse, error) {

// }

// func (h *Handler) ValidateToken(ctx context.Context, req *upb.ValidateTokenRequest) (*upb.ValidateTokenResponse, error) {

// }

// func (h *Handler) GetUserBatch(ctx context.Context, req *upb.GetUserBatchRequest) (*upb.GetUserBatchResponse, error) {

// }
