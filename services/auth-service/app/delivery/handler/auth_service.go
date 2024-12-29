package handler

import (
	"auth-service/app/model/proto"
	"auth-service/app/usecase"
	"context"
)

type AuthService interface {
    proto.AuthServiceServer
    Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthenticatedResponse, error)
    Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthenticatedResponse, error)
    LoginWithGoogle(ctx context.Context, req *proto.LoginWithGoogleRequest) (*proto.AuthenticatedResponse, error)
    GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error)
    Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.AuthenticatedResponse, error)
}

type AuthServiceImpl struct {
    proto.UnimplementedAuthServiceServer
    authUseCase usecase.AuthUseCase
}

func NewAuthServiceImpl(authUseCase usecase.AuthUseCase) *AuthServiceImpl {
    return &AuthServiceImpl{
        authUseCase: authUseCase,
    }
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthenticatedResponse, error) {
    res, err := s.authUseCase.Register(ctx, req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthenticatedResponse, error) {
    res, err := s.authUseCase.Login(ctx, req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *AuthServiceImpl) LoginWithGoogle(ctx context.Context, req *proto.LoginWithGoogleRequest) (*proto.AuthenticatedResponse, error) {
    res, err := s.authUseCase.LoginWithGoogle(ctx, req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *AuthServiceImpl) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
    res, err := s.authUseCase.GetProfile(ctx, req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *AuthServiceImpl) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.AuthenticatedResponse, error) {
    res, err := s.authUseCase.Refresh(ctx, req)
    if err != nil {
        return nil, err
    }

    return res, nil
}
