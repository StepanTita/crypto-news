package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"

	"common/data/model"
	"twitter-bot/internal/oauth"
	"twitter-bot/internal/web"
	"twitter-bot/internal/web/ctx"
)

func GetOAuthCallback(w http.ResponseWriter, r *http.Request) {
	log := ctx.Log(r)
	reqCtx := r.Context()

	log.Debug("received oauth callback")

	oauthState, err := r.Cookie(web.OAuthState)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !oauth.VerifyOAuthState(r.URL.Query().Get("state"), oauthState.Value) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	twitter := ctx.OAuth(r)
	token, err := twitter.Exchange(reqCtx,
		code,
		// TODO: update for prod
		oauth2.SetAuthURLParam("code_verifier", "challenge"))
	if err != nil {
		log.WithError(err).Error("failed to get token for twitter API")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataProvider := ctx.DataProvider(r)
	// For now twitter bot can only be one. Easy to fix that if make a unique key here
	id, err := dataProvider.KVProvider().SetStruct(reqCtx, model.ToKey(token, false), token, 0)
	if err != nil {
		log.WithError(err).Error("failed to set access token to the redis storage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithField("key", id).Debug("Success setting the access key to the redis storage")
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "Success setting the authorization key",
		Status:  http.StatusOK,
	})
	if err != nil {
		log.WithError(err).Error("failed to marshal oauth callback response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		log.WithError(err).Error("failed to write oauth callback response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
