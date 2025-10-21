package sdk

import (
	"context"
	"log/slog"

	"github.com/MattDevy/es-todoify/internal/todo"
)

type serviceContextKey struct{}
type repoContextKey struct{}
type loggerContextKey struct{}

// WithService adds the service to the context, this is useful to pass the service to the sub-commands.
func WithService(ctx context.Context, service *todo.Service) context.Context {
	return context.WithValue(ctx, serviceContextKey{}, service)
}

// WithRepo adds the repository to the context, this is useful to pass the repository to the sub-commands.
func WithRepo(ctx context.Context, repo todo.Repository) context.Context {
	return context.WithValue(ctx, repoContextKey{}, repo)
}

// WithLogger adds the logger to the context, this is useful to pass the logger to the sub-commands.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// GetService gets the service from the context, this is useful to get the service from the context in sub-commands.
func GetService(ctx context.Context) *todo.Service {
	return ctx.Value(serviceContextKey{}).(*todo.Service)
}

// GetRepo gets the repository from the context, this is useful to get the repository from the context in sub-commands.
func GetRepo(ctx context.Context) todo.Repository {
	return ctx.Value(repoContextKey{}).(todo.Repository)
}

// GetLogger gets the logger from the context, this is useful to get the logger from the context in sub-commands.
func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerContextKey{}).(*slog.Logger)
}
