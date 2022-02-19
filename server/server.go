package server

import (
	"context"
	"fas/config"
	"fas/data"
	"fas/storage"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Server defines the full backend structure
type Server struct {
	// mongodb conn
	mongoClient *mongo.Client
	serverConf config.Server

	logger echo.Logger
}

// initialize the server instance
func (s *Server) init(e *echo.Echo)  {

	// init config
	conf, err := config.LoadConf()
	if err!=nil {
		e.Logger.Fatalf("failed to load config: %v", err)
	}
	s.serverConf = conf.Server
	s.logger = e.Logger

	// init db connection
	client, err := storage.NewClient(conf.Mongodb)
	if err!=nil {
		e.Logger.Fatalf("failed to connect mongodb: %v", err)
	}
	s.mongoClient = client

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// plug in validator
	e.Validator = data.NewValidator()
	// routes registration
	e.GET("/flow", s.AggregateFlowDataByHour)
	e.POST("/flow", s.ReportNetworkFlow)
}

// shutdown hook to gracefully disconnect from mongodb
func (s *Server) shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.mongoClient.Disconnect(ctx)
}

// Start starts the application
func Start()  {
	e := echo.New()
	server := Server{}
	server.init(e)

	addr := fmt.Sprintf("%s:%d", server.serverConf.Host, server.serverConf.Port)

	defer server.shutdown()
	e.Logger.Fatal(e.Start(addr))
}