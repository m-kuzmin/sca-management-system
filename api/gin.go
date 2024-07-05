package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(server *Server) *gin.Engine {
	router := gin.New()

	g := router.Group("/")
	g.Use(LoggerMiddleware())

	cat := g.Group("cat")

	cat.POST("/", server.CreateCat)
	cat.GET("by-id", server.GetCatByID)
	cat.GET("list", server.ListCatsPaginated)
	cat.POST("salary", server.UpdateCatSalaryByID)
	cat.DELETE("/", server.DeleteCatByID)

	mission := g.Group("mission")

	mission.POST("/", server.CreateMission)
	mission.GET("by-id", server.GetMissionByID)
	mission.GET("list", server.ListMissionsPaginated)
	mission.POST("complete", server.CompleteMission)
	mission.POST("target", server.AddTargetsToMission)
	mission.POST("assign", server.AssignCatToMission)
	mission.DELETE("/", server.DeleteMission)

	target := g.Group("target")

	target.POST("complete", server.CompleteTarget)
	target.POST("notes", server.UpdateTargetNotes)
	target.DELETE("/", server.DeleteTarget)

	return router
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		latency := end.Sub(start)

		log.Printf("[GIN] %s %s -> %d (%v)\n",
			c.Request.Method,
			c.Request.URL,
			c.Writer.Status(),
			latency,
		)
	}
}
