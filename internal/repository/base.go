package repository

import (
	"context"
	"time"
)

// HealthStatus represents the overall health state of a repository backend.
type HealthStatus string

const (
	// HealthStatusHealthy indicates the backend is fully operational.
	HealthStatusHealthy HealthStatus = "healthy"

	// HealthStatusDegraded indicates the backend is operational but with reduced capacity or performance.
	HealthStatusDegraded HealthStatus = "degraded"

	// HealthStatusUnhealthy indicates the backend is not operational or experiencing critical issues.
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthInfo contains health check information from a repository backend.
// This structure is designed to be backend-agnostic while allowing backend-specific details.
type HealthInfo struct {
	// Status is the overall health status of the backend.
	Status HealthStatus `json:"status"`

	// Available indicates whether the backend is reachable and responding.
	Available bool `json:"available"`

	// ResponseTime is the duration taken to perform the health check.
	ResponseTime time.Duration `json:"responseTime"`

	// NodeCount is the number of nodes in a clustered backend (nil if not applicable).
	// Examples: ES cluster nodes, MongoDB replica set members, PostgreSQL streaming replicas.
	NodeCount *int `json:"nodeCount,omitempty"`

	// ActiveConnections is the number of active connections to the backend (nil if not applicable).
	// Examples: ES HTTP connections, PostgreSQL active connections, MongoDB connections.
	ActiveConnections *int `json:"activeConnections,omitempty"`

	// Version is the version string of the backend.
	Version string `json:"version,omitempty"`

	// Details contains backend-specific health information.
	// This allows each implementation to provide additional metrics without breaking the interface.
	//
	// Examples:
	//   Elasticsearch: cluster_name, active_shards, relocating_shards, unassigned_shards
	//   MongoDB: replica_set_name, oplog_window, replication_lag
	//   PostgreSQL: database_size, max_connections, cache_hit_ratio
	Details map[string]interface{} `json:"details,omitempty"`
}

// Base defines common operations that all repositories should support.
// This interface should be embedded in domain-specific repository interfaces.
type Base interface {
	// Health performs a health check on the repository backend and returns health information.
	// Returns an error if the health check fails critically.
	Health(ctx context.Context) (*HealthInfo, error)
}
