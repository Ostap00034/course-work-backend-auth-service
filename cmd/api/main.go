package main

import (
	"log"
	"net"
	"os"
	"time"

	pb "github.com/Ostap00034/course-work-backend/auth-service/api/auth/v1"
	userpb "github.com/Ostap00034/course-work-backend/user-service/api/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Ostap00034/course-work-backend/auth-service/db"
	authrepo "github.com/Ostap00034/course-work-backend/auth-service/internal/repo"
	authsrv "github.com/Ostap00034/course-work-backend/auth-service/internal/server"
	authsvc "github.com/Ostap00034/course-work-backend/auth-service/internal/service"
)

func main() {
	// 1) Инициализируем БД и миграции
	client := db.NewClient(os.Getenv("DB_CONN_STRING"))

	// 2) gRPC-клиент для UserService
	userConn, err := grpc.NewClient(os.Getenv("USER_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to user service: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	// 3) Репозитории и сервис
	repo := authrepo.NewRepo(client)
	svc := authsvc.NewService(repo, userClient, 24*time.Hour)

	// 4) Запуск gRPC-сервера
	lis, _ := net.Listen("tcp", ":50051")
	srv := grpc.NewServer()
	authSrv := authsrv.NewAuthServer(svc)
	pb.RegisterAuthServiceServer(srv, authSrv)

	log.Println("AuthService listening on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
