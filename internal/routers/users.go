package routers

import (
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/hub"
	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/service"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userServ *service.UserService
	gameHub  *hub.Hub
}

func NewUserRouter(userService *service.UserService) chi.Router {

	userHandl := userHandler{
		userServ: userService,
	}
	router := chi.NewRouter()

	router.Get("/", userHandl.getUsers)
	router.Get("/profile", userHandl.getUserProfile)

	return router
}

func (uh *userHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	userResp, err := uh.userServ.GetUserById(userId, ctx)

	if err != nil {
		util.BuildErrResponse(w, err)
		return
	}

	util.BuildResponseWithBody(w, &userResp)
}

func (uh *userHandler) getUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	username := queryParams.Get("search")

	users, err := uh.userServ.GetUserByUsernameNotFriendWith(ctx, userId, username)
	if err != nil {
		util.BuildErrResponse(w, errors.New("Invalid Token"))
		return
	}

	util.BuildResponseWithBody(w, &users)

}
