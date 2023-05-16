package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	utils2 "github.com/gotalk/api/utils"
	"github.com/gotalk/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// @Summary      Ws connection with chat room
// @Description  join room
// @Tags         ws
// @Param token header string true "token getting after joining room"
// @Success      200  {string} string
// @Router       /ws/{id} [get]
func (h *Handler) wsConnection(w http.ResponseWriter, r *http.Request) {
	// get body and header params
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Error("invalid room id", zap.Error(err))
		handleError(http.StatusBadRequest, "invalid room id", err, w)
		return
	}

	// get jwt token from header
	token, err := getJWTToken(r)
	if err != nil {
		h.logger.Error("invalid jwt token", zap.Error(err))
		handleError(http.StatusUnauthorized, "invalid jwt token", err, w)
		return
	}

	// make authentication
	userId, err := authenticate(roomId, token)
	if err != nil {
		h.logger.Error("authenticate user", zap.Error(err))
		handleError(http.StatusInternalServerError, "authenticate user", err, w)
		return
	}

	user, err := h.service.GetUserById(userId)
	if err != nil {
		h.logger.Error("get user by id", zap.String("user id", userId), zap.Error(err))
		handleError(http.StatusInternalServerError, "get user by id", err, w)
		return
	}

	// make websocket connection
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket connection", zap.Error(err))
		handleError(http.StatusInternalServerError, "websocket connection", err, w)
		return
	}
	h.service.MakeWSConnection(conn, roomId, user.Email)
}

// @Summary      Join chat room
// @Description  join room
// @Tags         ws
// @Param token header string true "auth token"
// @Success      200  {string} string
// @Router       /ws/{id}/join [post]
func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) {
	var input *models.JoinRoomInput
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error("parse input body", zap.Error(err))
			handleError(http.StatusBadRequest, "parse input body", err, w)
			return
		}
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Error("invalid room id", zap.Error(err))
		handleError(http.StatusBadRequest, "invalid room id", err, w)
		return
	}

	// get user id from jwt token in request header
	userId, err := verifyUserId(r)
	if err != nil {
		h.logger.Error("invalid jwt token", zap.Error(err))
		handleError(http.StatusUnauthorized, "invalid jwt token", err, w)
		return
	}

	// authentication for chat room by room id and user id
	err = h.service.AuthenticateInRoom(input, roomId, userId)
	if err != nil {
		h.logger.Error("authenticate in room by room and user id", zap.Error(err))
		handleError(http.StatusUnauthorized, "authenticate in room by room and user id", err, w)
		return
	}

	// create token for ws connection
	tokenWS, err := createToken(roomId, userId)
	if err != nil {
		h.logger.Error("create jwt token for ws connection", zap.Error(err))
		handleError(http.StatusInternalServerError, "create jwt token for ws connection", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenWS))
}

func getJWTToken(r *http.Request) (string, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", errors.New("there is no provided token")
	}
	return token, nil
}

func createToken(roomId int, userId string) (string, error) {
	jwtFields := []utils2.JWTTokenField{
		{
			Type:  utils2.RoomId,
			Value: roomId,
		},
		{
			Type:  utils2.UserId,
			Value: userId,
		},
	}
	return utils2.CreateToken(jwtFields)
}

func authenticate(roomIdHeader int, token string) (string, error) {
	// get room id
	roomParam, err := utils2.VerifyToken(token, utils2.RoomId)
	if err != nil {
		return "", errors.Wrap(err, "invalid token")
	}
	roomId, ok := roomParam.(float64)
	if !ok {
		return "", errors.New("convert room id in jwt token from interface to int")
	}

	if roomIdHeader != int(roomId) {
		return "", errors.New("invalid room id")
	}

	// get user id
	userIdParam, err := utils2.VerifyToken(token, utils2.UserId)
	if err != nil {
		return "", errors.Wrap(err, "invalid token")
	}
	userId, ok := userIdParam.(string)
	if !ok {
		return "", errors.New("convert user userId in jwt token from interface to string")
	}

	return userId, nil
}
