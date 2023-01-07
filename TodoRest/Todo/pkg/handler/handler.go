package handler

import (
	"Todo/pkg/service"
	"github.com/gin-gonic/gin"
)

// здесь мы смотрим на все возможные обращения к хосту (все параметры из строки)

type Handler struct { // типа dependency injection. На самом верху иерархии - хэндлер, он обрабатывает запросы к сервису,
	// вызывает необходимые функции в зависимости от текста в строке
	services *service.Service
}

func NewHandler(services *service.Service) *Handler { // конструктор (возвращает указатель)
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}
	return router
}
