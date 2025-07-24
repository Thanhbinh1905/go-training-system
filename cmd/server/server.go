package main

import (
	"log"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Thanhbinh1905/go-training-system/internal/config"
	"github.com/Thanhbinh1905/go-training-system/internal/graph"
	"github.com/Thanhbinh1905/go-training-system/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/internal/service"
	"github.com/Thanhbinh1905/go-training-system/pkg/db"
	"github.com/Thanhbinh1905/go-training-system/pkg/logger"
	"github.com/Thanhbinh1905/go-training-system/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load configuration")
	}

	logger.InitLogger(cfg.Production)

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	userRepo := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepo)
	resolver := &graph.Resolver{
		UserService: userService,
		JWTSecret:   cfg.JWTSecret,
	}
	srv := graphqlhandler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r := gin.Default()

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// GraphQL Playground
	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphQL").ServeHTTP(c.Writer, c.Request)
	})

	r.Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
	r.Use(middleware.ContextMiddleware())

	// GraphQL Query Handler
	r.POST("/graphQL", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// Protected routes group: yêu cầu auth
	authGroup := r.Group("/")
	authGroup.Use(middleware.RequiredAuthMiddleware(cfg.JWTSecret))

	teamRepo := repository.NewTeamRepository(conn)
	teamSvc := service.NewTeamService(teamRepo)
	teamHdl := handler.NewTeamHandler(teamSvc)

	// Routes cho team management chỉ dành cho manager
	teamGroup := authGroup.Group("/teams")
	teamGroup.Use(middleware.RequireManagerRole("manager"))
	{
		teamGroup.POST("/", teamHdl.CreateTeam)
		teamGroup.POST("/:teamId/members", teamHdl.AddMember)
		teamGroup.DELETE("/:teamId/members/:memberId", teamHdl.RemoveMember)
		teamGroup.POST("/:teamId/managers", teamHdl.AddManager)
		teamGroup.DELETE("/:teamId/managers/:managerId", teamHdl.RemoveManager)
	}

	logger.Log.Info("Starting server on port " + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}
