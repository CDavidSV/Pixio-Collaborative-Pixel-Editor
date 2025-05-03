package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/oklog/ulid/v2"
)

func ServerError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	slog.Error("Internal Server Error", "error", err.Error(), "method", method, "uri", uri)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "` + msg + `"}`))
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to write JSON response", "error", err.Error())
	}
}

func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		MaxAge:   maxAge,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func GenerateID() string {
	id := ulid.Make()
	return id.String()
}
