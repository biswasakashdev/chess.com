package handlers

import (
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/users"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userServ *users.UserService
}

func NewUserRouter(userService *users.UserService) chi.Router {

	userHandl := userHandler{
		userServ: userService,
	}
	router := chi.NewRouter()

	router.HandleFunc("/", userHandl.getUser)

	return router
}

type UserResp struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (uh *userHandler) getUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {

	}
	user, err := uh.userServ.GetUserById(userId, ctx)

	if err != nil {
		util.BuildErrResponse(w, err)
		return
	}

	userResp := UserResp{
		Id:        user.Id.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	util.BuildResponseWithBody(w, &userResp)
}

func (uh *userHandler) findUserByUsername(w http.ResponseWriter, r *http.Request) {

}
