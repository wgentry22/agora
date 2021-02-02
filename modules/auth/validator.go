package auth

import (
  "context"
  "errors"
  "fmt"
  "net/http"
  "os"
  "strings"

  firebase "firebase.google.com/go"
  "firebase.google.com/go/auth"
  "github.com/hashicorp/errwrap"
  "google.golang.org/api/option"
)

var (
  ErrFailedToCreateFirebaseTokenValidator = errors.New("failed to create FirebaseTokenValidator - `GOOGLE_APPLICATION_CREDENTIALS` environment variable is missing")
  ErrFailedToConfigureFirebase            = errors.New("failed to configure Firebase Authentication")
)

type TokenValidator interface {
  Validate(r *http.Request) (string, error)
}

type FirebaseTokenValidator struct {
  client *auth.Client
}

func newFirebaseTokenValidator(opts ...option.ClientOption) TokenValidator {
  app, err := firebase.NewApp(context.Background(), nil, opts...)
  if err != nil {
    return firebaseTokenValidatorFromEnv()
  }

  client, err := app.Auth(context.Background())
  if err != nil {
    panic(errwrap.Wrap(ErrFailedToConfigureFirebase, err))
  }

  return &FirebaseTokenValidator{client}
}

func firebaseTokenValidatorFromEnv() TokenValidator {
  _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
  if !ok {
    panic(ErrFailedToCreateFirebaseTokenValidator)
  }

  app, err := firebase.NewApp(context.Background(), nil)
  if err != nil {
    panic(errwrap.Wrap(ErrFailedToConfigureFirebase, err))
  }

  client, err := app.Auth(context.Background())

  if err != nil {
    panic(errwrap.Wrap(ErrFailedToConfigureFirebase, err))
  }

  return &FirebaseTokenValidator{client}
}

func (f *FirebaseTokenValidator) Validate(r *http.Request) (string, error) {
  header := r.Header.Get("Authorization")
  if header == "" || !strings.HasPrefix(header, "Bearer ") {
    return "", errors.New("authorization header missing")
  }

  jwt := strings.ReplaceAll(header, "Bearer ", "")

  token, err := f.client.VerifyIDToken(r.Context(), jwt)
  if err != nil {
    return "", err
  }

  fmt.Printf("Verified FirebaseToken: %+v\n", *token)

  return token.Subject, nil
}
