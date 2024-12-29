package usecase

import (
	"context"
	"user-service/app/helper"
	"user-service/app/model/dto"
	"user-service/app/model/entity"
	"user-service/app/model/proto"
	"user-service/app/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserUseCase interface {
    GetByEmail(ctx context.Context, gur *proto.GetUserByEmailRequest) (*entity.User, error)
    GetByID(ctx context.Context, gur *proto.GetUserByIDRequest) (*entity.User, error)
    Create(ctx context.Context, cur *proto.CreateUserRequest) (*entity.User, error)
    Update(ctx context.Context, uur *proto.UpdateUserRequest) (*entity.User, error)
}

type UserUseCaseImpl struct {
    DB              *gorm.DB
    Validate        *validator.Validate
    UserRepository  repository.UserRepository
    Log             *logrus.Logger
}

func NewUserUseCaseImpl(
    db *gorm.DB,
    validator *validator.Validate,
    userRepository repository.UserRepository,
    log *logrus.Logger,
) *UserUseCaseImpl {
    return &UserUseCaseImpl{
        DB:             db,
        Validate:       validator,
        UserRepository: userRepository,
        Log:            log,
    }
}

func (c *UserUseCaseImpl) GetByEmail(ctx context.Context, gur *proto.GetUserByEmailRequest) (*entity.User, error) {
    db := c.DB.WithContext(ctx)

    user := new(entity.User)
    if err := c.UserRepository.FindByEmail(db, user, gur.GetEmail()); err != nil {
        return nil, status.Error(codes.NotFound, "User not found")
    }

    return user, nil
}

func (c *UserUseCaseImpl) GetByID(ctx context.Context, gur *proto.GetUserByIDRequest) (*entity.User, error) {
    db := c.DB.WithContext(ctx)

    user := new(entity.User)
    if err := c.UserRepository.Find(db, user, int(gur.GetId())); err != nil {
        return nil, status.Error(codes.NotFound, "User not found")
    }

    return user, nil
}

func (c *UserUseCaseImpl) Create(ctx context.Context, cur *proto.CreateUserRequest) (*entity.User, error) {
    // Create request
    req := &dto.CreateUserRequest{
        Name:     cur.GetName(),
        Email:    cur.GetEmail(),
        Password: cur.GetPassword(),
        ImageURL: cur.GetImageUrl(),
    }

    // Validate request
    if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
        messages := helper.GetErrorMessages(errors)
        err := status.Newf(codes.InvalidArgument, "%s", messages[0])

        err, wde := err.WithDetails(cur)
        if wde != nil {
            return nil, status.Error(codes.Internal, "An error occurred while adding details to error")
        }

        return nil, err.Err()
    }

    db := c.DB.WithContext(ctx)

    // Validate if email already exists
    user := new(entity.User)
    if err := c.UserRepository.FindByEmail(db, user, req.Email); err == nil {
        return nil, status.Error(codes.AlreadyExists, "User already exists")
    }

    // Create user
    user = &entity.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: &req.Password,
        ImageURL: &req.ImageURL,
    }

    if user.Password != nil {
        hashedPassword, err := helper.HashPassword(*user.Password)
        if err != nil {
            return nil, status.Error(codes.Internal, "Failed to hash password")
        }

        user.Password = &hashedPassword
    }

    if err := c.UserRepository.Create(db, user); err != nil {
        return nil, status.Error(codes.Internal, "Failed to create user")
    }

    return user, nil
}

func (c *UserUseCaseImpl) Update(ctx context.Context, uur *proto.UpdateUserRequest) (*entity.User, error) {
    // Create request
    req := &dto.UpdateUserRequest{
        ID:       uint(uur.GetId()),
        Name:     uur.GetName(),
        Email:    uur.GetEmail(),
        Password: uur.GetPassword(),
        ImageURL: uur.GetImageUrl(),
    }

    // Validate request
    if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
        err := status.Newf(codes.InvalidArgument, "Validation errors: %v", errors)

        err, wde := err.WithDetails(uur)
        if wde != nil {
            return nil, status.Error(codes.Internal, "An error occurred while adding details to error")
        }

        return nil, err.Err()
    }

    db := c.DB.WithContext(ctx)

    // Validate if user exists
    user := new(entity.User)
    if err := c.UserRepository.Find(db, user, int(req.ID)); err != nil {
        return nil, status.Error(404, "User not found")
    }

    // Validate if email already exist but not the same user
    uniqueUser := new(entity.User)
    if err := c.UserRepository.FindByEmail(db, uniqueUser, req.Email); err == nil && uniqueUser.ID != req.ID {
        return nil, status.Error(codes.AlreadyExists, "Email already exists")
    }

    // Update user
    user.Name = req.Name
    user.Email = req.Email
    user.Password = &req.Password
    user.ImageURL = &req.ImageURL

    if req.Password != "" {
        hashedPassword, err := helper.HashPassword(req.Password)
        if err != nil {
            return nil, status.Error(codes.Internal, "Failed to hash password")
        }

        user.Password = &hashedPassword
    }

    if err := c.UserRepository.Update(db, user); err != nil {
        return nil, status.Error(codes.Internal, "Failed to update user")
    }

    return user, nil
}