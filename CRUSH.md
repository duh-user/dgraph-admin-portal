## Dgraph Client CRUSH

### Build, Run, Test

- **Build:** `make dev-build`
- **Run:** `make dev-start-api`
- **Test:** No tests found. If you add tests, run with `go test ./...`

### Code Style

- **Imports:** Grouped by standard library, third-party, and internal packages.
- **Formatting:** Use `gofmt` for consistent formatting.
- **Types:** Use descriptive names for structs and interfaces.
- **Naming:** Follow Go conventions (e.g., camelCase for variables, PascalCase for exported identifiers).
- **Error Handling:** Use `log` package for logging errors.
- **Dependencies:** Use `go mod` for dependency management.
