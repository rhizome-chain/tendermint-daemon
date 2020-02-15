package api

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"log"
	
	"github.com/gin-gonic/gin"
)

var VERSION = "v1"

type API interface {
	RelativePath() string
	SetHandlers(group *gin.RouterGroup)
}

// Server ..
type Server struct {
	router         *gin.Engine
	err            chan error
}

// NewServer create new API Server
func NewServer(dm *daemon.Daemon) (server *Server) {
	server = new(Server)
	server.err = make(chan error)
	server.router = gin.Default()
	
	dmapi := NewDaemonAPI(dm)
	server.AddAPI(dmapi)
	
	
	return server
}

func (server *Server) Error() <-chan error {
	return server.err
}

// Start ..
func (server *Server) Start(listenAddress string) {
	go func() {
		err := server.router.Run(listenAddress)
		if err != nil {
			log.Fatal("Cannot Start API Server")
		}
	}()
}

// Group : delegate *gin.RouterGroup.Group
func (server *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return server.router.Group(relativePath, handlers...)
}

func (server *Server) AddAPI(api API) {
	group := server.router.Group(VERSION + "/" + api.RelativePath())
	api.SetHandlers(group)
}
