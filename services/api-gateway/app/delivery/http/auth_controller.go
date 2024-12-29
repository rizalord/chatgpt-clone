package http

import (
	"api-gateway/app/model/dto"
	"api-gateway/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
    AuthUseCase usecase.AuthUseCase
}

func NewAuthController(authUseCase usecase.AuthUseCase) *AuthController {
    return &AuthController{
        AuthUseCase: authUseCase,
    }
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
    request := new(dto.RegisterRequest)
    if err := ctx.BodyParser(request); err != nil {
        return fiber.ErrBadRequest
    }

    response, err := c.AuthUseCase.Register(ctx.UserContext(), request)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.LoginResponse]{
        Message: "Register successful",
        Data:    response,
    })
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
    request := new(dto.LoginRequest)
    if err := ctx.BodyParser(request); err != nil {
        return fiber.ErrBadRequest
    }

    response, err := c.AuthUseCase.Login(ctx.UserContext(), request)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.LoginResponse]{
        Message: "Login successful",
        Data:    response,
    })
}

func (c *AuthController) LoginWithGoogle(ctx *fiber.Ctx) error {
    request := new(dto.LoginWIthGoogleRequest)
    if err := ctx.BodyParser(request); err != nil {
        return fiber.ErrBadRequest
    }

    response, err := c.AuthUseCase.LoginWithGoogle(ctx.UserContext(), request)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.LoginResponse]{
        Message: "Login successful",
        Data:    response,
    })
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
    request := new(dto.RefreshRequest)
    if err := ctx.BodyParser(request); err != nil {
        return fiber.ErrBadRequest
    }

    response, err := c.AuthUseCase.Refresh(ctx.UserContext(), request)
    if err != nil {
        return err
    }

    return ctx.JSON(dto.Response[dto.LoginResponse]{
        Message: "Token refreshed",
        Data:    response,
    })
}