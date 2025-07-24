package repository

import (
	"context"

	gqlmodel "github.com/Thanhbinh1905/go-training-system/internal/graph/model"
	"github.com/Thanhbinh1905/go-training-system/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FetchAll(ctx context.Context, role *model.UserRole) ([]*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, userID string, input *gqlmodel.UpdateUserInput) (*model.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FetchAll(ctx context.Context, role *model.UserRole) ([]*model.User, error) {
	var users []*model.User
	tx := r.db.WithContext(ctx).Model(&model.User{})
	if role != nil {
		tx = tx.Where("role = ?", *role)
	}
	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Update(ctx context.Context, userID string, input *gqlmodel.UpdateUserInput) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Role != nil {
		user.Role = model.UserRole(*input.Role)
	}

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
