package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/codebysmirnov/write-about/app/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT struct
type JWT struct {
	signingKey []byte
}

// NewJWT use secret key for generate token
func NewJWT(key string) *JWT {
	if len(key) <= 0 {
		panic("Empty jwt key")
	}
	return &JWT{signingKey: []byte(key)}
}

// Middleware check user auth
func (j *JWT) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Header["Token"]; ok {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(val[0], claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an arror")
				}
				return j.signingKey, nil
			})
			if err != nil {
				utils.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			ctx := context.WithValue(context.Background(), "user", claims)
			if token.Valid {
				handler.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			utils.RespondError(w, http.StatusUnauthorized, "Not authorized")
		}
	})
}

// Validate token
func (j *JWT) Validate(token string) (bool, error) {
	return true, nil
}

// Generate token for user auth
func (j *JWT) Generate(args ...interface{}) (string, error) {
	params := map[string]interface{}{}
	for _, arg := range args {
		switch arg.(type) {
		case map[string]interface{}:
			params = arg.(map[string]interface{})
		default:
			return "", errors.New("token generate fail")
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range params {
		claims[k] = v
	}

	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
