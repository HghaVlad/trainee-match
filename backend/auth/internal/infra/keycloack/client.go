package keycloack

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Nerzal/gocloak/v13"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

var ErrInvalidToken = errors.New("invalid access token")

type Client struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	realm  string

	clientID     string
	clientSecret string
	adminUser    string
	adminPass    string
}

func NewClient(clientURL, realm, clientID, clientSecret, adminUser, adminPass string) *Client {
	client := gocloak.NewClient(clientURL)
	keycloakClient := &Client{
		client:       client,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		adminUser:    adminUser,
		adminPass:    adminPass,
	}
	err := keycloakClient.loginAdmin(context.Background())
	if err != nil {
		panic(err)
	}
	slog.Debug("keycloakClient successfully started")
	return keycloakClient
}

func (kc *Client) loginAdmin(ctx context.Context) error {
	newToken, err := kc.client.LoginAdmin(ctx, kc.adminUser, kc.adminPass, kc.realm)
	if err != nil {
		return err
	}
	kc.token = newToken
	return nil
}

func (kc *Client) ensureAdminTokenValid(ctx context.Context) error {
	if kc.token == nil {
		if err := kc.loginAdmin(ctx); err != nil {
			return errors.New("keycloak loginAdmin failed")
		}
	}
	istResult, err := kc.client.RetrospectToken(ctx, kc.token.AccessToken, kc.clientID, kc.clientSecret, kc.realm)
	if err != nil {
		if err = kc.loginAdmin(ctx); err != nil {
			return errors.New("keycloak loginAdmin failed")
		}
	}

	if !*istResult.Active {
		if err = kc.loginAdmin(ctx); err != nil {
			return errors.New("keycloak loginAdmin failed")
		}
	}

	return nil
}

func (kc *Client) CreateUser(ctx context.Context, user domain.User, password string) (string, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return "", err
	}

	newUser := gocloak.User{
		FirstName: gocloak.StringP(user.FirstName),
		LastName:  gocloak.StringP(user.LastName),
		Email:     gocloak.StringP(user.Email),
		Username:  gocloak.StringP(user.Username),
		Enabled:   gocloak.BoolP(true),
	}

	userID, err := kc.client.CreateUser(ctx, kc.token.AccessToken, kc.realm, newUser)
	if err != nil {
		return "", parseError(err)
	}

	if err = kc.client.SetPassword(ctx, kc.token.AccessToken, userID, kc.realm, password, false); err != nil {
		_ = kc.client.DeleteUser(ctx, kc.token.AccessToken, kc.realm, userID)
		return "", parseError(err)
	}

	if err = kc.addRole(ctx, userID, user.Role); err != nil {
		_ = kc.client.DeleteUser(ctx, kc.token.AccessToken, kc.realm, userID)
		return "", parseError(err)
	}

	return userID, nil
}

func (kc *Client) addRole(ctx context.Context, userID, roleName string) error {
	roles := make([]gocloak.Role, 1)

	switch domain.UserRole(roleName) {
	case domain.UserCandidateRole:
		roles[0] = gocloak.Role{
			ID:   gocloak.StringP("15bd1c8f-1feb-4870-9f46-a847f0742be9"),
			Name: gocloak.StringP(string(domain.UserCandidateRole)),
		}
	case domain.UserCompanyRole:
		roles[0] = gocloak.Role{
			ID:   gocloak.StringP("2e90e50e-8db4-4881-8185-05a40220f759"),
			Name: gocloak.StringP(string(domain.UserCompanyRole)),
		}
	case domain.UserAdminRole:
		return fmt.Errorf("you can not set this role")
	default:
		return fmt.Errorf("the roleName is not valid: %s", roleName)
	}

	return kc.client.AddRealmRoleToUser(ctx, kc.token.AccessToken, kc.realm, userID, roles)
}

func (kc *Client) AddAdminRole(ctx context.Context, userID string) error {
	newRoles := []gocloak.Role{
		{
			ID:   gocloak.StringP("78b787b7-ccb1-46bb-ba4c-9eb74ab59ca7"),
			Name: gocloak.StringP(string(domain.UserAdminRole)),
		},
	}
	err := kc.deleteRoles(ctx, userID)
	if err != nil {
		return err
	}

	return kc.client.AddRealmRoleToUser(ctx, kc.token.AccessToken, kc.realm, userID, newRoles)
}

func (kc *Client) deleteRoles(ctx context.Context, userID string) error {
	roles := []gocloak.Role{
		{
			ID:   gocloak.StringP("15bd1c8f-1feb-4870-9f46-a847f0742be9"),
			Name: gocloak.StringP(string(domain.UserCandidateRole)),
		},
		{
			ID:   gocloak.StringP("2e90e50e-8db4-4881-8185-05a40220f759"),
			Name: gocloak.StringP(string(domain.UserCompanyRole)),
		},
		{
			ID:   gocloak.StringP("78b787b7-ccb1-46bb-ba4c-9eb74ab59ca7"),
			Name: gocloak.StringP(string(domain.UserAdminRole)),
		},
	}

	err := kc.client.DeleteRealmRoleFromUser(ctx, kc.token.AccessToken, kc.realm, userID, roles)
	if err != nil {
		return fmt.Errorf("deleteRoles failed: %s", err.Error())
	}
	return nil
}

func (kc *Client) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	token, err := kc.client.Login(ctx, kc.clientID, kc.clientSecret, kc.realm, username, password)

	return token, parseError(err)
}

func (kc *Client) Logout(ctx context.Context, token string) error {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return err
	}

	err := kc.client.Logout(ctx, kc.clientID, kc.clientSecret, kc.realm, token)
	return parseError(err)
}

func (kc *Client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	if err := kc.validateToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	token, err := kc.client.RefreshToken(ctx, refreshToken, kc.clientID, kc.clientSecret, kc.realm)
	return token, parseError(err)
}

func (kc *Client) validateToken(ctx context.Context, token string) error {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return err
	}

	istResult, err := kc.client.RetrospectToken(ctx, token, kc.clientID, kc.clientSecret, kc.realm)
	if err != nil {
		return ErrInvalidToken
	}
	if !(*istResult.Active) {
		return ErrInvalidToken
	}

	return nil
}

func (kc *Client) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	userInfo, err := kc.client.GetUserInfo(ctx, token, kc.realm)
	if err != nil {
		return nil, parseError(err)
	}

	user, err := kc.client.GetUserByID(ctx, kc.token.AccessToken, kc.realm, *userInfo.Sub)
	if err != nil {
		return nil, parseError(err)
	}

	return &domain.User{
		ID:        *userInfo.Sub,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Username:  *user.Username,
	}, nil
}

func (kc *Client) GetUserRole(ctx context.Context, _ string, userID string) (string, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return "", err
	}

	roles, err := kc.client.GetCompositeRealmRolesByUserID(ctx, kc.token.AccessToken, kc.realm, userID)
	if err != nil {
		return "", err
	}
	for _, role := range roles {
		switch domain.UserRole(*role.Name) {
		case domain.UserCandidateRole, domain.UserCompanyRole, domain.UserAdminRole:
			return *role.Name, nil
		}
	}

	return "", domain.ErrInvalidRole
}

func parseError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case strings.Contains(err.Error(), "User exists with same email"):
		return domain.ErrEmailAlreadyExists
	case strings.Contains(err.Error(), "User exists with same username"):
		return domain.ErrUserNameAlreadyExists
	case strings.Contains(err.Error(), "Invalid user credentials"):
		return domain.ErrIncorrectPassword
	case strings.Contains(err.Error(), "401 Unauthorized"):
		return domain.ErrUnauthorized
	default:
		return err
	}
}
