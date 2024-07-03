package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db"
)

// Server holds various globlals such as db connections that are used to respond to http requests
type Server struct {
	db db.Querier
}

func NewServer(db db.Querier) *Server {
	return &Server{db: db}
}

func (s *Server) CreateCat(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "name is empty",
		})
		return
	}

	breed := ctx.Query("breed")
	if breed == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "breed is empty",
		})
		return
	}

	str := ctx.Query("years_of_experience")
	yearsOfExperience, err := strconv.Atoi(str)
	if err != nil || yearsOfExperience < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid years_of_experience, should be a positive integer",
		})
		return
	}

	str = ctx.Query("salary")
	sallary, err := strconv.Atoi(str)
	if err != nil || sallary < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid salary, should be a positive integer",
		})
		return
	}

	id, err := s.db.CreateCat(ctx.Request.Context(), name, uint16(yearsOfExperience), breed, uint(sallary))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create a cat",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (s *Server) GetCatByID(ctx *gin.Context) {
	str := ctx.Query("id")
	id := uuid.UUID{}

	err := id.UnmarshalText([]byte(str))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid or missing cat id",
		})
		return
	}

	cat, err := s.db.GetCatByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "failed to lookup the cat by id",
		})
		return
	}

	ctx.JSON(http.StatusFound, cat)
}
