# CLAUDE.md - Guidelines for KarpelesLab/ini

## Build/Test Commands
- Build project: `make` (IMPORTANT: Always run this before committing)
- Run all tests: `make test` or `go test -v`
- Run single test: `go test -v -run TestName` (e.g., `go test -v -run TestIni`)
- Check test coverage: `go test -cover` or `go test -coverprofile=coverage.out && go tool cover -func=coverage.out`
- Install dependencies: `make deps`
- Format code: `goimports -w -l .`

## Workflow Guidelines
- Always run `make` before committing to ensure proper formatting and successful build
- Run `make test` to verify that tests still pass after your changes
- Consider checking code coverage with `go test -cover` to ensure adequate test coverage

## Code Style Guidelines
- Import formatting: Standard Go grouping (stdlib first, then external)
- Error handling: Return errors up the call stack, no panics
- Naming: Use idiomatic Go (camelCase for private, PascalCase for public)
- Functions should have comments in godoc format
- Variables: Lowercase section and key names (use `strings.ToLower()`)
- Indentation: Tabs (not spaces)
- Line length: Keep reasonable (<100 chars where possible)
- Use `map[string]map[string]string` structure for ini data
- Return explicit boolean success flags with values (e.g., `string, bool`)
- Implement standard Go interfaces (`io.ReaderFrom`, `io.WriterTo`) where applicable
- Mark deprecated methods with `// Deprecated:` comments and indicate alternatives