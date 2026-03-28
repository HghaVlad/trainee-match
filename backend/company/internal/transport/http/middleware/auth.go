package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type AuthMiddleware struct {
	JWKUrl string
	keys   jwk.Set
}

func NewAuthMiddleware(conf *config.Config) (*AuthMiddleware, error) {
	m := &AuthMiddleware{
		JWKUrl: conf.HTTP.JWKUrl,
	}

	err := m.getPublicKey()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if cookies == nil {
			helpers.RespondErrorMsg(w, http.StatusUnauthorized, "missing cookies")
			return
		}

		tokenString := getAccessTokenFromCookies(cookies)
		if tokenString == "" {
			helpers.RespondErrorMsg(w, http.StatusUnauthorized, "couldn't get jwt from cookies")
			return
		}

		var claims CustomClaims
		token, err := jwt.ParseString(
			tokenString,
			jwt.WithKeySet(m.keys),
			jwt.WithValidate(true),
			jwt.WithTypedClaim("realm_access", &claims.RealmAccess),
		)
		if err != nil {
			helpers.RespondError(w, http.StatusUnauthorized, err)
			return
		}

		ident, err := getIdentityFromToken(token, &claims)
		if err != nil {
			helpers.RespondError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), identityKey, ident)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) getPublicKey() error {
	keys, err := jwk.Fetch(context.Background(), m.JWKUrl)
	if err != nil {
		return err
	}

	m.keys = keys
	return nil
}

type ctxIdentityKeyT struct{}

//nolint:gochecknoglobals // ctx key
var identityKey = ctxIdentityKeyT{}

func IdentityFromContext(ctx context.Context) *identity.Identity {
	id, ok := ctx.Value(identityKey).(*identity.Identity)
	if !ok {
		panic("identity not found in context: auth middleware is not applied")
	}
	return id
}

func getAccessTokenFromCookies(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value
		}
	}
	return ""
}

func getIdentityFromToken(token jwt.Token, _ *CustomClaims) (*identity.Identity, error) {
	ident := new(identity.Identity)

	sub, ok := token.Subject()
	if !ok {
		return nil, errors.New("invalid jwt: sub not found")
	}

	subID, err := uuid.Parse(sub)
	if err != nil {
		return nil, errors.New("invalid jwt: sub was expected to be uuid format")
	}
	ident.UserID = subID

	var realmAccess RealmAccess
	if err := token.Get("realm_access", &realmAccess); err != nil {
		return nil, errors.New("invalid jwt: realm_access invalid")
	}

	for _, role := range realmAccess.Roles {
		grole := identity.GlobalRole(role)

		switch grole {
		case identity.RoleCandidate,
			identity.RoleHR,
			identity.RoleAdmin:
			ident.Role = grole
		}
	}

	if ident.Role == "" {
		return nil, errors.New("invalid jwt: no valid role was found")
	}

	return ident, nil
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type CustomClaims struct {
	RealmAccess RealmAccess `json:"realm_access"`
}
