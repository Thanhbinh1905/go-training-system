package main

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

	userRPC "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/handler/grpc"
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
	pb.RegisterUserServiceServer(grpcServer, userRPC.NewUserRPCHandler(userService))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.InitLogger(cfg.Production)
	defer logger.Log.Sync()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	userRepo := repository.NewUserRepository(conn)
	token := token.NewJWTManager(cfg.JWTSecret, cfg.JWTSecret, time.Hour*24, time.Hour*24*7)
	userService := service.NewUserService(userRepo, token)

	go RunGRPCServer(userService, cfg.GRPCPort)

	userHandler := htppHandler.NewUserHandler(userService)
	gqlHandler := graphqlHandler(userService)

	// Setting up Gin
	r := gin.Default()
	r.POST("/graphql", gqlHandler)
	r.GET("/", playgroundHandler())
	v1 := r.Group("/api/v1") // ✅ Tạo group

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1.POST("/users", userHandler.CreateUserFromFile)
	r.Run(":8080")
}
