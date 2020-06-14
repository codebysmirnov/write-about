package jwt

import (
	"context"
	"errors"
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT struct
type JWT struct {
	signingKey    []byte
	defaultExpire time.Duration
}

// Default token expire time is 30 minutes
func NewJWT(opts ...Option) *JWT {
	options := newOptions(opts...)

	if len(options.signingKey) <= 0 {
		panic("Empty JWT key")
	}

	obj := JWT(options)
	return &obj
}

// Middleware check user auth
// TODO: Brake this method to validate() and middleware()
func (j *JWT) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Header["Token"]; ok {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(val[0], claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("error token decryption")
				}
				return j.signingKey, nil
			})
			if err != nil {
				utils.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			ctx := context.WithValue(context.Background(), "user", auth.Meta(claims))
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
// TODO: Fix pars params
func (j *JWT) Generate(args ...interface{}) (string, error) {
	params := auth.Meta{}
	for _, arg := range args {
		switch arg.(type) {
		case auth.Meta, map[string]interface{}:
			params = arg.(auth.Meta)
		default:
			return "", errors.New("token generate fail")
		}
	}

	// use custom expire time if exists
	var expire = j.defaultExpire
	if val, ok := params["expire"]; ok {
		switch val.(type) {
		case time.Duration:
			expire = val.(time.Duration)
		default:
			return "", errors.New("invalid value type of token expire duration")
		}
		delete(params, "expire")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range params {
		claims[k] = v
	}

	claims["expire"] = time.Now().Add(expire).Unix()

	tokenString, err := token.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
