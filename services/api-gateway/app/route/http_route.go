package route

import (
	"api-gateway/app/delivery/http"
	"api-gateway/app/delivery/http/middleware"

	socketio "github.com/doquangtan/gofiber-socket.io"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type HttpRouter struct {
	Router *fiber.App
}

func NewHttpRouter(
	router *fiber.App, 
	auth *middleware.AuthMiddleware,
	websocket *middleware.WebsocketMiddleware,
	authController *http.AuthController,
	chatController *http.ChatController,
) *HttpRouter {
	// Middleware
	logFormat := `{"time": "${time}", "status": "${status}", "latency": "${latency}", "ip": "${ip}", "method": "${method}", "path": "${path}", "error": "${error}"}` + "\n"
	router.Use(logger.New(logger.Config{Format: logFormat}))
	router.Use(cors.New())

	// Routes
	router.Post("/auth/register", authController.Register)
	router.Post("/auth/login/google", authController.LoginWithGoogle)
	router.Post("/auth/login", authController.Login)
	router.Post("/auth/refresh", authController.RefreshToken)

	router.Get("/chats", auth.Handler, chatController.GetChats)
	router.Get("/chats/:chat_id/messages", auth.Handler, chatController.GetMessages)

	// Websocket
	socketIoRoute(router, websocket, chatController)

	http := &HttpRouter{
		Router: router,
	}

	return http
}

func socketIoRoute(router *fiber.App, middleware *middleware.WebsocketMiddleware, chatController *http.ChatController) {
	io := socketio.New()

	io.OnAuthorization(middleware.Handler)

	io.OnConnection(func(socket *socketio.Socket) {
		// middleware.Handler(socket)
		socket.On("join_chat_room", chatController.JoinChatRoom)
		socket.On("create_message", chatController.CreateMessage)
	})

	router.Use("/", io.Middleware)
	router.Route("/socket.io", io.Server)
}