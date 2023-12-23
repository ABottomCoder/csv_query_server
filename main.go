package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"zh.com/ms_coding2/internal/utils"

	"zh.com/ms_coding2/internal/handler"
	"zh.com/ms_coding2/internal/repository"
	"zh.com/ms_coding2/internal/server"
)

func main() {
	filePath := utils.GetFilePath()
	if filePath == "" {
		filePath = repository.DefaultFilePath
	}
	p := os.Getenv("CSV_FILE_PATH")
	if p != "" {
		filePath = p
	}

	fmt.Printf("in main, filePath: %s\n", filePath)
	repository.InitFile(filePath)
	pth, _ := os.Getwd()
	fmt.Printf("path: %s\n", pth)
	queryServer := server.NewServer(":9527", handler.RegisterQueryHandler)

	go func() {
		err := queryServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server on port 9527:", err)
		}
	}()

	modifyServer := server.NewServer(":7259", handler.RegisterModifyHandler)

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
