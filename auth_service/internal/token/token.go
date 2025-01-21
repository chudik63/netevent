package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/chudik63/netevent/auth_service/internal/models"
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
	tkn := jwt.NewWithClaims(jwtMethod, jwt.MapClaims{
		"userId":   id,
		"userRole": role,
		"exp":      time.Now().Add(expTime).Unix(),
		"iat":      time.Now().Unix(),
	})

	token, err := tkn.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidToken(token string) error {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	switch {
	case tkn.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return models.ErrTokenExpired
	case errors.Is(err, jwt.ErrSignatureInvalid):
		return models.ErrSignatureInvalid
	}

	if err != nil {
		return err
	}

	return errors.New("unexpected error with token")
}

func GetIdToken(token string) (int64, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if userID, ok := claims["userId"].(float64); ok {
			return int64(userID), nil
		}
		return 0, models.ErrGetFromClaims
	}

	return 0, errors.New("invalid token claims")
}

func GetRoleToken(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		role, ok := claims["userRole"].(string)
		if ok {
			return role, nil
		}

		return "", models.ErrGetFromClaims
	}

	return "", errors.New("invalid token claims")
}
