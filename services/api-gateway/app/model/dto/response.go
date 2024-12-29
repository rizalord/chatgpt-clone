package dto

type Response[T any] struct {
	Data	*T 				`json:"data,omitempty"`
	Message	string      	`json:"message"`
}

type ErrorResponse[T any] struct {
	Message	*string      	`json:"message,omitempty"`
	Errors	*T 				`json:"errors,omitempty"`
}