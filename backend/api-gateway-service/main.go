package main

import (
	"log"
	"net/http"

	"api-gateway-service/internal/controller"
	"api-gateway-service/internal/controller/http_handlers"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/ws"

	"github.com/go-chi/chi/v5"
)

func main() {
	coreClient, err := grpc_client.NewCoreClient("localhost", 50051)
	if err != nil {
		log.Fatalf("failed to connect to core-service: %v", err)
	}
	defer func(coreClient *grpc_client.CoreClient) {
		err := coreClient.Close()
		if err != nil {
			log.Printf("failed to close core client: %v", err)
		}
	}(coreClient)

	chatClient, err := grpc_client.NewChatClient("localhost", 50089)
	if err != nil {
		log.Fatalf("failed to connect to chat-service: %v", err)
	}
	defer func(chatClient *grpc_client.ChatClient) {
		err := chatClient.Close()
		if err != nil {
			log.Printf("failed to close chat client: %v", err)
		}
	}(chatClient)

	hub := ws.NewHub()

	router := chi.NewRouter()
	router.Use(corsMiddleware)

	controller.RegisterHandlers(
		router,
		http_handlers.NewCreateTicketHandler(coreClient),
		http_handlers.NewListUserTicketsHandler(coreClient),
		http_handlers.NewSpecialistListTicketsHandler(coreClient),
		http_handlers.NewGetTicketHandler(coreClient),
		http_handlers.NewAssignTicketHandler(coreClient),
		http_handlers.NewCloseTicketHandler(coreClient),

		http_handlers.NewGetChatMessagesHandler(chatClient),
		http_handlers.NewOpenTicketChatHandler(chatClient),
		http_handlers.NewSendMessageHandler(chatClient),

		http_handlers.NewChatWSHandler(chatClient, hub),

		http_handlers.NewRegisterHandler(coreClient),
		http_handlers.NewLoginHandler(coreClient),
		http_handlers.NewGetUserHandler(coreClient),
	)

	log.Println("api-gateway started on :8078")

	if err := http.ListenAndServe(":8078", router); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
