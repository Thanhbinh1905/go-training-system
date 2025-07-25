package graph

import "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service service.UserService
}
