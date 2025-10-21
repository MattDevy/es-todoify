package service

import (
	"context"

	"github.com/MattDevy/es-todoify/internal/repository"
)

// Base defines common operations that all services should support.
// This interface should be embedded in domain-specific service interfaces.
type Base interface {
	// Health performs a health check on the underlying repository and returns health information.
	// This delegates to the repository's Health method.
	Health(ctx context.Context) (*repository.HealthInfo, error)
}
