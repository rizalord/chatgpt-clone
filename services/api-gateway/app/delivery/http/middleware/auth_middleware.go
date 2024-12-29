package middleware

import (
	"api-gateway/app/delivery/client"
	"api-gateway/app/model/proto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	Handler fiber.Handler
}

func NewAuthMiddleware(auth *client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: func(ctx *fiber.Ctx) error {
			splitToken := strings.Split(ctx.Get("Authorization"), "Bearer ")
			if len(splitToken) < 2 {
				return fiber.ErrUnauthorized
			}

			accessToken := splitToken[1]

			user, err := auth.Service.GetProfile(ctx.Context(), &proto.GetProfileRequest{
				AccessToken: accessToken,
			})
			if err != nil {
				return fiber.ErrUnauthorized
			}

			ctx.Locals("userId", user.GetId())
			ctx.Locals("email", user.GetEmail())
			ctx.Locals("name", user.GetName())
			ctx.Locals("imageURL", user.GetImageUrl())
			ctx.Locals("accessToken", accessToken)

			return ctx.Next()
		},
	}
}