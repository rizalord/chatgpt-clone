package dto

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginWIthGoogleRequest struct {
    IdToken string `json:"id_token" validate:"required"`
}

type GetProfileRequest struct {
    AccessToken string `json:"access_token" validate:"required"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}
