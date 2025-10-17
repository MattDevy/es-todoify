package todo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service provides business logic for Todo operations.
// This is the application service layer in DDD.
type Service struct {
	repo Repository
}

// NewService creates a new Todo service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateTodo creates a new todo item.
func (s *Service) CreateTodo(ctx context.Context, title, description string, labels []string) (*Todo, error) {
	// Validate and create domain object
	todo, err := NewTodo(title, description, labels)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Persist via repository
	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// GetTodo retrieves a todo by ID.
func (s *Service) GetTodo(ctx context.Context, id string) (*Todo, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	// Validate ID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: invalid id format", ErrInvalidInput)
	}

	todo, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// UpdateTodo updates an existing todo.
func (s *Service) UpdateTodo(ctx context.Context, id string, update UpdateTodo) (*Todo, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	// Validate ID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: invalid id format", ErrInvalidInput)
	}

	// Retrieve existing todo
	todo, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates using domain logic
	if err := todo.Update(update); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Persist changes
	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// ChangeStatus changes the status of a todo.
func (s *Service) ChangeStatus(ctx context.Context, id string, newStatus Status) (*Todo, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	// Validate ID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: invalid id format", ErrInvalidInput)
	}

	// Retrieve existing todo
	todo, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply status change using domain logic (validates business rules)
	if err := todo.ChangeStatus(newStatus); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidStatus, err)
	}

	// Persist changes
	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to update todo status: %w", err)
	}

	return todo, nil
}

// DeleteTodo removes a todo by ID.
func (s *Service) DeleteTodo(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	// Validate ID format
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid id format", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

// ListTodos retrieves todos with filtering and pagination.
func (s *Service) ListTodos(ctx context.Context, filter ListFilter) ([]*Todo, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("%w: invalid filter", err)
	}

	// Apply sensible defaults if not provided
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000 // Max limit to prevent resource exhaustion
	}

	// Apply default sorting if not provided
	if filter.SortBy == "" {
		filter.SortBy = SortFieldCreateTime
	}
	if filter.SortOrder == "" {
		filter.SortOrder = SortOrderDesc
	}

	return s.repo.List(ctx, filter)
}

// CountTodos returns the total count of todos matching the filter.
func (s *Service) CountTodos(ctx context.Context, filter ListFilter) (int, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return 0, fmt.Errorf("%w: invalid filter", err)
	}

	return s.repo.Count(ctx, filter)
}
