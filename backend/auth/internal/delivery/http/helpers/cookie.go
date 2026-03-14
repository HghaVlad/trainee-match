package helpers

import (
	"net/http"
	"time"
)

func SetTokenPairToCookies(w http.ResponseWriter, accessToken, refreshToken string, accessTokenExpires, refreshTokenExpires int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Duration(accessTokenExpires) * time.Second),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(refreshTokenExpires) * time.Second),
		HttpOnly: true,
		Path:     "/",
	})
}

func GetRefreshTokenFromCookies(r *http.Request) string {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			return cookie.Value
		}
	}
	return ""
}

func GetAccessTokenFromCookies(r *http.Request) string {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value
		}
	}
	return ""
}
