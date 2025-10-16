package todo

// Status represents the state of a Todo item.
// This is a value object in DDD terminology.
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
	StatusBlocked    Status = "blocked"
)

// IsValid checks if the status is one of the defined values.
func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled, StatusBlocked:
		return true
	}
	return false
}

// String returns the string representation of the Status.
func (s Status) String() string {
	return string(s)
}

// AllStatuses returns all valid status values.
func AllStatuses() []Status {
	return []Status{
		StatusPending,
		StatusInProgress,
		StatusCompleted,
		StatusCancelled,
		StatusBlocked,
	}
}
