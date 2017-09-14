package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenClaims defines a common structure used for JWT claims in rig
type TokenClaims struct {
	IssuerType string `json:"ist,omitempty"`
	jwt.StandardClaims
}

// NewUserToken generates a JWT for users and signs with given secret
func NewUserToken(userid int, secret string) (string, error) {
	claims := TokenClaims{
		"user",
		jwt.StandardClaims{
			Issuer:    string(userid),
			Audience:  "rig.xplex.me",
			ExpiresAt: time.Now().AddDate(0, 0, 28).Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// NewInviteToken generates a JWT for user invites
func NewInviteToken(senderid int, email string, secret string) (string, error) {
	claims := TokenClaims{
		"invite",
		jwt.StandardClaims{
			Issuer:    string(senderid),
			Subject:   email,
			Audience:  "rig.xplex.me",
			ExpiresAt: time.Now().AddDate(0, 0, 14).Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// ParseToken attempts to verify a signed JWT issued for user auth
func ParseToken(t string, secret string) (*TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(t, &TokenClaims{}, func(ti *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := parsed.Claims.(*TokenClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, err
}
