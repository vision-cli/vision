package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atos-digital/NHSS-scigateway/internal/config"
	"github.com/atos-digital/NHSS-scigateway/internal/server"
)

func main() {
	conf := config.New()
	srv, err := server.New(conf)
	if err != nil {
		log.Fatalf("server: %v\n", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here, databases etc
		// TODO(steve): this does not seem to trigger
		log.Println("Shutting down DB connection", srv.TeardownDB())
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
