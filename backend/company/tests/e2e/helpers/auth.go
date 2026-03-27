package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/google/uuid"
)

func GetAuthClient(authServiceBaseUrl string) *http.Client {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	username := "usertest" + strings.ReplaceAll(uuid.NewString(), "-", "")
	email := username + "@gmail.com"

	registerBody := fmt.Sprintf(`{
		"first_name":"Test",
		"last_name":"User",
		"email":"%s",
		"username":"%s",
		"password":"testpass",
		"role":"%s"
	}`, email, username, identity.RoleHR)

	resp, err := client.Post(
		authServiceBaseUrl+"/auth/register",
		"application/json",
		strings.NewReader(registerBody),
	)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		log.Fatalf("register error: %d %s %s", resp.StatusCode, resp.Status, errResp.Error)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	loginBody := fmt.Sprintf(`{
		"username":"%s",
		"password":"testpass"
	}`, username)

	resp, err = client.Post(
		authServiceBaseUrl+"/auth/login",
		"application/json",
		strings.NewReader(loginBody),
	)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	for _, c := range resp.Cookies() {
		log.Println("cookie:", c.Name, c.Value)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return client
}
