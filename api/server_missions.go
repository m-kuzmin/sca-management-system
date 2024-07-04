package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db"
)

func (s *Server) CreateMission(ctx *gin.Context) {
	var targets []db.CreateTargetParams

	err := ctx.ShouldBindJSON(&targets)
	if err != nil && !errors.Is(err, io.EOF) { // body is optional
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to read JSON body"))
		return
	}

	for _, target := range targets {
		if target.Name == "" {
			ctx.JSON(http.StatusInternalServerError, errorRespJSON("must specify a name for a target"))
			return
		}

		if len(target.Country) == 0 || len(target.Country) > 3 {
			ctx.JSON(http.StatusInternalServerError, errorRespJSON("country must be a 3 letter code"))
			return
		}

		target.Country = strings.ToLower(target.Country)
	}

	id, err := s.db.CreateMissionWithTargets(ctx.Request.Context(), targets)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to create a mission"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (s *Server) GetMissionByID(ctx *gin.Context) {
	str := ctx.Query("id")
	id := uuid.UUID{}

	if id.UnmarshalText([]byte(str)) != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing id"))
		return
	}

	mission, err := s.db.GetMissionByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorRespJSON("failed to find mission with this id"))
		return
	}

	ctx.JSON(http.StatusFound, mission)
}

func (s *Server) ListMissionsPaginated(ctx *gin.Context) {
	str := ctx.Query("page")
	page, err := strconv.Atoi(str)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid page, should be more than 1",
		})
		return
	}

	str = ctx.Query("limit")
	limit, err := strconv.Atoi(str)
	// In this project the max number of cats is small, but an upper limit would be a good idea
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid limit, should be more than 1",
		})
		return
	}

	missions, err := s.db.ListMissions(ctx.Request.Context(), db.PaginationParams{PageNumber: uint32(page), Limit: uint32(limit)})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "failed to list missions",
		})
		return
	}

	ctx.JSON(http.StatusFound, missions)
}

func (s *Server) CompleteMission(ctx *gin.Context) {
	str := ctx.Query("id")
	id := uuid.UUID{}

	err := id.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing mission id"))
		return
	}

	err = s.db.CompleteMission(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("failed to mark mission as complete"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
