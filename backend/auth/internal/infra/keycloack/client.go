package keycloack

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

var ErrorInvalidToken = errors.New("invalid access token")

type Client struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	realm  string

	clientID     string
	clientSecret string
	adminUser    string
	adminPass    string
}

func NewClient(clientUrl, realm, clientID, clientSecret, adminUser, adminPass string) *Client {
	client := gocloak.NewClient(clientUrl)
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
			return err
		}
	}
	istResult, err := kc.client.RetrospectToken(ctx, kc.token.AccessToken, kc.clientID, kc.clientSecret, kc.realm)
	if err != nil {
		if err = kc.loginAdmin(ctx); err != nil {
			return err
		}
	}

	if !*istResult.Active {
		if err = kc.loginAdmin(ctx); err != nil {
			return err
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

	userId, err := kc.client.CreateUser(ctx, kc.token.AccessToken, kc.realm, newUser)
	if err != nil {
		return "", err
	}

	if err = kc.client.SetPassword(ctx, kc.token.AccessToken, userId, kc.realm, password, false); err != nil {
		_ = kc.client.DeleteUser(ctx, kc.token.AccessToken, kc.realm, userId)
		return "", err
	}

	if err = kc.addRole(ctx, userId, user.Role); err != nil {
		_ = kc.client.DeleteUser(ctx, kc.token.AccessToken, kc.realm, userId)
		return "", err
	}

	return userId, err
}

func (kc *Client) addRole(ctx context.Context, userId, roleName string) error {
	roles := make([]gocloak.Role, 1)

	if roleName == "Candidate" {
		roles[0] = gocloak.Role{
			ID:   gocloak.StringP("15bd1c8f-1feb-4870-9f46-a847f0742be9"),
			Name: gocloak.StringP("Candidate"),
		}
	} else if roleName == "Company" {
		roles[0] = gocloak.Role{
			ID:   gocloak.StringP("2e90e50e-8db4-4881-8185-05a40220f759"),
			Name: gocloak.StringP("Company"),
		}
	} else {
		return fmt.Errorf("the roleName is not valid: %s", roleName)
	}

	return kc.client.AddRealmRoleToUser(ctx, kc.token.AccessToken, kc.realm, userId, roles)

}

func (kc *Client) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	token, err := kc.client.Login(ctx, kc.clientID, kc.clientSecret, kc.realm, username, password)

	return token, err
}

func (kc *Client) Logout(ctx context.Context, token string) error {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return err
	}

	err := kc.client.Logout(ctx, kc.clientID, kc.clientSecret, kc.realm, token)
	return err
}

func (kc *Client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	if err := kc.validateToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	token, err := kc.client.RefreshToken(ctx, refreshToken, kc.clientID, kc.clientSecret, kc.realm)
	return token, err
}

func (kc *Client) validateToken(ctx context.Context, token string) error {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return err
	}

	istResult, err := kc.client.RetrospectToken(ctx, token, kc.clientID, kc.clientSecret, kc.realm)
	if err != nil {
		return ErrorInvalidToken
	}
	if !(*istResult.Active) {
		return ErrorInvalidToken
	}

	return nil
}

func (kc *Client) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return nil, err
	}

	user_info, err := kc.client.GetUserInfo(ctx, token, kc.realm)
	if err != nil {
		return nil, err
	}

	user, err := kc.client.GetUserByID(ctx, kc.token.AccessToken, kc.realm, *user_info.Sub)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        *user_info.Sub,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Username:  *user.Username,
	}, nil
}

func (kc *Client) GetUserRole(ctx context.Context, token string, userId string) (string, error) {
	if err := kc.ensureAdminTokenValid(ctx); err != nil {
		return "", err
	}

	roles, err := kc.client.GetCompositeRealmRolesByUserID(ctx, kc.token.AccessToken, kc.realm, userId)
	if err != nil {
		return "", err
	}
	for _, role := range roles {
		if *role.Name == "Candidate" {
			return "Candidate", nil
		} else if *role.Name == "Company" {
			return "Company", nil
		}
	}

	return "", fmt.Errorf("user role not found")
}
