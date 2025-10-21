package todo

import (
	"context"
	"time"

	"github.com/MattDevy/es-todoify/internal/repository"
)

// Repository defines the interface for Todo persistence operations.
// This is a repository pattern implementation that abstracts storage details.
type Repository interface {
	// Embed base repository interface for common operations (e.g., Health)
	repository.Base

	// Create persists a new Todo.
	Create(ctx context.Context, todo *Todo) error

	// Get retrieves a Todo by ID.
	// Returns ErrNotFound if the todo doesn't exist.
	Get(ctx context.Context, id string) (*Todo, error)

	// Update updates an existing Todo.
	// Returns ErrNotFound if the todo doesn't exist.
	Update(ctx context.Context, todo *Todo) error

	// Delete removes a Todo by ID.
	// Returns ErrNotFound if the todo doesn't exist.
	Delete(ctx context.Context, id string) error

	// List retrieves todos with optional filtering and pagination.
	List(ctx context.Context, filter ListFilter) ([]*Todo, error)

	// Count returns the total number of todos matching the filter.
	Count(ctx context.Context, filter ListFilter) (int, error)
}

// ListFilter defines filtering and pagination options for listing todos.
type ListFilter struct {
	// Status filters by todo status (empty = all)
	Status Status

	// Labels filters todos that have all specified labels
	Labels []string

	// SearchQuery performs full-text search on title and description
	SearchQuery string

	// FromDate filters todos created on or after this date
	FromDate *time.Time

	// ToDate filters todos created on or before this date
	ToDate *time.Time

	// Limit is the maximum number of results to return
	Limit int

	// Offset is the number of results to skip (for pagination)
	Offset int

	// SortBy specifies the field to sort by
	SortBy SortField

	// SortOrder specifies ascending or descending order
	SortOrder SortOrder
}

// Validate checks if the ListFilter has valid values.
func (f ListFilter) Validate() error {
	// Validate status if provided
	if f.Status != "" && !f.Status.IsValid() {
		return ErrInvalidInput
	}

	// Validate sort field if provided
	if f.SortBy != "" && !f.SortBy.IsValid() {
		return ErrInvalidInput
	}

	// Validate sort order if provided
	if f.SortOrder != "" && !f.SortOrder.IsValid() {
		return ErrInvalidInput
	}

	// Validate limit (must be positive if provided)
	if f.Limit < 0 {
		return ErrInvalidInput
	}

	// Validate offset (must be non-negative)
	if f.Offset < 0 {
		return ErrInvalidInput
	}

	// Validate date range (FromDate must be before ToDate)
	if f.FromDate != nil && f.ToDate != nil && f.FromDate.After(*f.ToDate) {
		return ErrInvalidInput
	}

	return nil
}

// DefaultListFilter returns a ListFilter with sensible defaults.
func DefaultListFilter() ListFilter {
	return ListFilter{
		Limit:     50,
		Offset:    0,
		SortBy:    SortFieldCreateTime,
		SortOrder: SortOrderDesc,
	}
}
