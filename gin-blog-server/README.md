## Refactor

- Use `slog` as the logging library. For Go 1.21+, `slog` is in the standard library. For earlier versions, import `exp/slog` (TODO: log persistence).
- Adopt the official Golang project layout.
- Generate API docs with Swagger (TODO: enrich comments for each API; consider alternatives to go-swagger).
- Replace global variables with dependency injection; initialize in `main.go`.
- Restructure the `utils` package to simplify the architecture.
- Use SQLite for unit tests in development; support MySQL in production.
- Remove Casbin and implement RBAC-based access control.
- Previous versions used panic for error codes; now use a Gin middleware to catch panics globally.
- TODO: Standardize error handling via `errors`.

Run the new version (run/init data): all related files are under `cmd/`.

1. MySQL: set `DbType = "mysql"` in `config.yml`, update MySQL connection info, and import `assets/gvb.sql` (you can also initialize data manually similar to SQLite flow).

2. SQLite: set `DbType = "sqlite"` in `config.yml`, then initialize data:
- `create_superadmin.sh` creates a super admin (all permissions), username `superadmin`, password `superadmin`.
- `generate_data.sh` initializes three default roles (admin, user, guest) + three default users (admin, user, guest; all passwords `123456`), initializes config, pages, and resource data (TODO).

---

The following are notes from older versions and can be skipped.

## Deployment

Use docker-compose for one-command deployment.

`./config/config.docker.toml` is the config file used in deployment; some values are overridden by docker-compose environment variables.

See `deploy/start/docker-compose.yml` for details.

Config precedence: environment variables > values in `config.docker.toml`. When using docker-compose, modify the `environment` block.

## Development Guidelines

Model layer returns `error`; Service layer returns a `code`; Controller uses `GetMsg(code)` to return messages to frontend.

- Error codes are maintained in `global/errmsg`.

For JSON, use snake_case; in Go code, use camelCase.

## Database

For MySQL boolean fields, use `tinyint`. In Go structs, define such fields as pointers so you can distinguish zero values via `nil`.

## Tests

Unit tests are important, especially for refactoring and iterative changes. They reduce manual API calls and ensure correctness.

## Gin

### Validator and zero values

Gin uses `validator` for parameter validation. If a field is tagged `required`, it must not receive the type's zero value.
- Strings: cannot be empty.
- Int: cannot be 0.
- Bool: cannot be false.

Sometimes a field is required yet 0 is a valid value (e.g., `sex` 0=female, 1=male). Use a pointer to the type: the pointer's zero value is `nil`.

```go
When decoding JSON into `any` in Go, the types map as follows:
bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null
```

### POST vs PUT

POST:
- Submits requests to create or update resources; not idempotent.
- For user registration, each request creates a new account; use POST.

PUT:
- Updates resources at a specific URL; idempotent.
- For changing a password, each request overwrites the same user's password; use PUT.

### Tree data for menus/resources

Three approaches:
1. Query a tree directly from MySQL (custom functions or other techniques).
2. Build the tree with recursion in code.
3. Use a single pass with a `map` to build the tree.
> Note: the `map` approach may only easily handle two levels. TODO: investigate.

### Logging

Log startup to console; write runtime logs to files.


# Nginx Deployment

HTTPS reference: [Nginx SSL certificate setup](https://cloud.tencent.com/document/product/400/35244)
