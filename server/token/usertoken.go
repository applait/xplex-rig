package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type customClaim struct {
	IssuerType string `json:"ist,omitempty"`
	jwt.StandardClaims
}

// UserAuthToken generates a JWT for users and signs with given secret
func UserAuthToken(userid int, secret string) (string, error) {
	claims := customClaim{
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

// InviteToken generates a JWT for user invites
func InviteToken(senderid int, email string, secret string) (string, error) {
	claims := customClaim{
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
