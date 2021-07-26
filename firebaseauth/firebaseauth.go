// Mostly copied from encore.dev docs
package firebaseauth

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/rlog"
	firebase "firebase.google.com/go/v4"
	fbauth "firebase.google.com/go/v4/auth"
	"go4.org/syncutil"
	"google.golang.org/api/option"
)

type AuthData struct {
	Email string
	Name string
}

//encore:authhandler
func ValidateToken(ctx context.Context, token string) (auth.UID, *AuthData, error) {
	if err := setupFB(); err != nil {
		return "", nil, err
	}
	tok, err := fbAuth.VerifyIDToken(ctx, token)
	if err != nil {
		return "", nil, err
	}

	email, _ := tok.Claims["email"].(string)
	name, _ := tok.Claims["name"].(string)
	uid := auth.UID(tok.UID)

	usr := &AuthData{
		Email:   email,
		Name:    name,
	}
	rlog.Info("auth name", "msg", string(name))
	return uid, usr, nil
}

var (
	fbAuth    *fbauth.Client
	setupOnce syncutil.Once
)

func setupFB() error {
	return setupOnce.Do(func() error {
		opt := option.WithCredentialsJSON([]byte(secrets.FirebasePrivateKey))
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err == nil {
			fbAuth, err = app.Auth(context.Background())
		}
		return err
	})
}

var secrets struct {
	FirebasePrivateKey string
}
