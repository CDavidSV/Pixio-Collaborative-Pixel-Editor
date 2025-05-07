package handlers

import (
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
)

func (h *Handler) CreateCanvasPost(w http.ResponseWriter, r *http.Request) {
	createCanvasDTO, ok := utils.DecodeJSONAndValidate[types.CreateCanvasDTO](w, r)
	if !ok {
		return
	}

	userID := r.Context().Value(utils.UserIDKey).(string)
	canvas, err := h.services.CanvasService.CreateCanvas(createCanvasDTO.Title, createCanvasDTO.Descrition, createCanvasDTO.Width, createCanvasDTO.Height, userID)
	if err != nil {
		utils.ServerError(w, r, err, "Failed to create canvas")
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"canvas_id":   canvas.ID,
		"created_at":  canvas.CreatedAt,
		"access_type": canvas.AccessType,
		"pixel_data":  canvas.PixelData,
	})
}

func (h *Handler) DeleteCanvasPost(w http.ResponseWriter, r *http.Request) {

}
