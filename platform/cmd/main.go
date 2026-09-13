package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	upb "github.com/H0lyDiv3r/ticketing-grpc/gen/user"
	"github.com/H0lyDiv3r/ticketing-grpc/platform/handler"
	"github.com/H0lyDiv3r/ticketing-grpc/platform/service"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	userSvcURL := getEnv("USER_SERVICE_URL", "localhost:50051")
	port := getEnv("PORT", "8080")

	conn, err := grpc.NewClient(userSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service at %s: %v", userSvcURL, err)
	}
	defer conn.Close()

	userClient := upb.NewUserServiceClient(conn)
	svc := service.NewService(userClient)
	h := handler.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/signup", h.Signup)

	log.Printf("Platform HTTP server listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
