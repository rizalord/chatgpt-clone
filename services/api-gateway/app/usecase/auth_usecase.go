package usecase

import (
	"api-gateway/app/delivery/client"
	"api-gateway/app/helper"
	"api-gateway/app/model"
	"api-gateway/app/model/dto"
	"api-gateway/app/model/proto"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUseCase interface {
    Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error)
    Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
    LoginWithGoogle(ctx context.Context, req *dto.LoginWIthGoogleRequest) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error)
}

type AuthUseCaseImpl struct {
    Validate        *validator.Validate
    Log             *logrus.Logger
	Auth 			*client.AuthClient
}

func NewAuthUseCaseImpl(
	validate *validator.Validate,
	log *logrus.Logger,
	auth *client.AuthClient,
) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		Validate: validate,
		Log:      log,
		Auth:     auth,
	}
}

func (a *AuthUseCaseImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validate request
    if errors := helper.Validate(a.Validate, req); len(errors) > 0 {
        return nil, model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
    }

	// Get response from auth service
	registerCtx, cancelRegister := context.WithTimeout(ctx, 10 * time.Second)
	defer cancelRegister()

	res, err := a.Auth.Service.Register(registerCtx, &proto.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, model.NewError(model.StatusBadRequest, s.Message(), nil)
			} 
			if s.Code() == codes.AlreadyExists {
				return nil, model.NewError(model.StatusConflict, s.Message(), nil)
			} 
			return nil, model.NewError(model.StatusInternalServerError, s.Message(), nil)
		}

		return nil, model.NewError(model.StatusInternalServerError, "Failed to register", nil)
	}

	imageURL := res.GetImageUrl()

	return &dto.LoginResponse{
		User: dto.CredentialData{
			Name:     res.GetName(),
			Email:    res.GetEmail(),
			ImageURL: &imageURL,
		},
		Token: dto.TokenData{
			AccessToken: res.GetToken().GetAccessToken(),
			RefreshToken: res.GetToken().GetRefreshToken(),
			ExpiredAt: res.GetToken().GetExpiredAt(),
		},
	}, nil
}

func (a *AuthUseCaseImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validate request
	if errors := helper.Validate(a.Validate, req); len(errors) > 0 {
		return nil, model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
	}

	// Get response from auth service
	loginCtx, cancelLogin := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelLogin()

	res, err := a.Auth.Service.Login(loginCtx, &proto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, model.NewError(model.StatusBadRequest, s.Message(), nil)
			} 
			if s.Code() == codes.NotFound {
				return nil, model.NewError(model.StatusNotFound, s.Message(), nil)
			} 
			return nil, model.NewError(model.StatusInternalServerError, s.Message(), nil)
		}

		return nil, model.NewError(model.StatusInternalServerError, "Failed to login", nil)
	}

	imageURL := res.GetImageUrl()

	return &dto.LoginResponse{
		User: dto.CredentialData{
			Name:     res.GetName(),
			Email:    res.GetEmail(),
			ImageURL: &imageURL,
		},
		Token: dto.TokenData{
			AccessToken: res.GetToken().GetAccessToken(),
			RefreshToken: res.GetToken().GetRefreshToken(),
			ExpiredAt: res.GetToken().GetExpiredAt(),
		},
	}, nil
}

func (a *AuthUseCaseImpl) LoginWithGoogle(ctx context.Context, req *dto.LoginWIthGoogleRequest) (*dto.LoginResponse, error) {
	// Validate request
	if errors := helper.Validate(a.Validate, req); len(errors) > 0 {
		return nil, model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
	}

	// Get response from auth service
	loginWithGoogleCtx, cancelLoginWithGoogle := context.WithTimeout(ctx, 5 * time.Second)
	defer cancelLoginWithGoogle()

	res, err := a.Auth.Service.LoginWithGoogle(loginWithGoogleCtx, &proto.LoginWithGoogleRequest{
		IdToken: req.IdToken,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, model.NewError(model.StatusBadRequest, s.Message(), nil)
			} 
			if s.Code() == codes.NotFound {
				return nil, model.NewError(model.StatusNotFound, s.Message(), nil)
			} 
			if s.Code() == codes.AlreadyExists {
				return nil, model.NewError(model.StatusConflict, s.Message(), nil)
			} 
			return nil, model.NewError(model.StatusInternalServerError, s.Message(), nil)
		}

		return nil, model.NewError(model.StatusInternalServerError, "Failed to login with google", nil)
	}

	imageURL := res.GetImageUrl()

	return &dto.LoginResponse{
		User: dto.CredentialData{
			Name:     res.GetName(),
			Email:    res.GetEmail(),
			ImageURL: &imageURL,
		},
		Token: dto.TokenData{
			AccessToken: res.GetToken().GetAccessToken(),
			RefreshToken: res.GetToken().GetRefreshToken(),
			ExpiredAt: res.GetToken().GetExpiredAt(),
		},
	}, nil
}

func (a *AuthUseCaseImpl) Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error) {
	// Validate request
	if errors := helper.Validate(a.Validate, req); len(errors) > 0 {
		return nil, model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
	}

	// Get response from auth service
	refreshCtx, cancelRefresh := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelRefresh()

	res, err := a.Auth.Service.Refresh(refreshCtx, &proto.RefreshRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, model.NewError(model.StatusBadRequest, s.Message(), nil)
			} 
			if s.Code() == codes.NotFound {
				return nil, model.NewError(model.StatusNotFound, s.Message(), nil)
			} 
			return nil, model.NewError(model.StatusInternalServerError, s.Message(), nil)
		}

		return nil, model.NewError(model.StatusInternalServerError, "Failed to refresh token", nil)
	}

	imageURL := res.GetImageUrl()

	return &dto.LoginResponse{
		User: dto.CredentialData{
			Name:     res.GetName(),
			Email:    res.GetEmail(),
			ImageURL: &imageURL,
		},
		Token: dto.TokenData{
			AccessToken: res.GetToken().GetAccessToken(),
			RefreshToken: res.GetToken().GetRefreshToken(),
			ExpiredAt: res.GetToken().GetExpiredAt(),
		},
	}, nil
}