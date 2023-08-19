package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/slavik22/imageAPI/db/sqlc"
	"github.com/slavik22/imageAPI/token"
	"github.com/slavik22/imageAPI/util"
)

type Server struct {
	config   util.Config
	store    db.Store
	router   *gin.Engine
	jwtMaker token.JWTMaker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	jwt, err := token.NewJWTMaker(config.SecretKey)

	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker %w", err)
	}

	server := &Server{
		store:    store,
		jwtMaker: *jwt,
		config:   config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", server.loginUser)
	router.POST("/register", server.createUser)

	_ = router.Group("/").Use(authMiddleware(server.jwtMaker))

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
