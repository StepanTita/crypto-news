package handlers

import (
	"encoding/json"
	"net/http"

	"twitter-bot/internal/web/ctx"
)

func GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	log := ctx.Log(r)

	log.WithField("healthcheck", "[SUCCESS]").Debug("Healthcheck...")

	w.WriteHeader(http.StatusOK)

	body, err := json.Marshal(struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "Server runs successfully",
		Status:  http.StatusOK,
	})
	if err != nil {
		log.WithError(err).Error("failed to marshal healthcheck response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		log.WithError(err).Error("failed to write healthcheck response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
