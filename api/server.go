package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "years_of_experience is not a number",
		})
		return
	}

	str = ctx.Query("salary")
	sallary, err := strconv.Atoi(str)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "salary is not a number",
		})
		return
	}

	id, err := s.db.CreateCat(context.Background(), name, uint16(yearsOfExperience), breed, uint(sallary))
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
