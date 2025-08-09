package main

import (
	"fmt"
	"github.com/dhope-nagesh/titanic-go-service/internal/config"
	"github.com/dhope-nagesh/titanic-go-service/internal/data"
	"github.com/dhope-nagesh/titanic-go-service/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

// @title           Titanic Passenger API
// @version         1.0
// @description     A web service to query Titanic passenger data.
// @host            127.0.0.1:8080
// @BasePath        /api/v1
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	var repo data.PassengerRepository
	switch cfg.Data.Source {
	case "csv":
		repo, err = data.NewCSVRepository(cfg.Data.CSVFile)
		log.Println("Using CSV data source")
	case "sqlite":
		repo, err = data.NewSQLiteRepository(cfg.Data.DBFile)
		log.Println("Using SQLite data source")
	default:
		log.Fatalf("invalid data source in config: %s", cfg.Data.Source)
	}

	if err != nil {
		log.Fatalf("could not initialize repository: %v", err)
	}

	router := gin.Default()
	apiHandler := handler.NewAPIHandler(repo)
	apiHandler.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
