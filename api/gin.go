package api

import (
	"github.com/gin-gonic/gin"
)

func NewGinRouter(server *Server) *gin.Engine {
	router := gin.New()

	router.POST("/cat", server.CreateCat)
	router.GET("/cat/by-id", server.GetCatByID)
	router.GET("/cat/list", server.ListCatsPaginated)
	router.POST("/cat/salary", server.UpdateCatSalaryByID)
	router.DELETE("/cat", server.DeleteCatByID)

	router.POST("/mission", server.CreateMission)
	router.GET("/mission/by-id", server.GetMissionByID)
	router.GET("/mission/list", server.ListMissionsPaginated)
	router.POST("/mission/complete", server.CompleteMission)
	router.POST("/mission/target", server.AddTargetsToMission)
	router.POST("/mission/assign", server.AssignCatToMission)
	router.DELETE("/mission", server.DeleteMission)

	router.POST("/target/complete", server.CompleteTarget)
	router.POST("/target/notes", server.UpdateTargetNotes)
	router.DELETE("/target/", server.DeleteTarget)

	return router
}
