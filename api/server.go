package api

import (
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

func errorRespJSON(message string) gin.H {
	return gin.H{"message": message}
}
