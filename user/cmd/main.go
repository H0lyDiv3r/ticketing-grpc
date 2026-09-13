package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"buf.build/go/protovalidate"
	upb "github.com/H0lyDiv3r/ticketing-grpc/gen/user"
	"github.com/H0lyDiv3r/ticketing-grpc/user/domain"
	"github.com/H0lyDiv3r/ticketing-grpc/user/handler"
	"github.com/H0lyDiv3r/ticketing-grpc/user/interceptors"
	"github.com/H0lyDiv3r/ticketing-grpc/user/repository"
	"github.com/H0lyDiv3r/ticketing-grpc/user/service"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func init() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "user_app")
	password := getEnv("DB_PASSWORD", "secretpassword")
	dbname := getEnv("DB_NAME", "user_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	db.AutoMigrate(&domain.User{})

	log.Println("Database connection established successfully")
}

func main() {

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create validator: %v", err)
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(*repo)
	handler := handler.NewHandler(svc)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen on :50051: %v\n", err)
		return
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.ValidationInterceptors(validator), interceptors.LoggerInterceptor))
	upb.RegisterUserServiceServer(grpcServer, handler)

	fmt.Println("User gRPC server is listening on :50051...")
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}

	log.Println("User service initialized")
}
