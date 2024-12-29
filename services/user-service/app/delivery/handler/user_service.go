package handler

import (
	"context"
	"user-service/app/model/proto"
	"user-service/app/usecase"
)

type UserService interface {
    proto.UserServiceServer
    GetUserByEmail(ctx context.Context, req *proto.GetUserByEmailRequest) (*proto.GetUserResponse, error)
    GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserResponse, error)
    CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.GetUserResponse, error)
    UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.GetUserResponse, error)
}

type UserServiceImpl struct {
    proto.UnimplementedUserServiceServer
    userUseCase usecase.UserUseCase
}

func NewUserServiceImpl(userUseCase usecase.UserUseCase) *UserServiceImpl {
    return &UserServiceImpl{
        userUseCase: userUseCase,
    }
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, req *proto.GetUserByEmailRequest) (*proto.GetUserResponse, error) {
    user, err := s.userUseCase.GetByEmail(ctx, req)
    if err != nil {
        return nil, err
    }

    return &proto.GetUserResponse{
        Id:    int64(user.ID),
        Name:  user.Name,
        ImageUrl: *user.ImageURL,
        Email: user.Email,
        Password: *user.Password,
    }, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserResponse, error) {
    user, err := s.userUseCase.GetByID(ctx, req)
    if err != nil {
        return nil, err
    }

    return &proto.GetUserResponse{
        Id:    int64(user.ID),
        Name:  user.Name,
        ImageUrl: *user.ImageURL,
        Email: user.Email,
        Password: *user.Password,
    }, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.GetUserResponse, error) {
    user, err := s.userUseCase.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    return &proto.GetUserResponse{
        Id:    int64(user.ID),
        Name:  user.Name,
        ImageUrl: *user.ImageURL,
        Email: user.Email,
        Password: *user.Password,
    }, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.GetUserResponse, error) {
    user, err := s.userUseCase.Update(ctx, req)
    if err != nil {
        return nil, err
    }

    return &proto.GetUserResponse{
        Id:    int64(user.ID),
        Name:  user.Name,
        ImageUrl: *user.ImageURL,
        Email: user.Email,
        Password: *user.Password,
    }, nil
}