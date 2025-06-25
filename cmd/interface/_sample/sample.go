package sample

import "fmt"

// UserService provides user management functionality
type UserService struct {
	users map[string]string
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]string),
	}
}

// GetUser retrieves a user by ID
func (u *UserService) GetUser(id string) (string, error) {
	user, exists := u.users[id]
	if !exists {
		return "", fmt.Errorf("user not found: %s", id)
	}
	return user, nil
}

// CreateUser creates a new user
func (u *UserService) CreateUser(id, name string) error {
	if _, exists := u.users[id]; exists {
		return fmt.Errorf("user already exists: %s", id)
	}
	u.users[id] = name
	return nil
}

// UpdateUser updates an existing user
func (u *UserService) UpdateUser(id, name string) error {
	if _, exists := u.users[id]; !exists {
		return fmt.Errorf("user not found: %s", id)
	}
	u.users[id] = name
	return nil
}

// DeleteUser deletes a user
func (u *UserService) DeleteUser(id string) error {
	if _, exists := u.users[id]; !exists {
		return fmt.Errorf("user not found: %s", id)
	}
	delete(u.users, id)
	return nil
}

// ListUsers returns all users
func (u *UserService) ListUsers() map[string]string {
	result := make(map[string]string)
	for k, v := range u.users {
		result[k] = v
	}
	return result
}

// privateMethod is not exported and should not be included in interface
func (u *UserService) privateMethod() {
	// internal implementation
}