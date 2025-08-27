package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/graph"
	httpHandler "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/http"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/util/token"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/pb"
	"github.com/Thanhbinh1905/go-training-system/shared/db/postgres"
	"github.com/Thanhbinh1905/go-training-system/shared/graceful"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/service"
	"google.golang.org/grpc"

	grpcHandler "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/grpc"
	ginzap "github.com/gin-contrib/zap"
)

const (
	GRATEFUL_TIMEOUT          = 30 * time.Second
	ACCESS_TOKEN_EXPIRE_TIME  = 24 * time.Hour
	REFRESH_TOKEN_EXPIRE_TIME = 7 * 24 * time.Hour
)

// Defining the Graphql handler
func graphqlHandler(userService service.UserService) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Service: userService,
	}}))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func RunGRPCServer(userService service.UserService, port string) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcHandler.NewUserRPCHandler(userService))

	log.Println("gRPC server listening on :50051")

	// Start server in goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	return grpcServer, lis, nil
}

func RunHTTPServer(userService service.UserService, log *zap.Logger) (*http.Server, error) {
	userHandler := httpHandler.NewUserHandler(userService)
	gqlHandler := graphqlHandler(userService)

	r := gin.Default()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	// GraphQL
	r.POST("/graphql", gqlHandler)
	r.GET("/", playgroundHandler())

	// REST
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		v1.POST("/users", userHandler.CreateUserFromFile)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Info("Starting HTTP server", zap.String("port", "8080"))

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to run HTTP server", zap.Error(err))
		}
	}()

	return server, nil
}

func Run(cfg *config.Config) {
	// Init Logger
	log := logger.InitLogger("logs/user-service.log", "user-service")
	defer log.Sync()

	// DB connection
	conn, err := postgres.Connect(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal("DB connection failed", zap.Error(err))
	}
	defer postgres.Close(conn, log)

	// Init dependencies
	userRepo := repository.NewUserRepository(conn)
	jwtManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTSecret, ACCESS_TOKEN_EXPIRE_TIME, REFRESH_TOKEN_EXPIRE_TIME)
	userService := service.NewUserService(userRepo, jwtManager)

	// Start gRPC server
	grpcServer, grpcListener, err := RunGRPCServer(userService, cfg.GRPCPort)
	if err != nil {
		log.Fatal("Failed to start gRPC server", zap.Error(err))
	}

	// Start HTTP server
	httpServer, err := RunHTTPServer(userService, log)
	if err != nil {
		log.Fatal("Failed to start HTTP server", zap.Error(err))
	}

	// Create graceful shutdown server
	gracefulServer := graceful.NewGracefulServer(GRATEFUL_TIMEOUT)

	// Add services for graceful shutdown
	gracefulServer.AddService(&gracefulService{
		name:         "gRPC Server",
		grpcServer:   grpcServer,
		grpcListener: grpcListener,
		httpServer:   httpServer,
		log:          log,
	})

	// Start graceful shutdown listener
	gracefulServer.Start()
}

// gracefulService implements Shutdownable interface
type gracefulService struct {
	name         string
	grpcServer   *grpc.Server
	grpcListener net.Listener
	httpServer   *http.Server
	log          *zap.Logger
}

func (s *gracefulService) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down " + s.name)

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Shutdown gRPC server
	s.grpcServer.GracefulStop()

	// Close gRPC listener
	if err := s.grpcListener.Close(); err != nil {
		s.log.Error("gRPC listener close error", zap.Error(err))
	}

	s.log.Info(s.name + " shutdown completed")
	return nil
}
