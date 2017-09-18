package token

import (
	"fmt"
	"time"

	"github.com/applait/xplex-rig/models"
	jwt "github.com/dgrijalva/jwt-go"
)

// Claims defines a common structure used for JWT claims in rig
type Claims struct {
	IssuerType string `json:"ist,omitempty"`
	jwt.StandardClaims
}

// NewUserToken generates a JWT for users and signs with given secret
func NewUserToken(u *models.User, secret string) (string, error) {
	claims := Claims{
		"user",
		jwt.StandardClaims{
			Issuer:    fmt.Sprintf("%d", u.ID),
			Subject:   u.Username,
			Audience:  "rig.xplex.me",
			ExpiresAt: time.Now().AddDate(0, 0, 28).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// NewInviteToken generates a JWT for user invites
func NewInviteToken(senderid string, email string, secret string) (string, error) {
	claims := Claims{
		"invite",
		jwt.StandardClaims{
			Issuer:    senderid,
			Subject:   email,
			Audience:  "rig.xplex.me",
			ExpiresAt: time.Now().AddDate(0, 0, 14).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// ParseToken attempts to verify a signed JWT issued for user auth
func ParseToken(t string, secret string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(t, &Claims{}, func(ti *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid && claims.VerifyAudience("rig.xplex.me", true) {
		return claims, nil
	}
	return nil, err
}
