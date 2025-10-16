# Todoify

## What is this repo?

This repo contains a Golang todo application, `todoify`, that uses the [go-elasticsearch](https://github.com/elastic/go-elasticsearch) SDK to create, update, list and delete TODOs with Elasticsearch as the backend.

## Getting started

### Initial setup

We need to install/start [start-local](https://github.com/elastic/start-local)

```bash
make es-setup
```

If already installed, but not running

```bash
make es-start
```

## Core features

- [ ] Create TODO
- [ ] Update TODO
- [ ] List TODOs (with filters and pagination)
- [ ] Delete TODO
- [ ] Search TODOs
- [ ] Bulk TODO upload
  - [ ] jsonl
  - [ ] csv
- [ ] TODO stats

## Extended features

- [ ] API
- [ ] UI work
- [ ] Tracing, metrics and log ingest to Elastic
- [ ] User management

## Implementation

This section covers the general plan for implementation.

1. Use [start-local](https://github.com/elastic/start-local) to give us Elasticsearch and Kibana locally
2. Use viper and cobra to give us a cli and configuration
3. Use a model & repository pattern to abstract our business logic from the backend (Elasticsearch)
4. Implement the Elasticsearch backend using [go-elasticsearch](https://github.com/elastic/go-elasticsearch)
