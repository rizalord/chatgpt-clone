package dto

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password"`
    ImageURL string `json:"image_url"`
}

type UpdateUserRequest struct {
    ID       uint   `json:"id" validate:"required"`
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password"`
    ImageURL string `json:"image_url"`
}