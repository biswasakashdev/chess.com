package routers

import (
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/hub"
	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

func NewGameRouter(h *hub.Hub) chi.Router {
	router := chi.NewRouter()

	router.Get("/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, ok := ctx.Value(middleware.UserIDKey).(string)

		if !ok || userId == "" {
			util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
		}
		// Extract the path variable using chi.URLParam
		roomId := chi.URLParam(r, "roomId")

		game, err := h.GetGameDetails(roomId, userId)

		if err != nil {
			util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		}

		util.BuildResponseBodyWithCode(w, game, http.StatusOK)

	})

	router.Post("/{roomId}", func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		userId, ok := ctx.Value(middleware.UserIDKey).(string)

		if !ok || userId == "" {
			util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
		}
		// Extract the path variable using chi.URLParam
		roomId := chi.URLParam(r, "roomId")

		err := h.DeleteRoom(userId, roomId)

		if err != nil {
			util.BuildErrResponse(
				w,
				err,
			)
		}

		util.BuildResponseBodyWithCode(w, map[string]string{
			"message": "Room deleted successfully",
		}, http.StatusAccepted)

	})

	return router
}
