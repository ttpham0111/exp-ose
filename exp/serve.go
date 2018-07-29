package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ttpham0111/exp-ose/exp/database"
	"github.com/ttpham0111/exp-ose/exp/services"

	"github.com/ttpham0111/exp-ose/exp/v1/events"
	"github.com/ttpham0111/exp-ose/exp/v1/experiences"
)

const (
	version = "0.0.1"
)

type Server struct {
	conf *config
	db   *database.Database
	yelp services.YelpService
}

func NewServer() (*Server, error) {
	conf := newConfig()
	db, err := database.NewDatabase(conf.dbUrl, conf.dbName)
	if err != nil {
		return nil, err
	}

	return &Server{
		conf: conf,
		db:   db,
		yelp: services.NewYelpService(conf.yelpApiKey),
	}, nil
}

func (server *Server) newRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"version": version,
		})
	})

	events.Register(router.Group("/v1/events"), &events.Service{
		Collection: server.db.EventCollection,
		Yelp:       server.yelp,
	})

	experiences.Register(router.Group("/v1/experiences"), &experiences.Service{
		ExperienceCollection: server.db.ExperienceCollection,
		EventCollection:      server.db.EventCollection,
	})

	return router
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer server.db.Close()

	log.Println("Serving on localhost:" + server.conf.port)
	log.Fatal(server.newRouter().Run(":" + server.conf.port))
}
