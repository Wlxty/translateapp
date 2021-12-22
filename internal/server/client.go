package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"translateapp/internal/libretranslate"
	"translateapp/internal/logger"
	"translateapp/internal/translateapp"
)

func GetServer() *http.Server {
	listenAddr := ":8080"
	logger := logger.NewLogger("debug", true)
	ltHost := "http://libretranslate:5000/"
	client := libretranslate.NewClient(logger, ltHost)
	service := translateapp.Service{
		logger,
		*client,
	}
	api := translateapp.NewServer(&service)
	api.HandleRequests(":8080")
	server := http.Server{
		Addr:         listenAddr,
		Handler:      api,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &server
}

func Run() error {
	server := GetServer()
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: API listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting server: %s", err)

	case <-shutdown:
		log.Println("main: Start shutdown")
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = server.Close()
		}
		if err != nil {
			return fmt.Errorf("main : could not stop server gracefully : %v", err)
		}
	}

	return nil
}
