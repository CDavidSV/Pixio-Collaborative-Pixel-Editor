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
type ClientErrorCode int

const (
	UserIDKey     contextKey = "userID"
	AccessRuleKey contextKey = "accessRule"
)

const (
	// 400 Bad Request
	ErrInvalidJSONBody  ClientErrorCode = 1000
	ErrCanvasIDRequired ClientErrorCode = 1001
	ErrInvalidID        ClientErrorCode = 1002

	// 401 Unauthorized
	ErrInvalidCredentials         ClientErrorCode = 1100
	ErrNotLoggedIn                ClientErrorCode = 1101
	ErrSessionExpired             ClientErrorCode = 1102
	ErrSessionNotFound            ClientErrorCode = 1103
	ErrInvalidRefreshToken        ClientErrorCode = 1104
	ErrLogoutFailed               ClientErrorCode = 1105
	ErrForbiddenCanvasAccess      ClientErrorCode = 1106
	ErrNotCanvasOwner             ClientErrorCode = 1107
	ErrAccessRulesUpdateForbidden ClientErrorCode = 1108

	// 404 Not Found
	ErrUserNotFound   ClientErrorCode = 1200
	ErrCanvasNotFound ClientErrorCode = 1201

	// 409 Conflict
	ErrUserAlreadyRegistered ClientErrorCode = 1300
)

var clientErrorCodes = map[ClientErrorCode]string{
	// 400 Bad Request
	ErrInvalidJSONBody:  "Invalid JSON body",
	ErrCanvasIDRequired: "Canvas ID must be provided",
	ErrInvalidID:        "ID must be provided and of valid format",

	// 401 Unauthorized
	ErrInvalidCredentials:         "Invalid email or password",
	ErrNotLoggedIn:                "User not logged in",
	ErrSessionExpired:             "Session expired",
	ErrSessionNotFound:            "Session not found",
	ErrInvalidRefreshToken:        "Invalid refresh token",
	ErrLogoutFailed:               "User cannot be logged out",
	ErrForbiddenCanvasAccess:      "You do not have permission to access this canvas",
	ErrNotCanvasOwner:             "User is not the owner",
	ErrAccessRulesUpdateForbidden: "User not allowd to update access rules",

	// 404 Not Found
	ErrUserNotFound:   "User not found",
	ErrCanvasNotFound: "Canvas does not exist",

	// 409 Conflict
	ErrUserAlreadyRegistered: "User already registered",
}

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

func ClientError(w http.ResponseWriter, status int, errCode ClientErrorCode) {
	errMsg := clientErrorCodes[errCode]

	WriteJSON(w, status, types.Map{
		"code":  errCode,
		"error": errMsg,
	})
}
