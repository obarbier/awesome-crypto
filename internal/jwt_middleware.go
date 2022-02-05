package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	JwtUserDetails = contextKey("userId")
)

type jwtMiddleWare struct {
}

func NewJwtMiddleWare() *jwtMiddleWare {
	return &jwtMiddleWare{}
}

func (j *jwtMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			RespondError(w, http.StatusBadRequest, errors.New("could not retrieve token information"))
			return

		}
		if token == "" {
			RespondError(w, http.StatusBadRequest, errors.New("bearer token need was not provided"))
			return
		}
		// Validation
		claims, err := validate(token)
		if err != nil {
			msg := "failed to verify bearer token"
			log.Printf("%s : %+v", msg, err)
			RespondError(w, http.StatusUnauthorized, errors.New(msg))
			return
		}

		// FIXME: clean this up
		r = r.Clone(context.WithValue(r.Context(), JwtUserDetails, claims["sub"]))
		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no JWT.
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func validate(tokenStr string) (jwt.MapClaims, error) {
	// TODO validation and parsing
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("my_secret_key"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("could not parse claims from token")
	}

}
