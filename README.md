# Todoify

A command-line todo application powered by Elasticsearch, built with Go.

## What is this?

Todoify is a CLI tool for managing todos with Elasticsearch as the backend. It leverages Elasticsearch's powerful full-text search, aggregations, and filtering capabilities to provide a robust todo management experience.

## Features

### Currently Implemented ✅

- **Flexible Configuration**: Multiple ways to configure Elasticsearch connection
  - Command-line flags
  - Environment variables
  - Configuration file (`~/.todoify.yaml`)
- **Elasticsearch Integration**:
  - Typed client using [go-elasticsearch v9](https://github.com/elastic/go-elasticsearch)
  - Support for multiple Elasticsearch addresses (cluster support)
  - Authentication via username/password or API key
  - Connection validation on startup
- **Index Mapping**: Production-ready mapping for todos with:
  - Full-text search on title and description
  - Filtering and aggregations on labels and status
  - Date range queries on createTime and updateTime

### Roadmap 🚧

- [ ] Create TODO
- [ ] Update TODO
- [ ] List TODOs (with filters and pagination)
- [ ] Delete TODO
- [ ] Search TODOs
- [ ] Bulk TODO upload
  - [ ] JSONL
  - [ ] CSV
- [ ] TODO stats and aggregations

### Extended Features (Future)

- [ ] REST API
- [ ] Web UI
- [ ] OpenTelemetry tracing, metrics, and log ingestion
- [ ] Multi-user support

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- Docker (for running Elasticsearch locally)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/MattDevy/es-todoify.git
cd es-todoify
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o todoify
```

### Setting Up Elasticsearch

We use [start-local](https://github.com/elastic/start-local) for local Elasticsearch development.

#### Initial Setup

Install and start Elasticsearch with OpenTelemetry collector:

```bash
make es-setup
```

#### Starting/Stopping Elasticsearch

If already installed but not running:

```bash
make es-start
```

To stop Elasticsearch:

```bash
make es-stop
```

To restart:

```bash
make es-restart
```

#### Get Elasticsearch Credentials

```bash
make es-creds
```

This will display:

- Username (default: `elastic`)
- Password
- API Key

### Creating the Todo Index

Once Elasticsearch is running, create the todo index with the mapping:

```bash
curl -X PUT "localhost:9200/todos" \
  -H 'Content-Type: application/json' \
  -u elastic:YOUR_PASSWORD \
  -d @mappings/todo.json
```

Or using the API key:

```bash
curl -X PUT "localhost:9200/todos" \
  -H 'Content-Type: application/json' \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -d @mappings/todo.json
```

## Configuration

Todoify uses a hierarchical configuration system (highest to lowest priority):

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefixed with `TODOIFY_`)
3. **Configuration file** (`~/.todoify.yaml`)
4. **Default values** (lowest priority)

### Configuration Options

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--es-addrs` | `TODOIFY_ES_ADDRS` | `http://localhost:9200` | Elasticsearch addresses (comma-separated) |
| `--es-username` | `TODOIFY_ES_USERNAME` | - | Elasticsearch username |
| `--es-password` | `TODOIFY_ES_PASSWORD` | - | Elasticsearch password |
| `--es-api-key` | `TODOIFY_ES_API_KEY` | - | Elasticsearch API key |
| `--es-index` | `TODOIFY_ES_INDEX` | `todos` | Elasticsearch index name |
| `--config` | - | `~/.todoify.yaml` | Config file path |

### Configuration Examples

#### Using Command-Line Flags

```bash
todoify create \
  --es-addrs=http://localhost:9200 \
  --es-username=elastic \
  --es-password=changeme \
  --es-index=todos
```

#### Using Environment Variables

```bash
export TODOIFY_ES_ADDRS=http://localhost:9200
export TODOIFY_ES_USERNAME=elastic
export TODOIFY_ES_PASSWORD=changeme
export TODOIFY_ES_INDEX=todos

todoify create
```

#### Using Configuration File

Create `~/.todoify.yaml`:

```yaml
es-addrs:
  - http://localhost:9200
  - http://localhost:9201  # Optional: for cluster support
es-username: elastic
es-password: changeme
es-index: todos
```

Or using API key authentication:

```yaml
es-addrs:
  - http://localhost:9200
es-api-key: your-api-key-here
es-index: todos
```

Then run:

```bash
todoify create
```

### Authentication

Todoify supports two authentication methods (mutually exclusive):

1. **Username/Password**: Both must be provided together

   ```bash
   todoify create --es-username=elastic --es-password=changeme
   ```

2. **API Key**: More secure for production use

   ```bash
   todoify create --es-api-key=your-api-key
   ```

**Note**: You cannot use both authentication methods simultaneously.

## Usage

### General Help

```bash
todoify --help
```

### Command-Specific Help

```bash
todoify create --help
```

## Todo Data Model

Each todo document contains the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | text/keyword | Yes | Todo summary (searchable and sortable) |
| `description` | text | No | Extended description (searchable) |
| `labels` | keyword[] | No | Array of labels for categorization |
| `status` | keyword | Yes | Current status (`pending`, `in_progress`, `completed`, `cancelled`, `blocked`) |
| `createTime` | date | Yes | Creation timestamp |
| `updateTime` | date | Yes | Last update timestamp |

See [`mappings/README.md`](mappings/README.md) for detailed mapping documentation and query examples.

## Project Structure

```
es-todoify/
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command with ES client initialization
│   └── create.go          # Create command (in progress)
├── docs/                   # Documentation
│   └── planning/          # Implementation planning docs
├── elastic-start-local/    # Local Elasticsearch setup
├── internal/              # Internal packages (future)
├── mappings/              # Elasticsearch index mappings
│   ├── todo.json          # Todo index mapping
│   └── README.md          # Mapping documentation
├── go.mod                 # Go module definition
├── main.go               # Application entry point
├── Makefile              # Development tasks
└── README.md             # This file
```

## Development

### Building

```bash
go build -o todoify
```

### Dependencies

Key dependencies:

- [go-elasticsearch/v9](https://github.com/elastic/go-elasticsearch) - Official Elasticsearch Go client (typed API)
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management

### Architecture

Todoify follows a clean architecture approach:

1. **CLI Layer** (`cmd/`): Cobra commands handle user interaction
2. **Configuration**: Viper manages multi-source configuration
3. **Client**: go-elasticsearch typed client initialized once at startup
4. **Business Logic** (planned): Model & repository pattern for data operations
5. **Storage**: Elasticsearch backend

For detailed implementation plans, see [`docs/planning/00_setup.md`](docs/planning/00_setup.md).

## Troubleshooting

### Connection Issues

If you get authentication errors:

1. Check Elasticsearch is running:

   ```bash
   curl http://localhost:9200
   ```

2. Verify credentials:

   ```bash
   make es-creds
   ```

3. Test connection with curl:

   ```bash
   curl -u elastic:PASSWORD http://localhost:9200
   ```

### Index Issues

If the todo index doesn't exist:

```bash
# Check if index exists
curl http://localhost:9200/todos -u elastic:PASSWORD

# Create the index
curl -X PUT "localhost:9200/todos" \
  -H 'Content-Type: application/json' \
  -u elastic:PASSWORD \
  -d @mappings/todo.json
```

### Configuration Precedence

Remember the configuration hierarchy:

- Flags override environment variables
- Environment variables override config file
- Config file overrides defaults

Use `--help` to see current flag values and defaults.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

See [LICENSE](LICENSE) for details.

## Resources

- [Elasticsearch Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [go-elasticsearch Documentation](https://github.com/elastic/go-elasticsearch)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [start-local](https://github.com/elastic/start-local)
