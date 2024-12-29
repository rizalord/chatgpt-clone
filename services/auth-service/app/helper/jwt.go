package helper

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTHelper interface {
    GenerateTokens(userID uint, name, email string, imageURL string) (string, string, int64, error)
    ValidateAccessToken(tokenString string) (uint, string, string, string, error)
    ValidateRefreshToken(tokenString string) (uint, string, string, string, error)
}

type JWTHelperImpl struct {
    Config *viper.Viper
}

func NewJWTHelperImpl(config *viper.Viper) *JWTHelperImpl {
    return &JWTHelperImpl{
        Config: config,
    }
}

func (a *JWTHelperImpl) GenerateTokens(userID uint, name, email string, imageURL string) (string, string, int64, error) {
    iss := a.Config.GetString("HOSTNAME")
    aud := strings.Split(a.Config.GetString("AUDIENCES"), ",")
    exp := time.Now().Add(time.Hour * 2).Unix()
    expRefresh := time.Now().Add(time.Hour * 24).Unix()
    sub := userID

    accessKey := a.Config.GetString("JWT_ACCESS_KEY")
    refreshKey := a.Config.GetString("JWT_REFRESH_KEY")

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "iss": iss,
            "aud": aud,
            "exp": exp,
            "sub": sub,
            "name": name,
            "email": email,
            "imageURL": imageURL,
        })

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "iss": iss,
            "aud": aud,
            "exp": expRefresh,
            "sub": sub,
            "name": name,
            "email": email,
            "imageURL": imageURL,
        })

    signedAccessToken, err := accessToken.SignedString([]byte(accessKey))
    if err != nil {
        return "", "", 0, err
    }

    signedRefreshToken, err := refreshToken.SignedString([]byte(refreshKey))
    if err != nil {
        return "", "", 0, err
    }

    expiredAt := time.Now().Add(time.Hour * 2).Unix()

    return signedAccessToken, signedRefreshToken, expiredAt, nil
}

func (a *JWTHelperImpl) ValidateAccessToken(tokenString string) (uint, string, string, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrInvalidKey
        }
        return []byte(a.Config.GetString("JWT_ACCESS_KEY")), nil
    })

    if err != nil {
        return 0, "", "", "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["sub"].(float64))
        name := claims["name"].(string)
        email := claims["email"].(string)
        imageURL := claims["imageURL"].(string)
        return userID, name, email, imageURL, nil
    }

    return 0, "", "", "", jwt.ErrInvalidKey
}

func (a *JWTHelperImpl) ValidateRefreshToken(tokenString string) (uint, string, string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(a.Config.GetString("JWT_REFRESH_KEY")), nil
	})

	if err != nil {
		return 0, "", "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["sub"].(float64))
        name := claims["name"].(string)
        email := claims["email"].(string)
        imageURL := claims["imageURL"].(string)
		return userID, name, email, imageURL, nil
	}

	return 0, "", "", "", jwt.ErrInvalidKey
}