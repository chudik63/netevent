package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	secretKey   = []byte("secret key")
	Small       = 100
	Long        = 400
	ExpTime     = time.Duration(Small) * time.Second
	LongExpTime = time.Duration(Long) * time.Second
	jwtMethod   = jwt.SigningMethodHS256
)

func NewToken(name string) (string, error) {
	return GetToken(name, ExpTime)
}

func RefreshToken(name string) (string, error) {
	return GetToken(name, LongExpTime)
}

func GetToken(name string, expTime time.Duration) (string, error) {
	calms := jwt.NewWithClaims(jwtMethod, jwt.MapClaims{
		"sub": name,
		"exp": time.Now().Add(expTime).Unix(),
		"iat": time.Now().Unix(),
	})
	fmt.Println(calms)

	token, err := calms.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	fmt.Println(token)

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

func GetNameToken(tkn string) (string, error) {
	calms, err := jwt.Parse(tkn, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})
	if claims, ok := calms.Claims.(jwt.MapClaims); ok && calms.Valid {
		return claims["sub"].(string), nil
	}

	return "", err
}
