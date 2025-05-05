package handlers

import (
	"errors"
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

	hashedPassword, err := h.services.AuthService.HashPassword(userSignupDTO.Password)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to hash password")
		return
	}

	// Attempt to create the user
	user, err := h.queries.User.CreateUser(userSignupDTO.Username, userSignupDTO.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			utils.WriteJSON(w, http.StatusConflict, types.ErrorResponse{
				Error: "User already registered",
			})
			return
		}

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
	utils.SetCookie(w, "rt", session.RefreshToken, int(session.ExpiresAt.Sub(session.CreatedAt).Seconds())) // 30 days
	utils.WriteJSON(w, http.StatusCreated, types.Map{
		"token":      session.AccessToken,
		"expires_at": session.AccessTokenExpiresAt.Unix(),
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

	// Check if the user exists
	user, err := h.queries.User.GetUserByEmail(userLoginDTO.Email)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	// Attempt to authenticate the user
	if !h.services.AuthService.ValidPassword(userLoginDTO.Password, user.HashedPassword) {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	// Start a session for the user
	session, err := h.services.AuthService.CreateSession(user.ID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create session")
		return
	}

	// Set refresh token cookie
	utils.SetCookie(w, "rt", session.RefreshToken, int(session.ExpiresAt.Sub(session.CreatedAt).Seconds())) // 30 days
	utils.WriteJSON(w, http.StatusOK, types.Map{
		"token":      session.AccessToken,
		"expires_at": session.AccessTokenExpiresAt.Unix(),
		"user":       user,
	})
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("rt")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "User not logged in",
			})
			return
		}

		utils.ServerError(w, r, err, "Failed to get cookie")
		return
	}

	session, err := h.services.AuthService.RevalidateSession(refreshToken.Value)
	if err != nil {
		if errors.Is(err, types.ErrSessionExpired) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "Session expired",
			})
			return
		} else if errors.Is(err, types.ErrSessionNotFound) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "Session not found",
			})
			return
		} else if errors.Is(err, types.ErrInvalidToken) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "Invalid refresh token",
			})
			return
		}

		utils.ServerError(w, r, err, "Failed to create session")
		return
	}

	utils.SetCookie(w, "rt", session.RefreshToken, int(session.ExpiresAt.Sub(session.CreatedAt).Seconds()))
	utils.WriteJSON(w, http.StatusOK, types.Map{
		"token":      session.AccessToken,
		"expires_at": session.AccessTokenExpiresAt.Unix(),
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("rt")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "User cannot be logged out",
			})
			return
		}

		utils.ServerError(w, r, err, "Failed to get cookie")
		return
	}

	err = h.services.AuthService.CloseSession(refreshToken.Value)
	if err != nil {
		if errors.Is(err, types.ErrInvalidToken) {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "User cannot be logged out",
			})
			return
		}

		utils.ServerError(w, r, err, "Failed to logout user")
		return
	}

	utils.SetCookie(w, "rt", "", -1) // Delete the cookie
	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "User logged out successfully",
	})
}
