package main

import (
	"amanj/trustwallet/ethparser"
	"amanj/trustwallet/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	PollingInterval time.Duration = 5
	HTTPPort        string        = ":8080"
)

func main() {
	// creates a new parser which starts polling periodically from ethereum network
	parser := ethparser.NewEthParser()

	// Polls Parser for new transactions
	go InitEthBlockChainPolling(parser)

	// Initializes the services
	svc := services.NewServices(parser)

	// create http server with routes to interact with parser
	httpServer := NewHttpServer(svc)

	// Start the server in a goroutine
	go InitHttpServer(httpServer)
	log.Printf("HTTP server started ion port %v", HTTPPort)

	// Create a channel to listen for interrupt or terminate signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %s. Initiating graceful shutdown...\n", sig)

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v\n", err)
	} else {
		log.Println("Server shutdown gracefully.")
	}
}

func InitEthBlockChainPolling(parser ethparser.Parser) {
	for {
		log.Println("Polling Eth Blockchain")
		err := parser.PollBlockchain()
		if err != nil {
			log.Print("Error Polling Eth Blockchain: %v\n", err)
		}

		time.Sleep(PollingInterval * time.Second)
	}
}

func InitHttpServer(httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port 8080: %v\n", err)
	}
}

func NewHttpServer(svc *services.Services) *http.Server {
	mux := http.NewServeMux()

	// route to fetch current block in the block chain
	mux.HandleFunc("/current_block", svc.GetCurrentBlock)
	// route to subscribe an address
	mux.HandleFunc("/subscribe", svc.Subscribe)
	//  fetch transaction for an address
	mux.HandleFunc("/transactions", svc.GetTransactions)

	return &http.Server{
		Addr:    HTTPPort,
		Handler: mux,
	}
}
