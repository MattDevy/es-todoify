# Todo Index Mapping

This directory contains the Elasticsearch mapping definition for the `todos` index.

## Mapping Overview

The `todo.json` file defines the index settings and field mappings for todo documents stored in Elasticsearch.

## Field Definitions

### title (required)

- **Type**: `text` with `keyword` multi-field
- **Purpose**: The title or summary of the todo item
- **Features**:
  - Full-text search on the analyzed `title` field
  - Exact matching, sorting, and aggregations on `title.keyword`
  - `ignore_above: 256` on keyword field for very long titles

### description (optional)

- **Type**: `text`
- **Purpose**: Extended description of the todo item
- **Features**:
  - Full-text search with standard analyzer
  - Tokenized and analyzed for flexible searching

### labels (optional)

- **Type**: `keyword` (array)
- **Purpose**: Categorization and grouping of todos
- **Features**:
  - Exact matching for filtering
  - Efficient aggregations (terms aggregation)
  - Can contain multiple values per document
- **Example**: `["work", "urgent", "backend"]`

### status (required)

- **Type**: `keyword`
- **Purpose**: Current state of the todo item
- **Features**:
  - Exact matching for filtering
  - Efficient aggregations
  - Terms query for filtering by status

#### Valid Status Values

The following status values are recommended:

- **`pending`** - Todo has been created but not started
- **`in_progress`** - Work has begun on the todo
- **`completed`** - Todo has been finished
- **`cancelled`** - Todo was cancelled and won't be completed
- **`blocked`** - Todo is blocked by dependencies or external factors

### createTime (required)

- **Type**: `date`
- **Purpose**: Timestamp when the todo was created
- **Format**: `strict_date_optional_time||epoch_millis`
- **Features**:
  - Range queries for filtering by date ranges
  - Date histogram aggregations
  - Sorting by creation date
- **Accepted formats**:
  - ISO 8601: `2024-01-15T10:30:00Z`
  - Date only: `2024-01-15`
  - Epoch milliseconds: `1705318200000`

### updateTime (required)

- **Type**: `date`
- **Purpose**: Timestamp of the last update to the todo
- **Format**: `strict_date_optional_time||epoch_millis`
- **Features**: Same as createTime
- **Note**: Should be updated whenever any field in the todo is modified

## Index Settings

- **Shards**: 1 (suitable for small to medium datasets)
- **Replicas**: 1 (provides basic fault tolerance)
- **Default Analyzer**: Standard analyzer (good for general text in multiple languages)

## Usage Examples

### Creating the Index

```bash
curl -X PUT "localhost:9200/todos" -H 'Content-Type: application/json' -d @mappings/todo.json
```

### Sample Document

```json
{
  "title": "Implement user authentication",
  "description": "Add OAuth2 authentication with support for Google and GitHub providers. Include refresh token rotation and secure session management.",
  "labels": ["backend", "security", "urgent"],
  "status": "in_progress",
  "createTime": "2024-01-15T10:30:00Z",
  "updateTime": "2024-01-16T14:22:00Z"
}
```

### Query Examples

#### Full-text search on title and description

```json
{
  "query": {
    "multi_match": {
      "query": "authentication security",
      "fields": ["title", "description"]
    }
  }
}
```

#### Filter by status and labels

```json
{
  "query": {
    "bool": {
      "must": [
        { "term": { "status": "in_progress" } },
        { "terms": { "labels": ["backend", "urgent"] } }
      ]
    }
  }
}
```

#### Range query on createTime

```json
{
  "query": {
    "range": {
      "createTime": {
        "gte": "2024-01-01",
        "lte": "2024-01-31"
      }
    }
  }
}
```

#### Aggregations by status

```json
{
  "size": 0,
  "aggs": {
    "todos_by_status": {
      "terms": {
        "field": "status"
      }
    }
  }
}
```

#### Aggregations by labels

```json
{
  "size": 0,
  "aggs": {
    "popular_labels": {
      "terms": {
        "field": "labels",
        "size": 10
      }
    }
  }
}
```

#### Date histogram aggregation

```json
{
  "size": 0,
  "aggs": {
    "todos_over_time": {
      "date_histogram": {
        "field": "createTime",
        "calendar_interval": "day"
      }
    }
  }
}
```

## Best Practices

1. **Always set createTime and updateTime**: These are crucial for auditing and time-based queries
2. **Use appropriate status values**: Stick to the predefined status values for consistency
3. **Keep labels consistent**: Use lowercase, hyphen-separated labels (e.g., `high-priority` not `High Priority`)
4. **Title length**: Keep titles concise (under 256 characters) for optimal keyword field performance
5. **Update updateTime**: Always update this field when modifying any todo field

## Migration Notes

If you need to update this mapping in the future:

1. Most field type changes require reindexing
2. You can add new fields without reindexing
3. Use the Reindex API to migrate data to a new index with updated mappings
4. Consider using index aliases for zero-downtime migrations
