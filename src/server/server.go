package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	Addr    string
	Server  *http.Server
	Handler *gin.Engine
}

// DefaultNewApp returns a new App with default configurations.
//
// No parameters.
// Returns a pointer to an App.
func DefaultNewApp() *App {
	handler := gin.Default()
	return &App{
		Addr:    ":8080",
		Handler: handler,
	}
}

// NewApp creates a new App instance.
//
// No parameters.
// Returns a pointer to an App.
func NewApp() *App {
	handler := gin.New()
	return &App{
		Addr:    ":8080",
		Handler: handler,
	}
}

// ErrorHandler defines the error handling middleware for the App.
//
// No parameters.
// No return type.
func (a *App) ErrorHandler() {
	a.Handler.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{"error": "Page not found"})
		return
	})
	a.Handler.Use(func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
				return
			}
		}()
		c.Next()
	})
}

// Run runs the App server.
//
// No parameters.
// No return types.
func (a *App) Run() {
	a.Server = &http.Server{
		Addr:    a.Addr,
		Handler: a.Handler,
	}
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()
}

// Close closes the App, shutting down the server gracefully.
//
// No parameters.
// No return values.
func (a *App) Close() {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %s\n", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}

// Register registers a handler for the specified prefix.
//
// It takes a prefix string and a handler function, and does not return anything.
func (a *App) Register(prefix string, handler func(c *gin.RouterGroup)) {
	handler(a.Handler.Group(prefix))
}
