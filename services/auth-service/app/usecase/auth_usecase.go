package usecase

import (
	"auth-service/app/delivery/client"
	"auth-service/app/helper"
	"auth-service/app/model/dto"
	"auth-service/app/model/proto"
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUseCase interface {
    Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthenticatedResponse, error)
    Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthenticatedResponse, error)
    LoginWithGoogle(ctx context.Context, req *proto.LoginWithGoogleRequest) (*proto.AuthenticatedResponse, error)
	GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error)
	Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.AuthenticatedResponse, error)
}

type AuthUseCaseImpl struct {
    Validate        *validator.Validate
    Jwt             helper.JWTHelper
    Log             *logrus.Logger
	User 			*client.UserClient
	Oauth			helper.OAuth
}

func NewAuthUseCaseImpl(
	validate *validator.Validate,
	jwt helper.JWTHelper,
	log *logrus.Logger,
	user *client.UserClient,
	oauth helper.OAuth,
) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		Validate: validate,
		Jwt:      jwt,
		Log:      log,
		User:     user,
		Oauth:    oauth,
	}
}

func (a *AuthUseCaseImpl) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthenticatedResponse, error) {
	// Create request
	request := &dto.RegisterRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	// Validate request
	if errors := helper.Validate(a.Validate, request); len(errors) > 0 {
		messages := helper.GetErrorMessages(errors)
		return nil, status.Error(codes.InvalidArgument, messages[0])
	}

	// Validate if email already exists
	getUserCtx, cancelGetUser := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelGetUser()

	_, err := a.User.Service.GetUserByEmail(getUserCtx, &proto.GetUserByEmailRequest{Email: req.GetEmail()})
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}

	// Create user
	createUserCtx, cancelCreateUser := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelCreateUser()

	user, err := a.User.Service.CreateUser(createUserCtx, &proto.CreateUserRequest{
		Name:    	req.GetName(),
		Email:    	req.GetEmail(),
		Password: 	req.GetPassword(),
	})

	fmt.Println("User: ", user)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, status.Error(codes.InvalidArgument, s.Message())
			} 
			if s.Code() == codes.AlreadyExists {
				return nil, status.Error(codes.AlreadyExists, s.Message())
			} 
			return nil, status.Error(codes.Internal, s.Message())
		}

		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// Create JWT
	accessToken, refreshToken, expiredAt, err := a.Jwt.GenerateTokens(uint(user.GetId()), user.GetName(), user.GetEmail(), user.GetImageUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate tokens")
	}

	return &proto.AuthenticatedResponse{
		Name:           user.GetName(),
		Email:          user.GetEmail(),
		ImageUrl: 	 	user.GetImageUrl(),
		Token: &proto.Token{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: expiredAt,
		},
	}, nil
}

func (a *AuthUseCaseImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthenticatedResponse, error) {
	// Create request
	request := &dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	// Validate request
	if errors := helper.Validate(a.Validate, request); len(errors) > 0 {
		messages := helper.GetErrorMessages(errors)
		return nil, status.Error(codes.InvalidArgument, messages[0])
	}

	// Validate user
	getUserCtx, cancelGetUser := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelGetUser()

	res, err := a.User.Service.GetUserByEmail(getUserCtx, &proto.GetUserByEmailRequest{Email: req.GetEmail()})
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	fmt.Println("User: ", res)

	// Validate password
	if !helper.VerifyPassword(res.GetPassword(), req.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "Invalid password")
	}

	// Create JWT
	accessToken, refreshToken, expiredAt, err := a.Jwt.GenerateTokens(uint(res.GetId()), res.GetName(), res.GetEmail(), res.GetImageUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate tokens")
	}

	return &proto.AuthenticatedResponse{
		Name:           res.GetName(),
		Email:          res.GetEmail(),
		ImageUrl: 	 res.GetImageUrl(),
		Token: &proto.Token{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: expiredAt,
		},
	}, nil
}

func (a *AuthUseCaseImpl) LoginWithGoogle(ctx context.Context, req *proto.LoginWithGoogleRequest) (*proto.AuthenticatedResponse, error) {
	// Create request
	request := &dto.LoginWIthGoogleRequest{
		IdToken: req.GetIdToken(),
	}

	// Validate request
	if errors := helper.Validate(a.Validate, request); len(errors) > 0 {
		messages := helper.GetErrorMessages(errors)
		return nil, status.Error(codes.InvalidArgument, messages[0])
	}

	// Verify id token
	credential, err  := a.Oauth.VerifyIdToken(req.GetIdToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id token")
	}

	// If user exists, return JWT
	getUserCtx, cancelGetUser := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelGetUser()

	res, err := a.User.Service.GetUserByEmail(getUserCtx, &proto.GetUserByEmailRequest{Email: credential.Email})
	if err == nil {
		accessToken, refreshToken, expiredAt, err := a.Jwt.GenerateTokens(uint(res.GetId()), res.GetName(), res.GetEmail(), res.GetImageUrl())
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to generate tokens")
		}

		return &proto.AuthenticatedResponse{
			Name:           res.GetName(),
			Email:          res.GetEmail(),
			ImageUrl: 	 res.GetImageUrl(),
			Token: &proto.Token{
				AccessToken: accessToken,
				RefreshToken: refreshToken,
				ExpiredAt: expiredAt,
			},
		}, nil
	}

	// If user does not exist, create user and return JWT
	createUserCtx, cancelCreateUser := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelCreateUser()

	user, err := a.User.Service.CreateUser(createUserCtx, &proto.CreateUserRequest{
		Name:     credential.Name,
		Email:    credential.Email,
		ImageUrl: credential.Picture,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				return nil, status.Error(codes.InvalidArgument, s.Message())
			} 
			if s.Code() == codes.AlreadyExists {
				return nil, status.Error(codes.AlreadyExists, s.Message())
			} 
			return nil, status.Error(codes.Internal, s.Message())
		}

		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// Create JWT
	res, err = a.User.Service.GetUserByEmail(ctx, &proto.GetUserByEmailRequest{Email: credential.Email})
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	accessToken, refreshToken, expiredAt, err := a.Jwt.GenerateTokens(uint(res.GetId()), user.GetName(), user.GetEmail(), user.GetImageUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate tokens")
	}

	return &proto.AuthenticatedResponse{
		Name:          user.GetName(),
		Email:         user.GetEmail(),
		ImageUrl: 	   user.GetImageUrl(),
		Token: &proto.Token{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: expiredAt,
		},
	}, nil
}

func (a *AuthUseCaseImpl) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
	// Create request
	request := &dto.GetProfileRequest{
		AccessToken: req.GetAccessToken(),
	}

	// Validate request
	if errors := helper.Validate(a.Validate, request); len(errors) > 0 {
		messages := helper.GetErrorMessages(errors)
		return nil, status.Error(codes.InvalidArgument, messages[0])
	}

	// Validate JWT
	userId, name, email, imageUrl, err := a.Jwt.ValidateAccessToken(req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid access token")
	}

	return &proto.GetProfileResponse{
		Id: 	   	int64(userId),
		Name:     	name,
		Email:    	email,
		ImageUrl: 	imageUrl,
	}, nil
}

func (a *AuthUseCaseImpl) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.AuthenticatedResponse, error) {
	// Create request
	request := &dto.RefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	}

	// Validate request
	if errors := helper.Validate(a.Validate, request); len(errors) > 0 {
		messages := helper.GetErrorMessages(errors)
		return nil, status.Error(codes.InvalidArgument, messages[0])
	}

	// Validate JWT
	userId, name, email, imageUrl, err := a.Jwt.ValidateRefreshToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid refresh token")
	}

	// Create JWT
	accessToken, refreshToken, expiredAt, err := a.Jwt.GenerateTokens(userId, name, email, imageUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate tokens")
	}

	return &proto.AuthenticatedResponse{
		Name:           name,
		Email:          email,
		ImageUrl: 		imageUrl,
		Token: &proto.Token{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: expiredAt,
		},
	}, nil
}
