package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/validator"
	"github.com/oklog/ulid/v2"
)

type contextKey string

const (
	UserIDKey  contextKey = "userID"
	CanvasKey  contextKey = "canvas"
	AccessRule contextKey = "accessRule"
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

func DecodeJSONAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid json body",
		})
		return nil, false
	}

	result, err := validator.Validate(body)
	if err != nil {
		ServerError(w, r, err, "Error validating request body")
		return nil, false
	}

	if !result.IsValid {
		result.SendValidationError(w)
		return nil, false
	}

	return &body, true
}
