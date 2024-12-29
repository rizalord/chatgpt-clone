package middleware

import (
	"api-gateway/app/delivery/client"
	"api-gateway/app/model/proto"
	"context"
	// socketio "github.com/doquangtan/gofiber-socket.io"
)

type WebsocketMiddleware struct {
	// Handler func(socket *socketio.Socket)
	Handler func(params map[string]string) bool
}

func NewWebsocketMiddleware(auth *client.AuthClient) *WebsocketMiddleware {
	return &WebsocketMiddleware{
		// Handler: func (socket *socketio.Socket) {
		// 	// Validate token
		// 	accessToken := socket.Conn.Headers("Authorization", "Bearer ")
		// 	splitToken := strings.Split(accessToken, "Bearer ")
		// 	if len(splitToken) < 2 {
		// 		socket.Emit("error", "Unauthorized")
		// 		socket.Disconnect()
		// 		return 
		// 	}

		// 	accessToken = splitToken[1]

		// 	_, err := auth.Service.GetProfile(context.Background(), &proto.GetProfileRequest{
		// 		AccessToken: accessToken,
		// 	})
		// 	if err != nil {
		// 		socket.Emit("error", "Unauthorized")
		// 		socket.Disconnect()
		// 		return
		// 	}
		// },
		Handler: func(params map[string]string) bool {
			accessToken := params["token"]

			_, err := auth.Service.GetProfile(context.Background(), &proto.GetProfileRequest{
				AccessToken: accessToken,
			})
			return err == nil
		},
	}
}