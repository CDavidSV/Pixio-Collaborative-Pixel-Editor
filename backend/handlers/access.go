package handlers

import (
	"net/http"
)

func (h *Handler) PostCreateAccess(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value(utils.UserIDKey).(string)

	// createAccessDTO, ok := utils.DecodeJSONAndValidate[types.CreateAccessDTO](w, r)
	// if !ok {
	// 	return
	// }

	// createdUserAccess, err := h.queries.CreateUserAccess(createAccessDTO.CanvasID, types.CanvasObject, createAccessDTO.AccessRole)
}

func (h *Handler) PostDeleteAccess(w http.ResponseWriter, r *http.Request) {
	// deleteAccessDTO, ok := utils.DecodeJSONAndValidate[types.DeleteAccessDTO](w, r)
	// if !ok {
	// 	return
	// }
}

func (h *Handler) PutUpdateAccess(w http.ResponseWriter, r *http.Request) {
	// updateAccessDTO, ok := utils.DecodeJSONAndValidate[types.UpdateAccessDTO](w, r)
	// if !ok {
	// 	return
	// }
}

func (h *Handler) GetAccessRules(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) PutUpdateGlobalAccess(w http.ResponseWriter, r *http.Request) {
	// updateGlobalAccessDTO, ok := utils.DecodeJSONAndValidate[types.UpdateGlobalAccessDTO](w, r)
	// if !ok {
	// 	return
	// }
}
