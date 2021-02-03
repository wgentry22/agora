package auth

import (
  "errors"
  "net/http"
  "strings"
  "sync"

  "github.com/gin-gonic/gin"
  "github.com/wgentry22/agora/types/config"
)

var (
  ErrAuthorizationHeaderRequired = errors.New("authorization header missing")
  ErrAuthConfigurationRequired = errors.New("`RequiresTokenMiddleware` requires that you provide a config.Auth to auth.Use")
  m sync.Mutex
  validator TokenValidator
)

func Use(conf config.Auth) {
  m.Lock()
  defer m.Unlock()

  if conf.Vendor.String() == "firebase" {
    validator = newFirebaseTokenValidator()
  } else if conf.Vendor.String() == "mock" {
    validator = newMockTokenValidator()
  }
}

func RequiresTokenMiddleware(c *gin.Context) {
  if validator == nil {
    panic(ErrAuthConfigurationRequired)
  }

  sub, err := validator.Validate(c.Request)
  if err != nil {
    c.AbortWithStatus(http.StatusUnauthorized)
  } else {
    c.Set("subject", sub)
    c.Next()
  }
}

type mockTokenValidator struct {}

func (m *mockTokenValidator) Validate(r *http.Request) (string, error) {
  header := r.Header.Get("Authorization")
  if header != "" && strings.HasPrefix(header, "Bearer "){
    return "mock", nil
  }

  return "", ErrAuthorizationHeaderRequired
}

func newMockTokenValidator() TokenValidator {
  return &mockTokenValidator{}
}
