package dto

import "api-gateway/app/model/proto"

type GetChatsRequest struct {
    UserID uint `json:"user_id" validate:"required"`
}

type GetChatRequest struct {
    ChatID uint `json:"chat_id" validate:"required"`
    UserID uint `json:"user_id" validate:"required"`
}

type CreateMessageRequest struct {
    ChatID uint `json:"chat_id" validate:"omitempty"`
    Message string `json:"message" validate:"required"`
}

type GetMessagesRequest struct {
    ChatID uint `json:"chat_id" validate:"required"`
    UserID uint `json:"user_id" validate:"required"`
}

type ChatData struct {
    ID uint `json:"id"`
    UserID uint `json:"user_id"`
    Topic string `json:"topic"`
} 

type GetChatsResponse struct {
    Chats []ChatData `json:"chats"`
}

type MessageData struct {
    ID uint `json:"id"`
    ChatID uint `json:"chat_id"`
    UserID uint `json:"user_id"`
    Role string `json:"role"`
    Content string `json:"content"`
}

type GetMessagesResponse struct {
    Messages []MessageData `json:"messages"`
}

type MessagePart struct {
    ChatID uint `json:"chat_id"`
    Part string `json:"part"`
    Status proto.Status `json:"status"`
}
