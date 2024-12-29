package usecase

import (
	"auth-service/app/helper"
	"auth-service/app/model/dto"
	"auth-service/app/model/entity"
	"auth-service/app/model/proto"
	"auth-service/app/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ChatUseCase interface {
	GetChats(ctx context.Context, gcr *proto.GetChatsRequest) (*proto.GetChatsResponse, error)
	GetChatByIdAndUserId(ctx context.Context, gcr *proto.GetChatRequest) (*proto.GetChatResponse, error)
}

type ChatUseCaseImpl struct {
	DB              *gorm.DB
	ChatRepository  repository.ChatRepository
	Log             *logrus.Logger
	Validate        *validator.Validate
}

func NewChatUseCaseImpl(
	db *gorm.DB,
	chatRepository repository.ChatRepository,
	log *logrus.Logger,
	validate *validator.Validate,
) *ChatUseCaseImpl {
	return &ChatUseCaseImpl{
		DB:             db,
		ChatRepository: chatRepository,
		Log:            log,
		Validate:       validate,
	}
}

func (c *ChatUseCaseImpl) GetChats(ctx context.Context, gcr *proto.GetChatsRequest) (*proto.GetChatsResponse, error) {
	// create request
	req := &dto.GetChatsRequest{
		UserID: uint(gcr.GetUserId()),
	}

	// validate request
	if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
		err := status.Newf(codes.InvalidArgument, "Validation errors: %v", errors)

		err, wde := err.WithDetails(gcr)
		if wde != nil {
			c.Log.Errorf("An error occurred while adding details to error: %v", wde)
			return nil, status.Error(codes.Internal, "An error occurred while adding details to error")
		}

		return nil, err.Err()
	}

	db := c.DB.WithContext(ctx)

	// get chats
	chats := new([]entity.Chat)
	if err := c.ChatRepository.FindAllByUserID(db, chats, int(req.UserID)); err != nil {
		c.Log.Errorf("Failed to get chats: %v", err)
		return nil, status.Error(codes.NotFound, "Chats not found")
	}

	// create response
	res := &proto.GetChatsResponse{
		Chats: make([]*proto.Chat, 0),
	}

	for _, chat := range *chats {
		res.Chats = append(res.Chats, &proto.Chat{
			Id:     int64(chat.ID),
			UserId: int64(chat.UserID),
			Topic: *chat.Topic,
		})
	}

	return res, nil
}

func (c *ChatUseCaseImpl) GetChatByIdAndUserId(ctx context.Context, gcr *proto.GetChatRequest) (*proto.GetChatResponse, error) {
	// create request
	req := &dto.GetChatRequest{
		ChatID: uint(gcr.GetChatId()),
		UserID: uint(gcr.GetUserId()),
	}

	// validate request
	if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
		err := status.Newf(codes.InvalidArgument, "Validation errors: %v", errors)

		err, wde := err.WithDetails(gcr)
		if wde != nil {
			c.Log.Errorf("An error occurred while adding details to error: %v", wde)
			return nil, status.Error(codes.Internal, "An error occurred while adding details to error")
		}

		return nil, err.Err()
	}

	db := c.DB.WithContext(ctx)

	// get chat
	chat := new(entity.Chat)
	if err := c.ChatRepository.FindByIDAndUserID(db, chat, int(req.ChatID), int(req.UserID)); err != nil {
		c.Log.Errorf("Failed to get chat: %v", err)
		return nil, status.Error(codes.NotFound, "Chat not found")
	}

	// get messages
	messages := make([]*proto.FullMessage, 0)
	for _, message := range chat.Messages {
		messages = append(messages, &proto.FullMessage{
			Id:     int64(message.ID),
			ChatId: int64(message.ChatID),
			UserId: int64(message.UserID),
			Role:  message.Role,
			Content: message.Content,
		})
	}


	// create response
	res := &proto.GetChatResponse{
		Chat: &proto.Chat{
			Id:     int64(chat.ID),
			UserId: int64(chat.UserID),
			Topic: *chat.Topic,
		},
		Messages: messages,
	}

	return res, nil
}
