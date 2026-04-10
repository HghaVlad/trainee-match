package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	utilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type AuthMiddleware struct {
	JWKUrl string
	keys   jwk.Set
}

func NewAuthMiddleware(ctx context.Context, conf *config.Config) (*AuthMiddleware, error) {
	m := &AuthMiddleware{
		JWKUrl: conf.HTTP.JWKUrl,
	}

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	err := m.getPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := utilslog.FromContext(ctx)

		cookies := r.Cookies()
		if cookies == nil {
			logger.InfoContext(ctx, "http request unauthorized: missing cookies",
				"status", http.StatusUnauthorized)

			helpers.RespondErrorMsg(ctx, w, http.StatusUnauthorized, "missing cookies")
			return
		}

		tokenString := getAccessTokenFromCookies(cookies)
		if tokenString == "" {
			logger.InfoContext(ctx, "http request unauthorized: couldn't get jwt from cookies",
				"status", http.StatusUnauthorized)

			helpers.RespondErrorMsg(ctx, w, http.StatusUnauthorized, "couldn't get jwt from cookies")
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
			logger.InfoContext(ctx, "http request unauthorized: invalid keys or realm",
				"status", http.StatusUnauthorized)

			helpers.RespondError(ctx, w, http.StatusUnauthorized, err)
			return
		}

		ident, err := getIdentityFromToken(token, &claims)
		if err != nil {
			logger.InfoContext(ctx, "http request unauthorized",
				"status", http.StatusUnauthorized, "err", err)

			helpers.RespondError(ctx, w, http.StatusUnauthorized, err)
			return
		}

		logger = logger.With(
			"user_id", ident.UserID,
			"role", ident.Role,
		)

		ctx = context.WithValue(ctx, identityKey, ident)
		ctx = utilslog.WithLoggerContext(ctx, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) getPublicKey(ctx context.Context) error {
	keys, err := jwk.Fetch(ctx, m.JWKUrl)
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

	bestRole, err := getBestRole(realmAccess.Roles)
	if err != nil {
		return nil, err
	}

	ident.Role = bestRole
	return ident, nil
}

// in case there are multiple roles in claims, we take
// Admin > HR > Candidate
func getBestRole(roles []string) (identity.GlobalRole, error) {
	var rolePriority = map[identity.GlobalRole]int{
		identity.RoleCandidate: 1,
		identity.RoleHR:        2,
		identity.RoleAdmin:     3,
	}

	var bestRole identity.GlobalRole
	var bestPriority int

	for _, role := range roles {
		grole := identity.GlobalRole(role)

		p, ok := rolePriority[grole]
		if !ok {
			continue
		}

		if p > bestPriority {
			bestPriority = p
			bestRole = grole
		}
	}

	if bestRole == "" {
		return "", errors.New("no valid role")
	}

	return bestRole, nil
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type CustomClaims struct {
	RealmAccess RealmAccess `json:"realm_access"`
}
