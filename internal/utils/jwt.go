package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpired      = errors.New("token expired")
	ErrRevoked      = errors.New("token revoked")
	ErrNotFound     = errors.New("token not found")
)

type Claims struct {
	UserID int64
	JTI    string
	Exp    time.Time
}

type JWT interface {
	SignAccess(userID int64, jti string, now time.Time) (string, time.Time, error)
	ParseAndVerify(raw string, now time.Time) (*Claims, error)
}

type HMACJWT struct {
	Issuer     string
	AccessTTL  time.Duration
	SigningKey []byte
}

func (j HMACJWT) SignAccess(userID int64, jti string, now time.Time) (string, time.Time, error) {
	exp := now.Add(j.AccessTTL)

	rc := jwt.RegisteredClaims{
		Issuer:    j.Issuer,
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        jti,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	raw, err := t.SignedString(j.SigningKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return raw, exp, nil
}

func (j HMACJWT) ParseAndVerify(raw string, now time.Time) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(
		raw,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return j.SigningKey, nil
		},
		jwt.WithIssuer(j.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, ErrInvalidToken
	}

	rc, ok := parsed.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	if rc.Subject == "" || rc.ID == "" || rc.ExpiresAt == nil {
		return nil, ErrInvalidToken
	}

	uid, err := strconv.ParseInt(rc.Subject, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Claims{
		UserID: uid,
		JTI:    rc.ID,
		Exp:    rc.ExpiresAt.Time,
	}, nil
}
