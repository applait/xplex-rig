package account

import (
	"time"

	"github.com/applait/xplex-rig/common"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Claims defines a common structure used for JWT claims in rig
type Claims struct {
	IssuerType string `json:"ist,omitempty"`
	jwt.StandardClaims
}

// newUserToken generates a JWT for users and signs with given secret
func newUserToken(userid uuid.UUID, username string) (string, error) {
	claims := Claims{
		"user",
		jwt.StandardClaims{
			Issuer:    userid.String(),
			Subject:   username,
			Audience:  "rig.xplex.me",
			ExpiresAt: time.Now().AddDate(0, 0, 28).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return utoken.SignedString([]byte(common.Config.JWTKeys.Users))
}

// NewInviteToken generates a JWT for user invites
func NewInviteToken(senderid string, email string) (string, error) {
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
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return utoken.SignedString([]byte(common.Config.JWTKeys.Users))
}

// ParseUserToken attempts to verify a signed JWT issued for user auth
func ParseUserToken(t string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(t, &Claims{}, func(ti *jwt.Token) (interface{}, error) {
		return []byte(common.Config.JWTKeys.Users), nil
	})
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid && claims.VerifyAudience("rig.xplex.me", true) {
		return claims, nil
	}
	return nil, err
}
