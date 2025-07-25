package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Fetch(ctx context.Context, role *model.UserRole, limit, offset int32) (*model.PaginatedUsers, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
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

func (r *userRepository) Fetch(ctx context.Context, role *model.UserRole, limit, offset int32) (*model.PaginatedUsers, error) {
	var users []*model.User

	query := r.db.WithContext(ctx).Model(&model.User{})

	if role != nil {
		query.Where("role = ?", *role)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Limit(int(limit)).Offset(int(offset)).Find(&users).Error; err != nil {
		return nil, err
	}

	return &model.PaginatedUsers{
		Users:  users,
		Total:  int32(total),
		Limit:  limit,
		Offset: offset,
	}, nil
}
