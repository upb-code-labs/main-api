# Main API

[![Integration](https://github.com/upb-code-labs/main-api/actions/workflows/integration.yaml/badge.svg?branch=dev)](https://github.com/upb-code-labs/main-api/actions/workflows/integration.yaml)
[![Coverage](https://github.com/upb-code-labs/main-api/actions/workflows/coverage.yaml/badge.svg)](https://github.com/upb-code-labs/main-api/actions/workflows/coverage.yaml)
[![Release](https://github.com/upb-code-labs/main-api/actions/workflows/release.yaml/badge.svg)](https://github.com/upb-code-labs/main-api/actions/workflows/release.yaml)

## Development 🧑🏻‍💻

### Environment variables

| Name                   | Description                                |
| ---------------------- | ------------------------------------------ |
| `DB_CONNECTION_STRING` | Connection string to the postgres database |
| `DB_MIGRATIONS_PATH`   | Absolute path to the migrations folder     |

Environment variables have default values, but you can override them by exporting them in your shell.

```bash
set NAME=VALUE
```

## Coverage 🧪

| [![Sunburst](https://codecov.io/gh/upb-code-labs/main-api/graphs/sunburst.svg?token=Q9QHF616RS)](https://app.codecov.io/gh/upb-code-labs/main-api) | [![square](https://codecov.io/gh/upb-code-labs/main-api/graphs/tree.svg?token=Q9QHF616RS)](https://app.codecov.io/gh/upb-code-labs/main-api) |
| -------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
