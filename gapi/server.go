package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/mikeheft/go-backend/pb"
	"github.com/mikeheft/go-backend/token"
	"github.com/mikeheft/go-backend/util"
	"github.com/mikeheft/go-backend/worker"
)

// Serve all gRPC requests
type Server struct {
	pb.UnimplementedSimpleBankServer

	config          util.Config
	router          *gin.Engine
	store           db.Store
	taskDistributer worker.TaskDistributor
	tokenMaker      token.Maker
}

func NewServer(config util.Config, store db.Store, taskDistributer worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributer: taskDistributer,
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
