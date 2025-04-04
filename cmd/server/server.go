package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"operation-borderless/internal/api"
	"operation-borderless/internal/infrastructure/postgres"
	serviceimplementation "operation-borderless/internal/service-implementation"
	"operation-borderless/pkg/config"
	"operation-borderless/pkg/db"
)

func Start() {
	cfg, err := config.InitConfigs()
	if err != nil {
		log.Fatalf("cannnot init config: %v", err)
	}
	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("cannnot init database: %v", err)
	}

	log.Println("after database init")

	repository := postgres.NewDB(database)

	forexAPI := serviceimplementation.NewExternalAPIClient(cfg)
	wallet := serviceimplementation.NewServiceClient(cfg, repository)

	h := api.NewHandler(forexAPI, wallet, cfg)

	r := SetupRouter(h)

	PORT := fmt.Sprintf(":%s", cfg.ServicePort)
	if PORT == ":" {
		PORT = ":8080"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", PORT)
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
