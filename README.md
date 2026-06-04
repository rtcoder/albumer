# Albumer

A private catalog for a collection of CDs, vinyl records, cassettes, and books. The application is distributed as a single Go binary with server-side rendered views, a small amount of embedded JavaScript for filtering, and Neo4j running as a separate service.

The frontend is rendered by the Go backend, but it is kept in regular HTML, CSS, and JavaScript files:

- `internal/web/templates/*.gohtml`
- `internal/web/assets/app.css`
- `internal/web/assets/app.js`

Templates use the `.gohtml` extension so GoLand can recognize Go HTML templates and provide view data hints. The files are embedded in the binary using `go:embed`, so deployment still requires only the Go application.

## Requirements

- Go 1.22+
- Neo4j 5+

## Configuration

The application reads configuration from environment variables:

```sh
ADDR=:4000
NEO4J_URI=neo4j://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=neo4j
APP_SECRET=change-me
UPLOAD_DIR=./uploads
ADMIN_PASSWORD_HASH='$2a$10$...'
```

At startup, the application also loads a local `.env` file when one exists. Environment variables that are already set take precedence over values from `.env`.

If `ADMIN_PASSWORD_HASH` is empty, the default development password is `admin`. Use a bcrypt hash on a server. The interactive installer asks for a plain-text password and creates the hash automatically.

## Local Development

```sh
go mod tidy
go run ./cmd/albumer
```

Open `http://localhost:4000` after the application starts.

## Data Model

- `Person` represents a music artist or book author.
- `Group` represents a band and has `HAS_MEMBER` relationships to people.
- `MusicItem` represents a physical music item. Its required `medium` field accepts `CD`, `VINYL`, or `CASSETTE`.
- `Book` represents a book.
- Music items and books use one of the statuses: `owned`, `ordered`, or `wanted`.

The same album on CD and vinyl is represented by two separate `MusicItem` nodes. Releases of the same album on different media can be connected using the `SAME_RELEASE` relationship. The music form allows relationships to be selected manually and suggests existing items after a title is entered.

## systemd Deployment

On a Linux server with `systemd` and Neo4j installed, run the interactive installer:

```sh
./install.sh
```

The installer checks for Neo4j, asks for configuration values, verifies that the selected port is available, hashes the administrator password with bcrypt, and installs the application as the `albumer` service. Service configuration is stored in `/opt/albumer/.env`.

If no administrator password is provided, the application keeps the default password `admin`.

## Tests

```sh
go test ./...
```

## Importing Legacy JSON

The old `json/data.json` file can be loaded into Neo4j as a collection of CDs:

```sh
go run ./cmd/import-json -file json/data.json
```

The importer creates people, bands (`Band` nodes in Neo4j), band member relationships, and `MusicItem` nodes with `medium=CD`. Database records receive new UUIDs; old JSON IDs are used only temporarily while relationships are being built.

If the database was populated before friendly search was added, run the search backfill once:

```sh
go run ./cmd/backfill-search
```

This fills the `searchName` and `searchTitle` fields, allowing a query such as `ostr` to find `O.S.T.R`.

Existing records with old random IDs can be migrated to UUIDs:

```sh
go run ./cmd/migrate-ids
```

The migration stores previous IDs in `legacyId`, so relationships are preserved.
