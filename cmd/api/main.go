package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/internal/server"
	_ "github.com/lib/pq"
)

//		@title			Manka Api
//		@version		1.0
//		@description	Manga search
//	 @BasePath	/
func main() {
	srv := server.NewServer()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error listen server: %v\n", err)
		}
	}()
	slog.Info("Server is listening","port", srv.Addr())

	<-sigCh
	fmt.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	srv.Close(ctx)
	slog.Info("Server gracefully stopped")
	os.Exit(0)
}
