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

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

type CredentialData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    ImageURL *string `json:"image_url,omitempty"`
}

type TokenData struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiredAt int64 `json:"expired_at"`
}

type LoginResponse struct {
    User  CredentialData `json:"user"`
    Token TokenData      `json:"token"`
}
