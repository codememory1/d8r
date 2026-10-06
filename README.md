# d8r

`d8r` is a Go service for managing HTTP downloads. It stores tasks in PostgreSQL, inspects remote resources, chooses a download strategy, and processes downloads with background workers. Clients create and monitor tasks through the HTTP API. The CLI currently has one application command: `server`.

## Features

- Resource inspection and single-request, parallel-range, and streaming downloads
- Concurrent inspection and download workers
- Persistent tasks and cursor-based task listing
- Outbox-backed events and HTTP webhooks, including download progress
- YAML configuration, filesystem storage, and a Docker Hub image

## How it works

1. A client sends `POST /tasks`.
2. An inspection worker checks the resource and selects a download strategy.
3. A download worker saves the resource and updates the task state.
4. Webhooks deliver task events asynchronously.

```mermaid
stateDiagram-v2
    [*] --> pending
    pending --> inspecting
    inspecting --> ready_to_download
    ready_to_download --> downloading
    downloading --> completed
    inspecting --> failed
    downloading --> failed
```

## Run with Docker Compose

The image is available as `codememory/d8r:latest`. Create the following files locally; publishing a Compose file is not required. This example exposes the API on `http://127.0.0.1:7171`.

```bash
mkdir -p config files
```

Create `compose.yaml`:

```yaml
services:
  d8r:
    image: codememory/d8r:latest
    volumes:
      - ./config/main.yaml:/app/config/config.yaml
      - ./files:/app/files
    ports:
      - "7171:7171"
    environment:
      DATABASE_DSN: "postgres://admin:admin@postgres:5432/go?sslmode=disable"
    networks:
      - d8r_test
    depends_on:
      - postgres

  postgres:
    container_name: d8r-test-postgres
    image: postgres:16.4
    restart: unless-stopped
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin
      POSTGRES_DB: go
    ports:
      - "5490:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - d8r_test

networks:
  d8r_test:
    name: d8r_test

volumes:
  postgres_data:
```

Create `config/main.yaml` using the full example in [Configuration](#configuration). Its `postgres.host` is `postgres`, and its `http.port` is `7171`.

Pull the image, start PostgreSQL, and check that it is accepting connections:

```bash
docker compose pull d8r
docker compose up -d postgres
docker compose exec postgres pg_isready -U admin -d go
```

When `pg_isready` reports `accepting connections`, apply the migrations included in the image:

```bash
docker compose run --rm --no-deps \
  --entrypoint /app/bin/migrate \
  d8r up
```

Start the server:

```bash
docker compose up -d d8r
docker compose logs -f d8r
```

The `DATABASE_DSN` variable is read by `/app/bin/migrate`; the server reads the `postgres` section of `config/main.yaml`. Downloads are stored in the local `files/` directory and database data in the `postgres_data` volume. The host port `5490` is for connecting to PostgreSQL from your computer; containers use `postgres:5432`.

To fetch a newer image and recreate the server container:

```bash
docker compose pull d8r
docker compose up -d --force-recreate d8r
```

## Configuration

The server loads the YAML file specified by `--configuration`. The CLI default is `./config/config.yaml`; the image starts with `/app/config/config.yaml`. The file must exist, even though omitted settings receive defaults in code.

The following `config/main.yaml` matches the Compose example and documents **all configuration fields** in the Go configuration. The credentials are intended for local development.

```yaml
postgres:
  # PostgreSQL hostname. Inside Compose, use its service name.
  host: "postgres"
  # PostgreSQL port inside the Compose network.
  port: 5432
  # Database user.
  user: "admin"
  # Database password (local development only).
  password: "admin"
  # Database name.
  database: "go"

http:
  # Address on which the HTTP server listens.
  address: "0.0.0.0"
  # Server port inside the container; match the Compose port mapping.
  port: 7171
  # Maximum time for reading an HTTP request.
  read_timeout: "10s"
  # Maximum time for writing an HTTP response.
  write_timeout: "20s"
  # Maximum idle connection lifetime.
  idle_timeout: "20s"

download:
  # Minimum resource size at which parallel downloading is considered.
  min_parallel_size: "50MiB"
  # Buffer size used while copying response data to storage.
  buffer_size: "256KiB"
  # Target size of an individual download part.
  part_size: "5MiB"
  # Maximum number of download parts processed at once.
  max_parallel_parts: 5
  # Number of byte ranges into which a resource is divided.
  range_parts: 5
  retry:
    # Maximum number of download attempts.
    max_attempts: 3
    # Delay between attempts.
    delay: "10ms"
  # Interval between download progress publications.
  progress_interval: "5s"

workers:
  download:
    # Maximum number of download tasks processed concurrently.
    concurrency: 5
    # Maximum number of tasks claimed in one batch.
    limit: 20
  inspection:
    # Maximum number of inspections processed concurrently.
    concurrency: 5
    # Maximum number of tasks claimed in one batch.
    limit: 20
  outbox_event:
    # Maximum number of outbox events processed concurrently.
    concurrency: 3
    # Maximum number of events claimed in one batch.
    limit: 20
  webhook_delivery:
    # Maximum number of webhook deliveries processed concurrently.
    concurrency: 3
    # Maximum number of deliveries claimed in one batch.
    limit: 20

webhook:
  # Maximum duration of one webhook HTTP request.
  request_timeout: "10s"
  # Maximum number of webhook delivery attempts.
  max_attempts: 5
  # Initial delay before retrying a failed delivery.
  retry_delay: "5s"
```

The built-in defaults differ from this Compose-specific example in four places: `postgres.host` is `localhost`, `postgres.user` is `root`, `postgres.database` is `postgres`, and `http.port` is `8080`. The remaining example values match the defaults in the Go configuration. Use `localhost` as the database host when running the application on your computer, and `postgres` when running it in this Compose network. Durations use values such as `10ms` and `5s`; byte sizes use values such as `256KiB` and `50MiB`.

An older YAML example used `webhook.concurrency`. The Go configuration instead defines `workers.webhook_delivery.concurrency`; `webhook.concurrency` does not configure a worker.

## CLI

```bash
d8r server --configuration ./config/config.yaml
```

| Option | Description | Default |
| --- | --- | --- |
| `-c`, `--configuration` | Path to the YAML configuration file | `./config/config.yaml` |
| `-v`, `--version` | Print the application version | — |
| `-h`, `--help` | Show help | — |

The Docker image runs `server --configuration /app/config/config.yaml` by default. There is no `d8r download` command; downloads are created through the HTTP API.

To run the application from source, use the Go version declared in `go.mod`, prepare PostgreSQL, apply the database migrations, and run:

```bash
go run ./cmd/d8r server --configuration ./config/config.yaml
```

## HTTP API

All examples below use the Compose URL `http://127.0.0.1:7171`. Creating a task schedules background processing; the HTTP request does not wait for the download to finish.

### Create a task

`POST /tasks`

```bash
curl -X POST 'http://127.0.0.1:7171/tasks' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://httpbin.org/range/1024",
    "headers": {},
    "filename": null,
    "priority": -256
  }'
```

`headers` supplies request headers for the remote resource. `filename` may be `null`; `priority` controls task processing order.

### Get a task

`GET /tasks/{id}`

```bash
curl 'http://127.0.0.1:7171/tasks/TASK_ID'
```

Replace `TASK_ID` with the ID returned by task creation.

### List tasks

`GET /tasks?limit=...&cursor=...`

```bash
curl 'http://127.0.0.1:7171/tasks?limit=20'
```

For the next page, pass the cursor returned by the previous page as the `cursor` query parameter.

### Register a webhook

`POST /webhooks`

```bash
curl -X POST 'http://127.0.0.1:7171/webhooks' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://example.com/webhooks/d8r",
    "headers": {},
    "events": [
      "task.created",
      "task.inspection.started",
      "task.inspection.completed",
      "task.inspection.failed",
      "task.download.started",
      "task.download.completed",
      "task.download.failed",
      "task.download.progress"
    ]
  }'
```

Set `url` to an endpoint reachable from the d8r container. An address such as `127.0.0.1` in that URL refers to the container itself.

| Event | When it occurs |
| --- | --- |
| `task.created` | A task is created |
| `task.inspection.started` | Resource inspection starts |
| `task.inspection.completed` | Inspection succeeds |
| `task.inspection.failed` | Inspection fails |
| `task.download.started` | Downloading starts |
| `task.download.completed` | Downloading succeeds |
| `task.download.failed` | Downloading fails |
| `task.download.progress` | New download progress is reported |

## Project layout

```text
bin/migrate               Migration wrapper
cmd/d8r/                  CLI entry point
config/                   YAML configuration
internal/application/     Commands, queries, and application interfaces
internal/bootstrap/       Dependency construction
internal/cli/             Cobra command wiring
internal/domain/          Entities, value objects, and domain events
internal/infrastructure/  PostgreSQL, HTTP download, storage, and delivery adapters
migrations/               Database migrations
```

## Development

```bash
go test ./...
go vet ./...
```

The project is under active development.
