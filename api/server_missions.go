package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db"
)

const MaxTargets = 3

func (s *Server) CreateMission(ctx *gin.Context) {
	var targets []db.CreateTargetParams

	err := ctx.ShouldBindJSON(&targets)
	if err != nil && !errors.Is(err, io.EOF) { // body is optional
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to read JSON body"))
		return
	}

	if len(targets) > MaxTargets { // TODO: should it be len > 1 ?
		ctx.JSON(http.StatusBadRequest, errorRespJSON("failed to read JSON body"))
		return
	}

	for _, target := range targets {
		if !validateAndNormalizeTargetCreateParams(&target) {
			ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid target create data"))
			return
		}
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

func (s *Server) AddTargetsToMission(ctx *gin.Context) {
	str := ctx.Query("id")
	missionID := uuid.UUID{}

	if missionID.UnmarshalText([]byte(str)) != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing mission id"))
		return
	}

	var targets []db.CreateTargetParams

	err := ctx.ShouldBindJSON(&targets)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("body contains invalid JSON"))
		return
	}

	for _, target := range targets {
		if !validateAndNormalizeTargetCreateParams(&target) {
			ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid target create data"))
			return
		}
	}

	currentlyTargets, err := s.db.CountMissionTargets(ctx.Request.Context(), missionID)
	log.Println(currentlyTargets)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to count existing targets"))
		return
	}

	if len(targets) > MaxTargets-int(currentlyTargets) { // TODO: should it be len > 1 ?
		ctx.JSON(http.StatusBadRequest, errorRespJSON(fmt.Sprintf("a mission can only have a maximum %d targets", MaxTargets)))
		return
	}

	targetIDs, err := s.db.AddTargetsToMission(ctx.Request.Context(), missionID, targets)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to add targets to mission"))
		return
	}

	ctx.JSON(http.StatusCreated, targetIDs)
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

func (s *Server) CompleteTarget(ctx *gin.Context) {
	str := ctx.Query("id")
	id := uuid.UUID{}

	err := id.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing target id"))
		return
	}

	err = s.db.CompleteTarget(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("failed to mark target as complete"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (s *Server) UpdateTargetNotes(ctx *gin.Context) {
	str := ctx.Query("id")
	id := uuid.UUID{}

	err := id.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing target id"))
		return
	}

	var notes struct {
		Notes string `json:"notes"`
	}

	err = ctx.ShouldBindJSON(&notes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid body JSON"))
		return
	}

	isComplete, err := s.db.GetTargetCompleteStatus(ctx.Request.Context(), id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to get target's complete status"))
		return
	}

	if isComplete {
		ctx.JSON(http.StatusForbidden, errorRespJSON("target is already completed"))
		return
	}

	err = s.db.UpdateTargetNotes(ctx, id, notes.Notes)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorRespJSON("failed to update target notes"))
		return
	}

	ctx.JSON(http.StatusOK, successTrue())
}

func (s *Server) AssignCatToMission(ctx *gin.Context) {
	str := ctx.Query("mission")
	mission := uuid.UUID{}

	err := mission.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing mission id"))
		return
	}

	str = ctx.Query("cat")
	cat := uuid.UUID{}

	err = cat.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRespJSON("invalid or missing cat id"))
		return
	}

	err = s.db.AssignCatToMission(ctx.Request.Context(), db.AssignCatToMissionParams{
		Mission: mission,
		Cat:     cat,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRespJSON("failed to assign a cat to a mission"))
		return
	}

	ctx.JSON(http.StatusOK, successTrue())
}

func validateAndNormalizeTargetCreateParams(target *db.CreateTargetParams) bool {
	if target.Name == "" {
		return false
	}

	if len(target.Country) != 3 {
		return false
	}

	target.Country = strings.ToLower(target.Country)
	return true
}
