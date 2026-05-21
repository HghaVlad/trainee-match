package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
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

	var userID string
	err := token.Get("sub", &userID)
	if err != nil {
		return User{}, err
	}
	user.Id, err = uuid.Parse(userID)
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

	role, err := getRole(token)
	if err != nil {
		return User{}, err
	}
	user.Role = role

	return user, nil
}

func getRole(token jwt.Token) (string, error) {
	var realmAccess map[string]any

	err := token.Get("realm_access", &realmAccess)
	if err != nil {
		return "", err
	}

	if rolesRaw, exists := realmAccess["roles"]; exists {
		if rolesList, ok := rolesRaw.([]any); ok {
			for _, r := range rolesList {
				if r == "Candidate" {
					return "Candidate", nil
				}
				if r == "admin" {
					return "admin", nil
				}
			}
		}
	}
	return "", fmt.Errorf("role not found")
}
