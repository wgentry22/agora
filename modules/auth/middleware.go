package auth

import (
  "errors"
  "net/http"
  "sync"

  "github.com/gin-gonic/gin"
  "github.com/wgentry22/agora/types/config"
)

var (
  m sync.Mutex
  validator TokenValidator
)

func Use(conf config.Auth) {
  m.Lock()
  defer m.Unlock()

  if conf.Vendor.String() == "firebase" {
    validator = newFirebaseTokenValidator()
  }
}

func RequiresTokenMiddleware(c *gin.Context) {
  if validator == nil {
    panic(errors.New("cannot use `RequiresTokenMiddleware` without a `TokenValidator`"))
  }

  sub, err := validator.Validate(c.Request)
  if err != nil {
    c.AbortWithStatus(http.StatusUnauthorized)
  } else {
    c.Set("subject", sub)
    c.Next()
  }
}
