package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

		canvasAccessType, canvasAccessRole, err := m.queries.GetCanvasLinkAccess(canvasID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				utils.WriteJSON(w, http.StatusNotFound, types.ErrorResponse{
					Error: "This canvas does not exist",
				})
				return
			}

			utils.ServerError(w, r, err, "Unable to fetch canvas access")
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		if canvasAccessType == types.WithLink {
			ctx = context.WithValue(ctx, utils.AccessRuleKey, canvasAccessRole)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		userAccess, err := m.queries.GetUserAccess(canvasID, types.CanvasObject, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
					Error: "You do not have permission to access this canvas",
				})
				return
			}

			utils.ServerError(w, r, err, "Unable to fetch user access")
			return
		}

		ctx = context.WithValue(ctx, utils.AccessRuleKey, userAccess.AccessRole)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
