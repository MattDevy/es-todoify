package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/MattDevy/es-todoify/internal/todo"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"

	_ "embed"
)

//go:embed indicies/todo.json
var todoIndex []byte

// Repository is the implementation of the Repository interface for Elasticsearch.
type Repository struct {
	client    *elasticsearch.TypedClient
	indexName string
}

// NewRepository creates a new Repository.
func NewRepository(client *elasticsearch.TypedClient, indexName string) *Repository {
	return &Repository{
		client:    client,
		indexName: indexName,
	}
}

// CreateIndices creates the indices for the repository.
func (r *Repository) CreateIndices(ctx context.Context) error {
	cr := &create.Request{}
	if err := json.NewDecoder(bytes.NewReader(todoIndex)).Decode(cr); err != nil {
		return fmt.Errorf("failed to decode todo index: %w", err)
	}
	res, err := r.client.Indices.Create(r.indexName).Request(cr).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	if !res.Acknowledged {
		return fmt.Errorf("failed to create index: %s not acknowledged", res.Index)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, t *todo.Todo) error {
	_, err := r.client.Create(r.indexName, t.ID.String()).Document(t).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*todo.Todo, error) {
	res, err := r.client.GetSource(r.indexName, id).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	var t todo.Todo
	if err := json.NewDecoder(bytes.NewReader(res)).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode todo: %w", err)
	}

	return &t, nil
}
