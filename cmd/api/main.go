package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IlyaAGL/avito_autumn_2025/internal/app/handler"
	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/service"
	"github.com/IlyaAGL/avito_autumn_2025/internal/infrastructure/persistence/postgres"
	"github.com/IlyaAGL/avito_autumn_2025/pkg/bootstrap/connections"
	"github.com/gin-gonic/gin"
)

const timeToCloseServerConnection = 30

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	db_pg := connections.InitPostgres()
	defer func() {
		log.Fatal(db_pg.Close())
	}()

	userRepo := postgres.NewPostgresUserRepository(db_pg)
	prRepo := postgres.NewPostgresPullRequestRepository(db_pg)
	teamRepo := postgres.NewPostgresTeamRepository(db_pg)

	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPullRequestService(prRepo, userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo)

	userHandler := handler.NewUserHandler(userService)
	prHandler := handler.NewpullRequestHandler(prService)
	teamHandler := handler.NewTeamHandler(teamService)

	r := gin.Default()

	teams := r.Group("/team")
	{
		teams.GET("/get", teamHandler.GetTeam)
		teams.POST("/add", teamHandler.AddTeam)
		teams.POST("/bulk", teamHandler.BulkDeactivateUsers)
	}

	users := r.Group("/users")
    {
        users.POST("/setIsActive", userHandler.SetIsActive)
        users.GET("/getReview", userHandler.GetReview)
    }

	prs := r.Group("/pullRequest")
    {
        prs.POST("/create", prHandler.CreatePR)
        prs.POST("/merge", prHandler.MergePR)
        prs.POST("/reassign", prHandler.ReassignReviewer)
        prs.GET("/statistics", prHandler.GetStats)
    }

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeToCloseServerConnection*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error while shutdowning: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("Forced to close the server: %v", err)
		}
	}
}
