package main

import (
	"auth-service/config"
	pb "auth-service/generated/auth_service"
	"auth-service/service"
	"auth-service/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config.Load()
	fmt.Println(cfg.GRPC_PORT)
	listener, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	service := service.NewAuthService(postgres.NewUserRepo(db))

	pb.RegisterAuthServiceServer(s, service)

	log.Printf("server is running on %v...", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}