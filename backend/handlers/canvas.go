package handlers

import (
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) PostCreateCanvas(w http.ResponseWriter, r *http.Request) {
	createCanvasDTO, ok := utils.DecodeJSONAndValidate[types.CreateCanvasDTO](w, r)
	if !ok {
		return
	}

	userID := r.Context().Value(utils.UserIDKey).(string)
	pixelArr := make([]types.Pixel, createCanvasDTO.Width*createCanvasDTO.Height)
	pixelBytes, err := h.services.CanvasService.CompressPixelData(pixelArr)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create canvas")
		return
	}

	canvas, err := h.queries.CreateCanvas(createCanvasDTO.Title, createCanvasDTO.Description, userID, createCanvasDTO.Width, createCanvasDTO.Height, pixelBytes)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create canvas")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"canvas_id":   canvas.ID,
		"created_at":  canvas.CreatedAt,
		"access_type": canvas.LinkAccessType,
		"pixel_data":  canvas.PixelData,
	})
}

func (h *Handler) GetCanvas(w http.ResponseWriter, r *http.Request) {
	canvasID := chi.URLParam(r, "id")

	if canvasID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Canvas id must be provided",
		})
		return
	}

	// Get access rule for the user from context
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)

	canvas, err := h.queries.GetCanvas(canvasID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to fetch canvas")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"canvas":      canvas,
		"access_role": userAccess.AccessRole,
	})
}

func (h *Handler) PostDeleteCanvas(w http.ResponseWriter, r *http.Request) {
	deleteCanvasDTO, ok := utils.DecodeJSONAndValidate[types.DeleteCanvasDTO](w, r)
	if !ok {
		return
	}

	userID := r.Context().Value(utils.UserIDKey).(string)
	canvasOwnerID, err := h.queries.GetCanvasOwner(deleteCanvasDTO.CanvasID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to fetch canvas")
		return
	}

	if canvasOwnerID != userID {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "User is not the owner",
		})
		return
	}

	if err := h.queries.DeleteCanvas(deleteCanvasDTO.CanvasID); err != nil {
		utils.ServerError(w, r, err, "Failed to delete canvas")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "Canvas deleted successfully",
	})
}

func (h *Handler) PostJoinCanvasRoom(w http.ResponseWriter, r *http.Request) {
	canvasID := chi.URLParam(r, "id")
	userID := r.Context().Value(utils.UserIDKey).(string)

	if canvasID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Canvas id must be provided",
		})
		return
	}

	// Get access rule for the user from context
	userAccess := r.Context().Value(utils.AccessRuleKey).(types.UserAccess)

	canvas, err := h.queries.GetCanvas(canvasID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to fetch canvas")
		return
	}

	if err := h.websocket.JoinRoom(canvas.ID, userID, userAccess); err != nil {
		utils.ServerError(w, r, err, "Failed to join canvas room")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"canvas": canvas,
		"access": userAccess,
	})
}

func (h *Handler) PostLeaveCanvasRoom(w http.ResponseWriter, r *http.Request) {
	canvasID := chi.URLParam(r, "id")
	userID := r.Context().Value(utils.UserIDKey).(string)

	if canvasID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
			Error: "Canvas id must be provided",
		})
		return
	}

	h.websocket.LeaveRoom(canvasID, userID)
	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "Left canvas room successfully",
	})
}
