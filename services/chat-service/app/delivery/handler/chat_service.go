package handler

import (
	"auth-service/app/model/proto"
	"auth-service/app/usecase"
	"context"
)

type ChatService interface {
    proto.ChatServiceServer
    GetChats(ctx context.Context, req *proto.GetChatsRequest) (*proto.GetChatsResponse, error)
    GetChatByIdAndUserId(ctx context.Context, req *proto.GetChatRequest) (*proto.GetChatResponse, error)
    GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error)
    CreateMessage(cmr proto.ChatService_CreateMessageServer) error
}

type ChatServiceImpl struct {
    proto.UnimplementedChatServiceServer
    chatUseCase usecase.ChatUseCase
    messageUseCase usecase.MessageUseCase
}

func NewChatServiceImpl(
    chatUseCase usecase.ChatUseCase,
    messageUseCase usecase.MessageUseCase,
) *ChatServiceImpl {
    return &ChatServiceImpl{
        chatUseCase: chatUseCase,
        messageUseCase: messageUseCase,
    }
}

func (s *ChatServiceImpl) GetChats(ctx context.Context, req *proto.GetChatsRequest) (*proto.GetChatsResponse, error) {
    chats, err := s.chatUseCase.GetChats(ctx, req)
    if err != nil {
        return nil, err
    }

    return chats, nil
}

func (s *ChatServiceImpl) GetChatByIdAndUserId(ctx context.Context, req *proto.GetChatRequest) (*proto.GetChatResponse, error) {
    chat, err := s.chatUseCase.GetChatByIdAndUserId(ctx, req)
    if err != nil {
        return nil, err
    }

    return chat, nil
}

func (s *ChatServiceImpl) GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
    messages, err := s.messageUseCase.GetMessages(ctx, req)
    if err != nil {
        return nil, err
    }

    return messages, nil
}

func (s *ChatServiceImpl) CreateMessage(cmr proto.ChatService_CreateMessageServer) error {
    return s.messageUseCase.CreateMessage(cmr)
}
