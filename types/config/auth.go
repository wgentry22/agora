package config

import (
  "errors"
)

type AuthVendor int8

const (
  AuthVendorUnknown AuthVendor = iota
  AuthVendorFirebase
)

var (
  ErrAuthVendorRequired = errors.New("value for `auth.vendor` is expected")
  ErrUnknownAuthVendor  = errors.New("unknown auth vendor")
  authVendorDisplay     = []string{"unknown", "firebase"}
  authVendorLookup      = map[string]AuthVendor{
    "unknown":  AuthVendorUnknown,
    "firebase": AuthVendorFirebase,
  }
)

func (a AuthVendor) String() string {
  return authVendorDisplay[a]
}

func (a *AuthVendor) UnmarshalTOML(data interface{}) error {
  if val, ok := data.(string); ok {
    found, isKnown := authVendorLookup[val]
    if isKnown {
      *a = found
    }

    return ErrUnknownAuthVendor
  }

  return ErrAuthVendorRequired
}

type Auth struct {
  Vendor AuthVendor             `toml:"vendor"`
  Args   map[string]interface{} `toml:"args"`
}
