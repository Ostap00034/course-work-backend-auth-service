package main

import (
	"log"
	"net"
	"os"
	"time"

	authpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	userpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Ostap00034/course-work-backend-auth-service/db"
	"github.com/Ostap00034/course-work-backend-auth-service/internal/auth"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	dbString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		log.Fatal("not DB_CONN_STRING in .env file")
	}
	client := db.NewClient(dbString)
	defer client.Close()

	userAddr, exists := os.LookupEnv("USER_ADDR")
	if !exists {
		log.Fatal("not USER_ADDR in .env file")
	}
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to user service: %v", err)
	}
	userClient := userpbv1.NewUserServiceClient(userConn)

	repo := auth.NewRepo(client)
	svc := auth.NewService(repo, userClient, 24*time.Hour)

	lis, _ := net.Listen("tcp", ":50051")
	srv := grpc.NewServer()
	authSrv := auth.NewAuthServer(svc)
	authpbv1.RegisterAuthServiceServer(srv, authSrv)

	log.Println("AuthService listening on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
