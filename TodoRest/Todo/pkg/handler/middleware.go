package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// certain middle layer for extracting some things from context (parsing)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// header is a way of sending additional info with k-v
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader) // получаем хэддер авторизации
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ") // разделяем хэддер на имя и значение
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1]) // парсим токен на основе полученной части хэддера
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId) // создаем в контексте пару ключ-значение для передачи возвращаемого id
}

func getUserId(c *gin.Context) (int, error) { // middleware for proper getting userId from gin context
	id, ok := c.Get(userCtx)
	if !ok { // error finding val with key userCtx in the context
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok { // error converting id to int
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
