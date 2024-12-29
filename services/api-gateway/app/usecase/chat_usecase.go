package usecase

import (
	"api-gateway/app/delivery/client"
	"api-gateway/app/helper"
	"api-gateway/app/model"
	"api-gateway/app/model/dto"
	"api-gateway/app/model/proto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	socketio "github.com/doquangtan/gofiber-socket.io"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatUseCase interface {
	GetChats(ctx context.Context, userID *uint) (*dto.GetChatsResponse, error)
	GetMessages(ctx context.Context, chatID *string, userID *uint) (*dto.GetMessagesResponse, error)
	// InitChatConnection(kws *socketio.Websocket)
	JoinChatRoom(ep *socketio.EventPayload)
	CreateMessage(ep *socketio.EventPayload)
}

type ChatUseCaseImpl struct {
    Validate        *validator.Validate
    Log             *logrus.Logger
	Auth 			*client.AuthClient
	Chat 			*client.ChatClient
}

func NewChatUseCaseImpl(
	validate *validator.Validate,
	log *logrus.Logger,
	auth *client.AuthClient,
	chat *client.ChatClient,
) *ChatUseCaseImpl {
	return &ChatUseCaseImpl{
		Validate: validate,
		Log:      log,
		Auth:     auth,
		Chat:     chat,
	}
}

func (u *ChatUseCaseImpl) GetChats(ctx context.Context, userID *uint) (*dto.GetChatsResponse, error) {
	// Get response from chat service
	getChatsCtx, cancelGetChats := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelGetChats()

	res, err := u.Chat.Service.GetChats(getChatsCtx, &proto.GetChatsRequest{
		UserId: int64(*userID),
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

		return nil, model.NewError(model.StatusInternalServerError, "Failed to get chats", nil)
	}

	chats := make([]dto.ChatData, 0)
	for _, chat := range res.GetChats() {
		chats = append(chats, dto.ChatData{
			ID:   uint(chat.GetId()),
			UserID: uint(chat.GetUserId()),
			Topic: chat.GetTopic(),
		})
	}

	return &dto.GetChatsResponse{
		Chats: chats,
	}, nil
}

func (u *ChatUseCaseImpl) GetMessages(ctx context.Context, chatID *string, userID *uint) (*dto.GetMessagesResponse, error) {
	// Convert chatID to uint
	chatIDUint, err := strconv.ParseUint(*chatID, 10, 64)
	if err != nil {
		return nil, model.NewError(model.StatusBadRequest, "Chat ID must be a number", nil)
	}

	// Create request
	request := &dto.GetMessagesRequest{
		ChatID: uint(chatIDUint),
		UserID: *userID,
	}

	// Validate request
	if errors := helper.Validate(u.Validate, request); len(errors) > 0 {
		return nil, model.NewError(model.StatusBadRequest, "The given data was invalid", errors)
	}

	// Get response from chat service
	getMessagesCtx, cancelGetMessages := context.WithTimeout(ctx, 2 * time.Second)
	defer cancelGetMessages()

	res, err := u.Chat.Service.GetMessages(getMessagesCtx, &proto.GetMessagesRequest{
		ChatId: int64(chatIDUint),
		UserId: int64(*userID),
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

		return nil, model.NewError(model.StatusInternalServerError, "Failed to get messages", nil)
	}

	messages := make([]dto.MessageData, 0)
	for _, message := range res.GetMessages() {
		messages = append(messages, dto.MessageData{
			ID:      uint(message.GetId()),
			ChatID:  uint(message.GetChatId()),
			UserID:  uint(message.GetUserId()),
			Role:    message.GetRole(),
			Content: message.GetContent(),
		})
	}

	return &dto.GetMessagesResponse{
		Messages: messages,
	}, nil
}

func (u *ChatUseCaseImpl) getUser(ep *socketio.EventPayload) (*proto.GetProfileResponse, error) {
	accessToken := ep.Socket.Conn.Query("token")
	if accessToken == "" {
		return nil, errors.New("Unauthorized")
	}

	user, err := u.Auth.Service.GetProfile(context.Background(), &proto.GetProfileRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *ChatUseCaseImpl) JoinChatRoom(ep *socketio.EventPayload) {
	// Get user
	user, err := u.getUser(ep)
	if err != nil {
		ep.Socket.Emit("error", "Unauthorized")
		u.Log.Errorf("Failed to get user: %v", err)
		return
	}

	// Parse data
	dataMap, ok := ep.Data[0].(map[string]interface{})
	if !ok {
		ep.Socket.Emit("error", "Failed to parse data")
		return
	}

	var request = &dto.GetChatRequest{}
	chatId, err := strconv.ParseUint(fmt.Sprintf("%v", dataMap["chat_id"]), 10, 64)
	if err != nil {
		ep.Socket.Emit("error", "Chat ID must be a number")
		return
	}

	request.ChatID = uint(chatId)
	request.UserID = uint(user.GetId())

	// Validate request
	if errors := helper.Validate(u.Validate, request); len(errors) > 0 {
		message := fmt.Sprintf("The given data was invalid: %v", errors)
		ep.Socket.Emit("error", message)
		return
	}

	// Get response from chat service
	getChatCtx, cancelGetChat := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancelGetChat()

	_, err = u.Chat.Service.GetChatByIdAndUserId(getChatCtx, &proto.GetChatRequest{
		ChatId: int64(request.ChatID),
		UserId: user.GetId(),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.InvalidArgument {
				ep.Socket.Emit("error", s.Message())
				return
			} 
			if s.Code() == codes.NotFound {
				ep.Socket.Emit("error", s.Message())
				return
			} 
			ep.Socket.Emit("error", s.Message())
			return
		}

		ep.Socket.Emit("error", "Failed to get chat")
		return
	}

	// Join chat room
	room := fmt.Sprintf("chat_%d", request.ChatID)
	ep.Socket.Join(room)
	ep.Socket.Emit("message", fmt.Sprintf("Joined chat room %d", request.ChatID))
}

func (u *ChatUseCaseImpl) CreateMessage(ep *socketio.EventPayload) {
	// Get user
	user, err := u.getUser(ep)
	if err != nil {
		ep.Socket.Emit("error", "Unauthorized")
		u.Log.Errorf("Failed to get user: %v", err)
		return
	}

	// Parse data
	dataMap, ok := ep.Data[0].(map[string]interface{})
	if !ok {
		ep.Socket.Emit("error", "Failed to parse data")
		return
	}

	var request = &dto.CreateMessageRequest{}
	chatId, err := strconv.ParseUint(fmt.Sprintf("%v", dataMap["chat_id"]), 10, 64)
	if err != nil {
		ep.Socket.Emit("error", "Chat ID must be a number")
		return
	}

	message, ok := dataMap["message"].(string)
	if !ok {
		ep.Socket.Emit("error", "Message must be a string")
		return
	}

	request.ChatID = uint(chatId)
	request.Message = message

	// Validate request
	if errors := helper.Validate(u.Validate, request); len(errors) > 0 {
		message := fmt.Sprintf("The given data was invalid: %v", errors)
		ep.Socket.Emit("error", message)
		return
	}

	// Craete request message to chat service
	stream, err := u.Chat.Service.CreateMessage(context.Background())
	if err != nil {
		ep.Socket.Emit("error", "Failed to create message")
		u.Log.Errorf("Failed to create message: %v", err)
		return
	}

	// Send message to chat service
	stream.Send(&proto.Message{
		ChatId: int64(request.ChatID),
		UserId: user.GetId(),
		Message: request.Message,
	})

	// Close send
	if err := stream.CloseSend(); err != nil {
		ep.Socket.Emit("error", "Failed to close send")
		u.Log.Errorf("Failed to close send: %v", err)
		return
	}

	if request.ChatID != 0 {
		room := fmt.Sprintf("chat_%d", request.ChatID)
		ep.Socket.Join(room)
	}

	// Loop get response from chat service
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		part := &dto.MessagePart{
			ChatID: uint(res.GetChatId()),
			Part:   res.GetPart(),
			Status: res.GetStatus(),
		}

		if err != nil {
			ep.Socket.Emit("error", "Failed to receive message")
			u.Log.Errorf("Failed to receive message: %v", err)
			break
		}

		// Marshal response to bytes
		resJson, err := json.Marshal(part)
		if err != nil {
			ep.Socket.Emit("error", "Failed to marshal response")
			u.Log.Errorf("Failed to marshal response: %v", err)
			break
		}

		// If response status is start, emit start event to user
		if res.GetStatus() == proto.Status_START {
			ep.Socket.Emit("message_start", string(resJson))
			continue
		}

		// If response status is END, emit end event and break the loop
		if res.GetStatus() == proto.Status_END {
			ep.Socket.Emit("message_end", string(resJson))
			break
		}

		// Send message to user
		if request.ChatID != 0 {
			event := fmt.Sprintf("message_chat_%d", request.ChatID)
			ep.Socket.Emit(event, string(resJson))
		} else {
			ep.Socket.Emit("message", string(resJson))
		}
	}
}
