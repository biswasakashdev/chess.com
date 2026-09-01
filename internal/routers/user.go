package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/hub"
	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/users"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userServ *users.UserService
	gameHub  *hub.Hub
}

func NewUserRouter(userService *users.UserService) chi.Router {

	userHandl := userHandler{
		userServ: userService,
	}
	router := chi.NewRouter()

	router.Get("/", userHandl.getUsersByUsername)
	router.Get("/profile", userHandl.getUserProfile)
	router.Get("/friends", userHandl.getFriends)
	router.Post("/requests", userHandl.getRequests)

	return router
}

type UserResp struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (uh *userHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
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

func (uh *userHandler) getUsersByUsername(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	username := queryParams.Get("search")
	if username == "" {
		util.BuildErrResponse(w, errors.New("Invalid Token"))
		return
	}

	users, err := uh.userServ.GetUserByUsernameNotFriendWith(ctx, userId, username)
	if err != nil {
		util.BuildErrResponse(w, errors.New("Invalid Token"))
		return
	}

	util.BuildResponseWithBody(w, &users)

}

func (uh *userHandler) getFriends(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	// Type can be ['requests', 'active', '']
	friendType := queryParams.Get("type")

	userPayloads, err := uh.userServ.FetchAllFriends(ctx, userId, friendType)

	if err != nil {
		util.BuildErrResponse(w, err)
	}
	util.BuildResponseWithBody(w, userPayloads)
}

func (uh *userHandler) getRequests(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	// Type can be ['sent', '']
	friendType := queryParams.Get("type")

	userPayloads, err := uh.userServ.FetchAllRequests(ctx, userId, friendType)

	if err != nil {
		util.BuildErrResponse(w, err)
	}
	util.BuildResponseWithBody(w, userPayloads)
}

type UpdateFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
	Action       string `json:"action"` // "accept" or "block"
}

// UpdateFriendshipStatus handles accepting or blocking a friend request.
// POST /friends/status
func (uh *userHandler) UpdateFriendshipStatus(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), http.StatusUnauthorized)
		return
	}

	var req UpdateFriendshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BuildErrResponseWithCode(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.TargetUserID == "" || req.Action == "" {
		util.BuildErrResponseWithCode(w, errors.New("target_user_id and action are required"), http.StatusBadRequest)
		return
	}

	if err := uh.userServ.UpdateFriendShipStatus(r.Context(), userId, req.TargetUserID, req.Action); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseWithBodyAndCode(w, map[string]string{
		"message": "Friendship status updated successfully",
	}, http.StatusOK)
}

type DeleteFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
}

// DeleteFriendRequest handles declining a request, canceling a sent request, or unfriending.
// DELETE /friends/request (or DELETE /friends)
func (uh *userHandler) DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), http.StatusUnauthorized)
		return
	}

	var req DeleteFriendshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BuildErrResponseWithCode(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.TargetUserID == "" {
		util.BuildErrResponseWithCode(w, errors.New("target_user_id is required"), http.StatusBadRequest)
		return
	}

	if err := uh.userServ.DeleteRequest(r.Context(), userId, req.TargetUserID); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseWithBodyAndCode(w, map[string]string{
		"message": "Friendship request deleted successfully",
	}, http.StatusOK)
}

type SendFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
}

// SendFriendRequest handles sending a new friend request.
// POST /friends/request
func (uh *userHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), http.StatusUnauthorized)
		return
	}

	var req SendFriendshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BuildErrResponseWithCode(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.TargetUserID == "" {
		util.BuildErrResponseWithCode(w, errors.New("target_user_id is required"), http.StatusBadRequest)
		return
	}

	if err := uh.userServ.SendFriendRequest(r.Context(), userId, req.TargetUserID); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseWithBodyAndCode(w, map[string]string{
		"message": "Friend request sent successfully",
	}, http.StatusCreated)
}
