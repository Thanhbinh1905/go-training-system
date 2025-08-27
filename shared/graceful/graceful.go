package graceful

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Shutdownable interface for services that need graceful shutdown
type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

// GracefulServer manages graceful shutdown of multiple services
type GracefulServer struct {
	services []Shutdownable
	timeout  time.Duration
	wg       sync.WaitGroup
}

// NewGracefulServer creates a new graceful server
func NewGracefulServer(timeout time.Duration) *GracefulServer {
	return &GracefulServer{
		services: make([]Shutdownable, 0),
		timeout:  timeout,
	}
}

// AddService adds a service to be managed for graceful shutdown
func (g *GracefulServer) AddService(service Shutdownable) {
	g.services = append(g.services, service)
}

// Start starts the graceful shutdown listener
func (g *GracefulServer) Start() {
	// Create channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, starting graceful shutdown...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	// Shutdown all services concurrently
	g.wg.Add(len(g.services))
	for _, service := range g.services {
		go func(s Shutdownable) {
			defer g.wg.Done()
			if err := s.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down service: %v", err)
			}
		}(service)
	}

	// Wait for all services to shutdown or timeout
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All services shutdown successfully")
	case <-ctx.Done():
		log.Println("Shutdown timeout reached, forcing exit")
	}

	log.Println("Graceful shutdown completed")
}

// Wait waits for shutdown to complete
func (g *GracefulServer) Wait() {
	g.wg.Wait()
}

// SimpleShutdown provides a simple way to handle shutdown for basic services
func SimpleShutdown(timeout time.Duration, shutdownFuncs ...func(context.Context) error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(shutdownFuncs))

	for _, shutdownFunc := range shutdownFuncs {
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				log.Printf("Error during shutdown: %v", err)
			}
		}(shutdownFunc)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Graceful shutdown completed successfully")
	case <-ctx.Done():
		log.Println("Shutdown timeout reached")
	}
}
