package handler

import (
	todo "Todo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil { // проверяем нормально ли чаитаются данные из контекста
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input) // выполняем операцию создания юзера, передавая контекст дальше
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{ // все получилось, передаем статус ОК и мапу с id
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"` // (связывание 	)
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil { // получаем контекс и пробуем его спарсить в структуру
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // создаем токен
	// на основе полученных из контекста данных о юзернейме и пароле
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{ // все получилось, передаем созданный токен
		"token": token,
	})
}
