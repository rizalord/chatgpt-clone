export interface RegisterRequest {
    name: string
    email: string
    password: string
}

export interface LoginRequest {
    email: string
    password: string
}

export interface LoginWithGoogleRequest {
    id_token: string
}

export interface AuthResponse {
    user: {
        name: string
        email: string
        image_url: string
    }
    token: {
        access_token: string
        refresh_token: string
        expired_at: number
    }
}