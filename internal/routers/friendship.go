package routers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/service"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type friendshipHandler struct {
	frndshpService *service.FriendShipService
}

func NewFriendshipRouter(frndshpServ *service.FriendShipService) chi.Router {

	router := chi.NewRouter()

	frndshpHandl := friendshipHandler{
		frndshpService: frndshpServ,
	}

	router.Get("/", frndshpHandl.getFriends)
	router.Delete("/", frndshpHandl.RemoveFriend)
	router.Post("/", frndshpHandl.SendFriendRequest)
	router.Patch("/", frndshpHandl.UpdateFriendshipStatus)
	router.Get("/requests", frndshpHandl.getRequests)
	return router
}

/*
 * GET /
 * Return all friends
 * query(type, ['online','offline','all'])
 */
func (fh *friendshipHandler) getFriends(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	friendType := queryParams.Get("type")

	userPayloads, err := fh.frndshpService.FetchAllFriends(ctx, userId, friendType)

	if err != nil {
		util.BuildErrResponse(w, err)
	}
	util.BuildResponseWithBody(w, userPayloads)
}

type DeleteFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
}

// RemoveFriend handles declining a request, canceling a sent request, or unfriending.
// DELETE /
func (fh *friendshipHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
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

	if err := fh.frndshpService.DeleteRequest(r.Context(), userId, req.TargetUserID); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseBodyWithCode(w, map[string]string{
		"message": "Friendship request deleted successfully",
	}, http.StatusCreated)
}

type SendFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
}

// SendFriendRequest handles sending a new friend request.
// POST /
func (fh *friendshipHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
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

	if err := fh.frndshpService.SendFriendRequest(r.Context(), userId, req.TargetUserID); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseBodyWithCode(w, map[string]string{
		"message": "Friend request sent successfully",
	}, http.StatusCreated)
}

/*
 * GET /requests
 * Find all the requests or relationship if not accepted
 * query(type,['sent','blocked','pending'])
 */
func (fh *friendshipHandler) getRequests(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userId == "" {
		util.BuildErrResponseWithCode(w, errors.New("Invalid authentication"), 400)
	}

	queryParams := r.URL.Query()

	// Type can be ['sent', '']
	friendType := queryParams.Get("type")

	userPayloads, err := fh.frndshpService.FetchAllRequests(ctx, userId, friendType)

	if err != nil {
		util.BuildErrResponse(w, err)
	}
	util.BuildResponseWithBody(w, userPayloads)
}

type UpdateFriendshipRequest struct {
	TargetUserID string `json:"target_user_id"`
	Action       string `json:"action"` // "accept" or "block" or "unblock" or "cancel"
}

// UpdateFriendshipStatus handles accepting or blocking a friend request.
// PATCH /requests
func (fh *friendshipHandler) UpdateFriendshipStatus(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
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

	if err := fh.frndshpService.UpdateFriendShipStatus(r.Context(), userId, req.TargetUserID, req.Action); err != nil {
		util.BuildErrResponseWithCode(w, err, http.StatusBadRequest)
		return
	}

	util.BuildResponseBodyWithCode(w, map[string]string{
		"message": "Friendship status updated successfully",
	}, http.StatusAccepted)
}
