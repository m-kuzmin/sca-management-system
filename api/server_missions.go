package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db"
)

/*
Missions / Targets
Ability to delete a mission
A mission cannot be deleted if it is already assigned to a cat
Ability to update mission
Ability to mark it as completed
Ability to update mission targets
Ability to mark it as completed
Ability to update Notes
Notes cannot be updated if either the target or the mission is completed
Ability to delete targets from an existing mission
A target cannot be deleted if it is already completed
Ability to add targets to an existing mission
A target cannot be added if the mission is already completed
Ability to assign a cat to a mission
Ability to list missions
Ability to get a single mission

*/

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
