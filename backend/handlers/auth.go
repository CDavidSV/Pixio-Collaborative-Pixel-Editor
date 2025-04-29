package handlers

import (
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/CDavidSV/Pixio/validator"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, r, err, "Failed to parse form")
		return
	}

	userSignupDTO := types.UserSignupDTO{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	result, err := validator.Validate(userSignupDTO)
	if err != nil {
		utils.ServerError(w, r, err, "Error validating request body")
		return
	}

	if !result.IsValid {
		result.SendValidationError(w)
		return
	}

	// Attempt to create the user
	user, err := h.queries.User.CreateUser(userSignupDTO.Username, userSignupDTO.Email, userSignupDTO.Password)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create user")
		return
	}

	// Start a session for the user
	session, err := h.services.AuthService.CreateSession(user.ID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create session")
		return
	}

	// Set refresh token cookie
	utils.SetCookie(w, "refresh_token", session.RefreshToken, session.ExpiresAt) // 30 days
	utils.WriteJSON(w, http.StatusCreated, types.Map{
		"token":      session.AccessToken,
		"expires_at": session.AccessTokenExpiresAt,
		"user":       user,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, r, err, "Failed to parse form")
		return
	}

	userLoginDTO := types.UserLoginDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	result, err := validator.Validate(userLoginDTO)
	if err != nil {
		utils.ServerError(w, r, err, "Error validating request body")
		return
	}

	if !result.IsValid {
		result.SendValidationError(w)
		return
	}

	// Attempt to authenticate the user
	err = h.services.AuthService.Authenticate(userLoginDTO.Email, userLoginDTO.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	// Start a session for the user
	session, err := h.services.AuthService.CreateSession(userLoginDTO.Email)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create session")
		return
	}

	// Set refresh token cookie
	utils.SetCookie(w, "refresh_token", session.RefreshToken, session.ExpiresAt) // 30 days
	utils.WriteJSON(w, http.StatusOK, types.Map{
		"token":      session.AccessToken,
		"expires_at": session.AccessTokenExpiresAt,
		"user_id":    session.UserID,
	})
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

}
