package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mikeheft/go-backend/db/sqlc"
)

// Serve all http requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	router.POST("/accounts", server.createAccount)

	server.router = router

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
