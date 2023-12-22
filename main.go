package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"zh.com/ms_coding2/server"
)

func main() {
	server.InitFile()
	queryServer := server.NewServer(":9527", server.RegisterQueryHandler)

	go func() {
		err := queryServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server on port 9527:", err)
		}
	}()

	modifyServer := server.NewServer(":7259", server.RegisterModifyHandler)

	go func() {
		err := modifyServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server on port 7259:", err)
		}
	}()

	fmt.Println("Server started")

	// shutdown server when receive signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	fmt.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown server
	err := queryServer.Shutdown(ctx)
	if err != nil {
		fmt.Println("Failed to  shutdown server on port 9527:", err)
	}

	err = modifyServer.Shutdown(ctx)
	if err != nil {
		fmt.Println("Failed to  shutdown server on port 7259:", err)
	}

	fmt.Println("Server successfully stopped")
}
