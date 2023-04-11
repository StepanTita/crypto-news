package handlers

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"twitter-bot/internal/oauth"
	"twitter-bot/internal/web"
	"twitter-bot/internal/web/ctx"
)

func GetOAuth(w http.ResponseWriter, r *http.Request) {
	log := ctx.Log(r)

	log.Debug("received oauth login request")

	twitter := ctx.OAuth(r)

	oauthState, err := oauth.GenerateOAuthState()
	if err != nil {
		log.WithError(err).Error("failed to generate oauth state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authUrl := twitter.AuthCodeURL(oauthState,
		// TODO: tweak for production
		oauth2.SetAuthURLParam("code_challenge", "challenge"),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"))

	cookie := http.Cookie{Name: web.OAuthState, Value: oauthState, Expires: time.Now().Add(2 * time.Hour)}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}
