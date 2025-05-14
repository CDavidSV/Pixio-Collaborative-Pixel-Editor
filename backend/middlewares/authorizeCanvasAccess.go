package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/go-chi/chi/v5"
)

func (m *Middleware) AuthorizeCanvasAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		canvasID := chi.URLParam(r, "id")
		userID := r.Context().Value(utils.UserIDKey).(string)

		if canvasID == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
				Error: "canvas id is required",
			})
			return
		}

		userAccess, canvas, err := m.services.CanvasService.UserHasAccess(canvasID, userID)
		if err != nil {
			if errors.Is(err, types.ErrUserAccessDenied) {
				utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
					Error: "You do not have permission to access this canvas",
				})
				return
			}

			if errors.Is(err, types.ErrCanvasDoesNotExist) {
				utils.WriteJSON(w, http.StatusNotFound, types.ErrorResponse{
					Error: "Canvas does not exist",
				})
				return
			}

			utils.ServerError(w, r, err, "Unable to fetch user access")
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.CanvasKey, canvas)
		ctx = context.WithValue(ctx, utils.AccessRuleKey, userAccess.AccessRole)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
