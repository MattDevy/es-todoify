package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MattDevy/es-todoify/internal/todo"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"

	_ "embed"
)

//go:embed indices/todo.json
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
	res, err := r.client.Get(r.indexName, id).Do(ctx)
	if err != nil {
		// Check for 404 error
		if errors.Is(err, &types.ElasticsearchError{}) || res == nil {
			return nil, todo.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if !res.Found {
		return nil, todo.ErrNotFound
	}

	var t todo.Todo
	if err := json.Unmarshal(res.Source_, &t); err != nil {
		return nil, fmt.Errorf("failed to decode todo: %w", err)
	}

	return &t, nil
}

func (r *Repository) Update(ctx context.Context, t *todo.Todo) error {
	_, err := r.client.Index(r.indexName).
		Id(t.ID.String()).
		Document(t).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	res, err := r.client.Delete(r.indexName, id).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	// Check if document was found and deleted
	if res.Result.Name == "not_found" {
		return todo.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context, filter todo.ListFilter) ([]*todo.Todo, error) {
	// Build query
	query := buildQuery(filter)

	// Build sort
	sortOptions := buildSort(filter)

	// Execute search
	req := &search.Request{
		Query: query,
		Size:  &filter.Limit,
		From:  &filter.Offset,
	}

	if len(sortOptions) > 0 {
		req.Sort = sortOptions
	}

	res, err := r.client.Search().
		Index(r.indexName).
		Request(req).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to search todos: %w", err)
	}

	// Parse results
	var todos []*todo.Todo
	for _, hit := range res.Hits.Hits {
		var t todo.Todo
		if err := json.Unmarshal(hit.Source_, &t); err != nil {
			return nil, fmt.Errorf("failed to parse todo document: %w", err)
		}
		todos = append(todos, &t)
	}

	return todos, nil
}

func (r *Repository) Count(ctx context.Context, filter todo.ListFilter) (int, error) {
	query := buildQuery(filter)

	res, err := r.client.Count().
		Index(r.indexName).
		Query(query).
		Do(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return int(res.Count), nil
}

// buildQuery constructs an Elasticsearch query from a ListFilter.
func buildQuery(filter todo.ListFilter) *types.Query {
	var must []types.Query

	// Status filter
	if filter.Status != "" {
		must = append(must, types.Query{
			Term: map[string]types.TermQuery{
				"status": {Value: filter.Status.String()},
			},
		})
	}

	// Labels filter (must have all specified labels)
	for _, label := range filter.Labels {
		must = append(must, types.Query{
			Term: map[string]types.TermQuery{
				"labels": {Value: label},
			},
		})
	}

	// Full-text search
	if filter.SearchQuery != "" {
		must = append(must, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  filter.SearchQuery,
				Fields: []string{"title^2", "description"}, // Boost title matches
			},
		})
	}

	// Date range filter
	if filter.FromDate != nil || filter.ToDate != nil {
		dateRange := types.DateRangeQuery{}

		if filter.FromDate != nil {
			gte := filter.FromDate.Format("2006-01-02T15:04:05Z07:00")
			dateRange.Gte = &gte
		}
		if filter.ToDate != nil {
			lte := filter.ToDate.Format("2006-01-02T15:04:05Z07:00")
			dateRange.Lte = &lte
		}

		must = append(must, types.Query{
			Range: map[string]types.RangeQuery{
				"createTime": dateRange,
			},
		})
	}

	// If no filters, match all
	if len(must) == 0 {
		return &types.Query{
			MatchAll: &types.MatchAllQuery{},
		}
	}

	// Combine all must clauses
	return &types.Query{
		Bool: &types.BoolQuery{
			Must: must,
		},
	}
}

// buildSort constructs sort parameters from a ListFilter.
func buildSort(filter todo.ListFilter) []types.SortCombinations {
	if filter.SortBy == "" {
		return nil
	}

	order := sortorder.Desc
	if filter.SortOrder == todo.SortOrderAsc {
		order = sortorder.Asc
	}

	// Map domain sort fields to Elasticsearch fields
	field := string(filter.SortBy)

	// For text fields with .keyword subfield, use it for sorting
	if field == "title" {
		field = "title.keyword"
	}

	return []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				field: {Order: &order},
			},
		},
	}
}
