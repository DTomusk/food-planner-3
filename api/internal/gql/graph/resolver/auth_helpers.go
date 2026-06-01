package resolver

import (
	"net/http"
	"time"
)

func setRefreshTokenCookie(w http.ResponseWriter, token string, expiresAt int64) {
	if w == nil {
		return
	}

	maxAge := int(expiresAt - time.Now().Unix())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // make sure to set this to true in production
		SameSite: http.SameSiteNoneMode,
	})
}

func clearRefreshTokenCookie(w http.ResponseWriter) {
	if w == nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // make sure to set this to true in production
		SameSite: http.SameSiteNoneMode,
	})
}
