package helper

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/idtoken"
)

type AuthCredential struct {
	Email string
	Name  string
	Picture string
}

type OAuth interface {
	VerifyIdToken(token string) (*AuthCredential, error)
}

type OAuthImpl struct {
	Logger *logrus.Logger
	Config *viper.Viper
}

func NewOauthImpl(logger *logrus.Logger, config *viper.Viper) *OAuthImpl {
	return &OAuthImpl{
		Logger: logger,
		Config: config,
	}
}

func (o *OAuthImpl) VerifyIdToken(token string) (*AuthCredential, error) {
	googleClientId := o.Config.GetString("GOOGLE_CLIENT_ID")

	payload, err := idtoken.Validate(context.Background(), token, googleClientId)
	if err != nil {
		o.Logger.Error(err)
		return nil, err
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)

	return &AuthCredential{
		Email: email,
		Name:  name,
		Picture: picture,
	}, nil
}
