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
