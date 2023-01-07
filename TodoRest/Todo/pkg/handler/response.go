package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct { // кастомная структура для создания сообщений об ошибке
	Message string `json:"message"`
}

type StatusResponse struct { // кастомная структура для задания статуса
	Status string `json:"status"`
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)                              // использвуем logrus для логгирования
	ctx.AbortWithStatusJSON(statusCode, errorResponse{ // создаем ошибку с статусКодом и сообщением (формат json)
		Message: message,
	})
}
