package app

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/graph"
	htppHandler "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/http"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/util/token"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/pb"
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
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"google.golang.org/grpc"

	grpcHandler "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/grpc"
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

func RunGRPCServer(userService service.UserService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcHandler.NewUserRPCHandler(userService))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func RunHTTPServer(userService service.UserService) {
	userHandler := htppHandler.NewUserHandler(userService)
	gqlHandler := graphqlHandler(userService)

	r := gin.Default()

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

	logger.Log.Info("HTTP server listening", zap.String("port: ", "8080"))
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("failed to start HTTP server", zap.Error(err))
	}
}

func Run(cfg *config.Config) {
	// Init Logger
	logger.InitLogger(cfg.Production)
	defer logger.Log.Sync()

	// DB connection
	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("DB connection failed", zap.Error(err))
	}
	defer db.Close(conn)

	// Init dependencies
	userRepo := repository.NewUserRepository(conn)
	jwtManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTSecret, 24*time.Hour, 7*24*time.Hour)
	userService := service.NewUserService(userRepo, jwtManager)

	// Run gRPC in background
	go RunGRPCServer(userService, cfg.GRPCPort)

	// Start HTTP server (GraphQL + REST)
	RunHTTPServer(userService)
}
