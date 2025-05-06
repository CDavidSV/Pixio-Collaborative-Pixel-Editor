package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/CDavidSV/Pixio/validator"
)

func (h *Handler) CreateCanvasPost(w http.ResponseWriter, r *http.Request) {
	var createCanvasDTO types.CreateCanvasDTO
	if err := json.NewDecoder(r.Body).Decode(&createCanvasDTO); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid json body",
		})
		return
	}

	result, err := validator.Validate(createCanvasDTO)
	if err != nil {
		utils.ServerError(w, r, err, "Error validating request body")
		return
	}

	if !result.IsValid {
		result.SendValidationError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Map{
		"message": "Successfully created new canvas",
	})
}
