package handlers

import (
	"errors"
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) PostCreateAccess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(string)
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)
	canvasID := chi.URLParam(r, "id")

	if len(canvasID) != 26 {
		utils.ClientError(w, http.StatusBadRequest, utils.ErrInvalidID)
		return
	}

	if userAccess.AccessRole == types.Viewer {
		utils.ClientError(w, http.StatusUnauthorized, utils.ErrAccessRulesUpdateForbidden)
		return
	}

	createAccessDTO, ok := utils.DecodeJSONAndValidate[types.CreateAccessDTO](w, r)
	if !ok {
		return
	}

	userToGrantAccess, err := h.queries.GetUserByEmail(canvasID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.ClientError(w, http.StatusNotFound, utils.ErrUserNotFound)
			return
		}

		utils.ServerError(w, r, err, "Failed to fetch user")
		return
	}

	_, err = h.queries.CreateUserAccess(canvasID, types.CanvasObject, createAccessDTO.AccessRole, userToGrantAccess.ID, userID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed grant access to user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "Access granted to user",
		"user_id": userToGrantAccess.ID,
	})
}

func (h *Handler) PostDeleteAccess(w http.ResponseWriter, r *http.Request) {
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)
	canvasID := chi.URLParam(r, "id")

	if len(canvasID) != 26 {
		utils.ClientError(w, http.StatusBadRequest, utils.ErrInvalidID)
		return
	}

	if userAccess.AccessRole == types.Viewer {
		utils.ClientError(w, http.StatusUnauthorized, utils.ErrAccessRulesUpdateForbidden)
		return
	}

	deleteAccessDTO, ok := utils.DecodeJSONAndValidate[types.DeleteAccessDTO](w, r)
	if !ok {
		return
	}

	err := h.queries.DeleteUserAccess(canvasID, types.CanvasObject, deleteAccessDTO.UserID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to delete canvas access rule")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "Access removed from user",
		"user_id": deleteAccessDTO.UserID,
	})
}

func (h *Handler) PutUpdateAccess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(string)
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)
	canvasID := chi.URLParam(r, "id")

	if len(canvasID) != 26 {
		utils.ClientError(w, http.StatusBadRequest, utils.ErrInvalidID)
		return
	}

	if userAccess.AccessRole == types.Viewer {
		utils.ClientError(w, http.StatusUnauthorized, utils.ErrAccessRulesUpdateForbidden)
		return
	}

	updateAccessDTO, ok := utils.DecodeJSONAndValidate[types.UpdateAccessDTO](w, r)
	if !ok {
		return
	}

	err := h.queries.UpdateUserAccess(canvasID, types.CanvasObject, updateAccessDTO.AccessRole, userID, updateAccessDTO.UserID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to update user access rules")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message":     "Global canvas access rules updated",
		"access_role": updateAccessDTO.AccessRole,
	})
}

func (h *Handler) GetAccessRules(w http.ResponseWriter, r *http.Request) {
	canvasID := chi.URLParam(r, "id")

	if len(canvasID) != 26 {
		utils.ClientError(w, http.StatusBadRequest, utils.ErrInvalidID)
		return
	}

	userAccessRules, err := h.queries.GetAccessRules(canvasID, types.CanvasObject)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to fetch user access rules")
		return
	}

	utils.WriteJSON(w, http.StatusOK, userAccessRules)
}

func (h *Handler) PutUpdateGlobalAccess(w http.ResponseWriter, r *http.Request) {
	canvasID := chi.URLParam(r, "id")
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)

	if len(canvasID) != 26 {
		utils.ClientError(w, http.StatusBadRequest, utils.ErrInvalidID)
		return
	}

	if userAccess.AccessRole == types.Viewer {
		utils.ClientError(w, http.StatusUnauthorized, utils.ErrAccessRulesUpdateForbidden)
		return
	}

	updateGlobalAccessDTO, ok := utils.DecodeJSONAndValidate[types.UpdateGlobalAccessDTO](w, r)
	if !ok {
		return
	}

	if err := h.queries.UpdateLinkAccess(canvasID, updateGlobalAccessDTO.LinkAccessType, updateGlobalAccessDTO.LinkAccessRole); err != nil {
		utils.ServerError(w, r, err, "Failed to update global canvas access rules")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message":     "Global canvas access rules updated",
		"access_type": updateGlobalAccessDTO.LinkAccessType,
		"access_role": updateGlobalAccessDTO.LinkAccessRole,
	})
}
