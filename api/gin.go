package api

import (
	"github.com/gin-gonic/gin"
)

func NewGinRouter(server *Server) *gin.Engine {
	router := gin.New()

	router.POST("/cat", server.CreateCat)
	router.GET("/cat", server.GetCatByID)
	router.DELETE("/cat", server.DeleteCatByID)

	return router
}
