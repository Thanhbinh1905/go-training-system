package app

import (
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/client"
	grpcHandler "github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler/grpc"
	httpHandler "github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler/http"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/middleware"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	teampb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/team"
	userpb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/user"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pkg/logger"
	"github.com/Thanhbinh1905/go-training-system/shared/db"

	"net"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const userGRPCURL = "user-service:50051"

func Run(cfg *config.Config) {
	log := logger.NewLogger("logs/team-service.log", "team-service")
	defer log.Sync()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	teamRepo := repository.NewTeamRepository(conn)
	userClient := client.NewUserGRPCClient(userGRPCURL)
	teamService := service.NewTeamService(teamRepo, *userClient)

	go runGRPCServer(teamService, cfg.GRPCPort, log)
	runHTTPServer(teamService, *userClient, log)
}

func runGRPCServer(teamService service.TeamService, port string, log *zap.Logger) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}
	server := grpc.NewServer()
	teampb.RegisterTeamServiceServer(server, grpcHandler.NewTeamgRPCHandler(teamService))

	log.Info("gRPC server listening", zap.String("port", port))
	if err := server.Serve(lis); err != nil {
		log.Fatal("failed to serve gRPC", zap.Error(err))
	}
}

func runHTTPServer(teamService service.TeamService, userClient client.UserGRPCClient, log *zap.Logger) {
	r := gin.Default()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	// Metrics
	p := ginprometheus.NewPrometheus("team_service")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	teamHandler := httpHandler.NewTeamHandler(teamService)

	api := r.Group("/api/v1")
	teamGroup := api.Group("/teams")
	teamGroup.Use(middleware.AuthMiddleware(userClient.Client, userpb.Role_MANAGER))
	{
		teamGroup.POST("/", teamHandler.CreateTeam)
		teamGroup.POST("/:teamID/managers", teamHandler.AddManager)
		teamGroup.POST("/:teamID/members", teamHandler.AddMember)
		teamGroup.DELETE("/:teamID/managers/:managerID", teamHandler.RemoveManager)
		teamGroup.DELETE("/:teamID/members/:memberID", teamHandler.RemoveMember)
	}

	log.Info("HTTP server started", zap.String("port", "8080"))
	r.Run(":8080")
}
