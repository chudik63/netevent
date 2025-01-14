package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/chudik63/netevent/auth/internal/db/postgres/models"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	secretKey   = []byte("secret key")
	ExpTime     = 15 * time.Minute
	LongExpTime = 24 * time.Hour
	jwtMethod   = jwt.SigningMethodHS256
)

func NewTokens(id int64, role string) (*models.Token, error) {
	access, err := GetToken(id, role, ExpTime)
	if err != nil {
		return nil, err
	}

	refresh, err := GetToken(id, role, LongExpTime)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		AccessTkn:  access,
		AccessTtl:  int64(ExpTime),
		RefreshTkn: refresh,
		RefreshTtl: int64(LongExpTime),
	}, nil

}

func GetToken(id int64, role string, expTime time.Duration) (string, error) {
	calms := jwt.NewWithClaims(jwtMethod, jwt.MapClaims{
		"userId":   id,
		"userRole": role,
		"exp":      time.Now().Add(expTime).Unix(),
		"iat":      time.Now().Unix(),
	})

	token, err := calms.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidToken(token string) (bool, error) {
	valid, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	switch {
	case valid.Valid:
		return true, nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return false, errors.New("expired token")
	}
	if err != nil {
		return false, err
	}
	return false, errors.New("unexpected error with token")
}

func GetIdToken(tkn string) (int64, error) {
	calms, err := jwt.Parse(tkn, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})
	if claims, ok := calms.Claims.(jwt.MapClaims); ok && calms.Valid {
		return claims["userId"].(int64), nil
	}

	return 0, err
}

func GetRoleToken(tkn string) (string, error) {
	calms, err := jwt.Parse(tkn, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})
	if claims, ok := calms.Claims.(jwt.MapClaims); ok && calms.Valid {
		return claims["userRole"].(string), nil
	}

	return "", err
}
