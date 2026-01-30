package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"net/http"
)

type Middleware struct {
	JWKUrl string
	keys   jwk.Set
}

func NewMiddleware(jwkUrl string) *Middleware {
	m := &Middleware{
		JWKUrl: jwkUrl,
	}

	err := m.getPublicKey()
	if err != nil {
		panic(err)
	}

	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if cookies == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := getAccessTokenFromCookies(cookies)
		if tokenString == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseString(tokenString, jwt.WithKeySet(m.keys), jwt.WithValidate(true))
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := getUserFromToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) getPublicKey() error {
	keys, err := jwk.Fetch(context.Background(), m.JWKUrl)
	if err != nil {
		return err
	}

	m.keys = keys
	return nil
}

func getAccessTokenFromCookies(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value
		}
	}
	return ""
}

func getUserFromToken(token jwt.Token) (User, error) {
	var user User

	var userId string
	err := token.Get("sub", &userId)
	if err != nil {
		return User{}, err
	}
	user.Id, err = uuid.Parse(userId)
	if err != nil {
		return User{}, err
	}

	err = token.Get("first_name", &user.FirstName)
	if err != nil {
		return User{}, err
	}

	err = token.Get("last_name", &user.LastName)
	if err != nil {
		return User{}, err
	}

	err = token.Get("username", &user.Username)
	if err != nil {
		return User{}, err
	}

	err = token.Get("email", &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
