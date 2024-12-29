package dto

type CreateMessageRequest struct {
    ChatID uint `json:"chat_id" validate:"omitempty"`
    UserID uint `json:"user_id" validate:"required"`
    Message string `json:"message" validate:"required"`
}

type GetMessagesRequest struct {
    ChatID uint `json:"chat_id" validate:"required"`
    UserID uint `json:"user_id" validate:"required"`
}