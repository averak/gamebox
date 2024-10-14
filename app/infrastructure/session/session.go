package session

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Session struct {
	PrincipalID uuid.UUID
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

func NewSession(principalID uuid.UUID, issuedAt time.Time, expiresAt time.Time) Session {
	return Session{
		PrincipalID: principalID,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
	}
}

func (s Session) IsExpired(now time.Time) bool {
	return now.Equal(s.ExpiresAt) || now.After(s.ExpiresAt)
}

func EncodeSessionToken(session Session, secretKey []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   session.PrincipalID.String(),
		IssuedAt:  jwt.NewNumericDate(session.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(session.ExpiresAt),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}

func DecodeSessionToken(token string, secretKey []byte, now time.Time) (Session, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		jwt.WithTimeFunc(func() time.Time {
			return now
		}),
	)
	if err != nil {
		return Session{}, err
	}
	return Session{
		PrincipalID: uuid.MustParse(claims.Subject),
		IssuedAt:    claims.IssuedAt.Time,
		ExpiresAt:   claims.ExpiresAt.Time,
	}, nil
}
