package http

import (
	"log"
	"net/http"
	"time"

	"github.com/dosanma1/bluelabs_assessment/config"
	"github.com/dosanma1/bluelabs_assessment/internal/server"
	"github.com/dosanma1/bluelabs_assessment/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type httpServer struct {
	log     *logrus.Logger
	cfg     *config.Config
	pgxPool *pgxpool.Pool
}

func NewHttpServer(log *logrus.Logger, cfg *config.Config, pgxPool *pgxpool.Pool) server.Server {
	return &httpServer{
		log:     log,
		cfg:     cfg,
		pgxPool: pgxPool,
	}
}

func (s *httpServer) Run() {
	server := &http.Server{
		Addr:           s.cfg.HTTP.Port,
		Handler:        setupRouter(),
		ReadTimeout:    s.cfg.HTTP.ReadTimeout * time.Second,
		WriteTimeout:   s.cfg.HTTP.WriteTimeout * time.Second,
		IdleTimeout:    s.cfg.HTTP.MaxConnectionIdle * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(s.log.Writer(), "", 0),
	}
	if err := server.ListenAndServe(); err != nil {
		s.log.Fatalln(err)
	}

	s.log.Println("GRPC Server started!")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})

	return r
}
