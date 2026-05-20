package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"chat-service/internal/config"
	"chat-service/internal/dao"
	"chat-service/internal/domain/repositories"
	"chat-service/internal/grpc_handlers"
	"chat-service/internal/services"
	chatpb "proto/chat"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	chatAccessor := dao.NewChatAccessor(db)
	messageAccessor := dao.NewMessageAccessor(db)

	chatRepository := repositories.NewChatRepository(chatAccessor)
	messageRepository := repositories.NewMessageRepository(messageAccessor)

	chatService := services.NewChatService(
		chatRepository,
		messageRepository,
	)

	server := grpc_handlers.NewChatServer(chatService)

	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("%s started on %s", cfg.App.Name, addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
