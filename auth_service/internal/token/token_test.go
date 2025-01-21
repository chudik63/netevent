package token

import (
	"testing"
	"time"

	"github.com/chudik63/netevent/auth_service/internal/models"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidToken(t *testing.T) {
	testTable := []struct {
		name        string
		id          int64
		role        string
		expTime     time.Duration
		secretKey   []byte
		expectedErr error
	}{
		{
			name:        "OK",
			id:          1,
			role:        "user",
			expTime:     100 * time.Second,
			secretKey:   []byte("secret key"),
			expectedErr: nil,
		},
		{
			name:        "Wrong secret key",
			id:          1,
			role:        "user",
			expTime:     100 * time.Second,
			secretKey:   []byte("WRONG KEY"),
			expectedErr: models.ErrSignatureInvalid,
		},
		{
			name:        "Token is expired",
			id:          1,
			role:        "user",
			expTime:     1 * time.Second,
			secretKey:   []byte("secret key"),
			expectedErr: models.ErrTokenExpired,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":   testCase.id,
				"userRole": testCase.role,
				"exp":      time.Now().Add(testCase.expTime).Unix(),
				"iat":      time.Now().Unix(),
			})

			token, err := tkn.SignedString(testCase.secretKey)

			time.Sleep(1 * time.Second)

			err = ValidToken(token)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}

func TestGetIdToken(t *testing.T) {
	testTable := []struct {
		name        string
		id          interface{}
		expTime     time.Duration
		expectedId  int64
		expectedErr error
	}{
		{
			name:        "OK",
			id:          1,
			expTime:     100 * time.Second,
			expectedId:  1,
			expectedErr: nil,
		},
		{
			name:        "Wrong id type",
			id:          "aba",
			expTime:     100 * time.Second,
			expectedId:  0,
			expectedErr: models.ErrGetFromClaims,
		},
		{
			name:        "No id in claims",
			expTime:     100 * time.Second,
			expectedId:  0,
			expectedErr: models.ErrGetFromClaims,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"exp": time.Now().Add(testCase.expTime).Unix(),
				"iat": time.Now().Unix(),
			}

			if testCase.id != nil {
				claims["userId"] = testCase.id
			}

			tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			token, err := tkn.SignedString([]byte("secret key"))

			id, err := GetIdToken(token)
			assert.Equal(t, id, testCase.expectedId)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}

func TestGetRoleToken(t *testing.T) {
	testTable := []struct {
		name         string
		role         interface{}
		expTime      time.Duration
		expectedRole string
		expectedErr  error
	}{
		{
			name:         "OK",
			role:         "user",
			expTime:      100 * time.Second,
			expectedRole: "user",
			expectedErr:  nil,
		},
		{
			name:         "Wrong role type",
			role:         1,
			expTime:      100 * time.Second,
			expectedRole: "",
			expectedErr:  models.ErrGetFromClaims,
		},
		{
			name:         "No id in claims",
			expTime:      100 * time.Second,
			expectedRole: "",
			expectedErr:  models.ErrGetFromClaims,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"exp": time.Now().Add(testCase.expTime).Unix(),
				"iat": time.Now().Unix(),
			}

			if testCase.role != nil {
				claims["userRole"] = testCase.role
			}

			tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			token, err := tkn.SignedString([]byte("secret key"))

			id, err := GetRoleToken(token)
			assert.Equal(t, id, testCase.expectedRole)
			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}
