package http

import (
	"api-gateway/app/model/dto"
	"api-gateway/app/usecase"

	socketio "github.com/doquangtan/gofiber-socket.io"
	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
    ChatUseCase usecase.ChatUseCase
}

func NewChatController(chatUseCase usecase.ChatUseCase) *ChatController {
    return &ChatController{
        ChatUseCase: chatUseCase,
    }
}

func (c *ChatController) GetChats(ctx *fiber.Ctx) error {
    userId, ok := ctx.Locals("userId").(int64)
    if !ok {
        return fiber.ErrUnauthorized
    }
    userIdUint := uint(userId)

    response, err := c.ChatUseCase.GetChats(ctx.UserContext(), &userIdUint)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.GetChatsResponse]{
        Message: "Chats retrieved",
        Data:    response,
    })
}

func (c *ChatController) GetMessages(ctx *fiber.Ctx) error {
    id := ctx.Params("chat_id")
    userId, ok := ctx.Locals("userId").(int64)
    if !ok {
        return fiber.ErrUnauthorized
    }
    userIdUint := uint(userId)

    response, err := c.ChatUseCase.GetMessages(ctx.UserContext(), &id, &userIdUint)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.GetMessagesResponse]{
        Message: "Messages retrieved",
        Data:    response,
    })
}

func (c *ChatController) JoinChatRoom(ep *socketio.EventPayload) {
    c.ChatUseCase.JoinChatRoom(ep)
}

func (c *ChatController) CreateMessage(ep *socketio.EventPayload) {
    c.ChatUseCase.CreateMessage(ep)
}