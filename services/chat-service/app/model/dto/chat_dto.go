package dto

type GetChatsRequest struct {
    UserID uint `json:"user_id" validate:"required"`
}

type GetChatRequest struct {
    ChatID uint `json:"chat_id" validate:"required"`
    UserID uint `json:"user_id" validate:"required"`
}