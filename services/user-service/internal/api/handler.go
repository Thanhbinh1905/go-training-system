package handler

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUserFromFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var jobs []dto.Job
	line := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		line++

		// Bỏ qua dòng header
		if line == 1 {
			continue
		}

		// Basic validation
		if len(record) < 4 {
			logger.Log.Warn("Invalid CSV format", zap.Int("line", line))
			continue
		}

		jobs = append(jobs, dto.Job{
			User: dto.UserCSV{
				Username: record[0],
				Email:    record[1],
				Password: record[2],
				Role:     model.UserRole(record[3]),
			},
			LineNumber: line,
		})
	}

	if len(jobs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid user rows found in CSV"})
		return
	}

	numWorkers := 5
	jobChan := make(chan dto.Job, len(jobs))
	resultChan := make(chan dto.Result, len(jobs))

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go h.worker(c.Request.Context(), jobChan, resultChan)
	}

	// Push jobs
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Collect result
	var summary dto.ImportResult
	var failedDetails []dto.Result

	for i := 0; i < len(jobs); i++ {
		res := <-resultChan
		if res.Success {
			summary.Success++
		} else {
			summary.Failed++
			failedDetails = append(failedDetails, res)
		}
	}

	logger.Log.Info("Import completed", zap.Int("success", summary.Success), zap.Int("failed", summary.Failed))

	c.JSON(http.StatusOK, gin.H{
		"success": summary.Success,
		"failed":  summary.Failed,
		"errors":  failedDetails,
	})
}

func (h *UserHandler) worker(ctx context.Context, jobChan <-chan dto.Job, resultChan chan<- dto.Result) {
	for job := range jobChan {
		input := &dto.CreateUserInput{
			Username: job.User.Username,
			Email:    job.User.Email,
			Password: job.User.Password,
			Role:     model.UserRole(job.User.Role),
		}

		_, err := h.service.Register(ctx, input)
		if err != nil {
			logger.Log.Error("Failed to create user",
				zap.String("email", job.User.Email),
				zap.Int("line", job.LineNumber),
				zap.Error(err),
			)
			resultChan <- dto.Result{
				Success:    false,
				Error:      err,
				LineNumber: job.LineNumber,
				Email:      job.User.Email,
			}
			continue
		}

		resultChan <- dto.Result{Success: true}
	}
}
