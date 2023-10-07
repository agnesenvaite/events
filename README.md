# Events

## Setup

1. Install [docker](https://docs.docker.com/engine/installation/)
   and [docker-compose](https://docs.docker.com/compose/).
   Complete [post-install](https://docs.docker.com/engine/installation/linux/linux-postinstall/)
   instructions.
2. Run ```docker compose up -d``` to start development server.
3. API should now be accessible through http://localhost:8080/v1/status/ping
4. Run initial migrations script `make migrate-init`


## Migrations

Run migrations:

```make migrate```

Create migration:

```make new-migration```

## Tests

1. Generate mocks `make mocks`
2. Run tests `make test`

## Manual testing api

1. Docs are accessible through http://localhost:8080/v1/docs/
