package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"core-service/internal/config"
	"core-service/internal/dao"
	"core-service/internal/domain/repositories"
	"core-service/internal/grpc_handlers"
	"core-service/internal/services"
	corepb "proto/core"

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

	ticketAccessor := dao.NewTicketAccessor(db)
	specialistAccessor := dao.NewSpecialistAccessor(db)
	userAccessor := dao.NewUserAccessor(db)

	analysisAccessor, err := dao.NewAnalysisAccessor("localhost:50052")
	if err != nil {
		log.Fatalf("failed to connect analysis-service: %v", err)
	}
	defer analysisAccessor.Close()

	ticketRepository := repositories.NewTicketRepository(ticketAccessor)
	specialistRepository := repositories.NewSpecialistRepository(specialistAccessor)
	userRepository := repositories.NewUserRepository(userAccessor)

	ticketService := services.NewTicketService(
		ticketRepository,
		specialistRepository,
		analysisAccessor,
	)

	userService := services.NewUserService(userRepository, specialistRepository)

	server := grpc_handlers.NewServer(ticketService, userService)

	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	corepb.RegisterCoreServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("%s started on %s", cfg.App.Name, addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
