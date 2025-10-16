package todo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Todo represents a TODO item in our domain.
// This is our aggregate root in DDD terminology.
type Todo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Labels      []string  `json:"labels,omitempty"`
	Status      Status    `json:"status"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// NewTodo creates a new Todo with validation.
// Factory function ensures all required fields are set and valid.
func NewTodo(title string, description string, labels []string) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	now := time.Now()
	return &Todo{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Labels:      labels,
		Status:      StatusPending,
		CreateTime:  now,
		UpdateTime:  now,
	}, nil
}

type UpdateTodo struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1024"`
	Labels      []string `json:"labels,omitempty" validate:"omitempty,min=1,max=10"`
}

func (u UpdateTodo) Validate() error {
	return validate.Struct(u)
}

// Update updates the mutable fields of a Todo.
func (t *Todo) Update(update UpdateTodo) error {
	if err := update.Validate(); err != nil {
		return err
	}

	if update.Title != nil {
		t.Title = *update.Title
	}

	if update.Description != nil {
		t.Description = *update.Description
	}

	if update.Labels != nil {
		t.Labels = update.Labels
	}

	t.UpdateTime = time.Now()

	return nil
}

// ChangeStatus transitions the Todo to a new status.
func (t *Todo) ChangeStatus(newStatus Status) error {
	if !newStatus.IsValid() {
		return errors.New("invalid status")
	}

	// Business rule: completed todos cannot be set to blocked
	if t.Status == StatusCompleted && newStatus == StatusBlocked {
		return errors.New("cannot block a completed todo")
	}

	t.Status = newStatus
	t.UpdateTime = time.Now()

	return nil
}

// IsCompleted returns true if the todo is in a completed state.
func (t *Todo) IsCompleted() bool {
	return t.Status == StatusCompleted
}
