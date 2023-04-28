package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

// @Summary      Create Room
// @Description  create room
// @Tags         room
// @Param 		 input body models.Room true "room"
// @Param token header string true "auth token"
// @Success      200  {int}  int
// @Router       /room [post]
func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var room *models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		h.logger.Error("parse input body json", zap.Error(err))
		handleError(http.StatusBadRequest, "parse input body json", err, w)
		return
	}
	if room.Private && room.Password == "" {
		h.logger.Error("password not provided", zap.Error(err))
		handleError(http.StatusBadRequest, "password not provided", err, w)
		return
	}

	// get user userId from jwt token in header
	userId, err := verifyUserId(r)
	if err != nil {
		h.logger.Error("verify token", zap.Error(err))
		handleError(http.StatusUnauthorized, "verify token", err, w)
		return
	}
	room.CreatorId = userId

	roomId, err := h.service.CreateRoom(room)
	if err != nil {
		h.logger.Error("create room", zap.Error(err))
		handleError(http.StatusInternalServerError, "create room", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(roomId)))
}

// @Summary      Update room
// @Description  update room
// @Tags         room
// @Param 		 input body models.UpdateRoomInput true "update room"
// @Param token header string true "auth token"
// @Success      200  {string}  string
// @Router       /room/{id} [put]
func (h *Handler) updateRoomById(w http.ResponseWriter, r *http.Request) {
	var room *models.UpdateRoomInput
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		h.logger.Error("parse input body json", zap.Error(err))
		handleError(http.StatusBadRequest, "parse input body json", err, w)
		return
	}
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Error("invalid room id", zap.Error(err))
		handleError(http.StatusBadRequest, "invalid room id", err, w)
		return
	}

	// get user id from jwt token in header
	userId, err := verifyUserId(r)

	err = h.service.UpdateRoom(room, roomId, userId)
	if err != nil {
		log.Println(err)
		h.logger.Error("update room", zap.Error(err))
		handleError(http.StatusInternalServerError, "update room", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully updated the room"))
}

// @Summary      Get all rooms
// @Description  get all rooms
// @Tags         room
// @Param token header string true "auth token"
// @Success      200  {object}  []models.RoomResponse
// @Router       /room [get]
func (h *Handler) getAllRooms(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserId(r)
	if err != nil {
		h.logger.Error("invalid token", zap.Error(err))
		handleError(http.StatusUnauthorized, "invalid token", err, w)
		return
	}

	rooms, err := h.service.GetAllRooms()
	if err != nil {
		log.Println(err)
		h.logger.Error("get all rooms", zap.Error(err))
		handleError(http.StatusInternalServerError, "get all rooms", err, w)
		return
	}

	resultBody, err := json.MarshalIndent(rooms, "", " ")
	if err != nil {
		log.Println(err)
		h.logger.Error("parse rooms to json", zap.Error(err))
		handleError(http.StatusInternalServerError, "parse room to json", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultBody)
}

// @Summary      Get room by id
// @Description  get room by id
// @Tags         room
// @Param token header string true "auth token"
// @Success      200  {object}  models.RoomResponse
// @Router       /room/{id} [get]
func (h *Handler) getRoomById(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserId(r)
	if err != nil {
		h.logger.Error("invalid token", zap.Error(err))
		handleError(http.StatusUnauthorized, "invalid token", err, w)
		return
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Error("invalid room id", zap.Error(err))
		handleError(http.StatusBadRequest, "invalid room id", err, w)
		return
	}

	room, err := h.service.GetRoomById(roomId)
	if err != nil {
		h.logger.Error("get room by id", zap.Error(err))
		handleError(http.StatusInternalServerError, "get room by id", err, w)
		return
	}

	roomBody, err := json.MarshalIndent(room, "", " ")
	if err != nil {
		h.logger.Error("parse room to json", zap.Error(err))
		handleError(http.StatusInternalServerError, "parse room to json", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(roomBody)
}

// @Summary      Delete room by id
// @Description  delete room
// @Tags         room
// @Param token header string true "auth token"
// @Success      200  {string} string
// @Router       /room/{id} [delete]
func (h *Handler) deleteRoomById(w http.ResponseWriter, r *http.Request) {
	// verify user
	userId, err := verifyUserId(r)
	if err != nil {
		h.logger.Error("invalid token", zap.Error(err))
		handleError(http.StatusUnauthorized, "invalid token", err, w)
		return
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Error("invalid room id", zap.Error(err))
		handleError(http.StatusBadRequest, "invalid room id", err, w)
		return
	}

	err = h.service.DeleteRoomById(roomId, userId)
	if err != nil {
		log.Println(err)
		h.logger.Error("delete room by id", zap.Error(err), zap.String("room id", strconv.Itoa(roomId)), zap.String("user id", userId))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete successfully"))
}

func verifyUserId(r *http.Request) (string, error) {
	userToken, err := getJWTToken(r)
	if err != nil {
		return "", err
	}
	userIdStr, err := utils.VerifyToken(userToken, utils.UserId)
	if err != nil {
		return "", err
	}
	userId, ok := userIdStr.(string)
	if !ok {
		log.Println("invalid user id")
	}
	return userId, nil
}
