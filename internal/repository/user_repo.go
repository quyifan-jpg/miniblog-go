package repository

import (
	"miniblog/internal/model"
	"miniblog/internal/pkg/mysql"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return mysql.DB.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := mysql.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return mysql.DB.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id int64) error {
	return mysql.DB.Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(offset, limit int) ([]*model.User, error) {
	var users []*model.User
	err := mysql.DB.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// Search 搜索用户
func (r *UserRepository) Search(keyword string, offset, limit int) ([]*model.User, error) {
	var users []*model.User
	err := mysql.DB.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// GetByIDs 根据 ID 列表获取用户
func (r *UserRepository) GetByIDs(ids []int64) ([]*model.User, error) {
	var users []*model.User
	err := mysql.DB.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

