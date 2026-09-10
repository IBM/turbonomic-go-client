# Generated API Client

This directory contains the auto-generated Go client for the Turbonomic API using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Directory Structure

```
api/
├── README.md           # This file
├── config.yaml         # oapi-codegen configuration
└── generated/          # Generated code (do not edit manually)
    └── client.gen.go   # Generated client and types
```

The upstream OpenAPI specification is read from:
`../api/com.vmturbo.api/target/swagger/openapi3/openapi3.yaml`

## Prerequisites

Install oapi-codegen:

```bash
make install-tools
```

## Generating the Client

```bash
make generate
```

To use a different spec path:

```bash
make generate OPENAPI_SPEC=/path/to/openapi3.yaml
```

### How generation works

1. **`make preprocess-spec`** — runs [`scripts/preprocess-spec.sh`](../scripts/preprocess-spec.sh) which patches the upstream spec: renames `{entity_Uuid}` path parameters, strips `null` enum values, and merges duplicate path definitions.
2. **`make generate`** — calls `oapi-codegen -config api/config.yaml` on the preprocessed spec. `compatibility.always-prefix-enum-values` prevents enum constant name collisions.

## Configuration

[`config.yaml`](config.yaml) controls code generation:

```yaml
package: generated
generate:
  client: true              # Generate HTTP client
  models: true              # Generate data models
  embedded-spec: false      # Do not embed OpenAPI spec in the binary

compatibility:
  always-prefix-enum-values: true   # Prefix all enum constants with their type name
                                    # (e.g. GroupByCriterionGroupType, not GroupType)
                                    # Prevents clashes with other generated types.

output: api/generated/client.gen.go
output-options:
  skip-prune: true
  response-type-suffix: Resp        # Rename client response wrappers to *Resp
                                    # to avoid clashes with schema types named *Response
```

## Notes

- The `generated/` directory is gitignored — it is regenerated on demand
- Always regenerate after updating the OpenAPI spec
- Do not manually edit the generated code; put customisations in the `v2` wrapper package
- Use the `v2` wrapper package for authentication, session management, and custom headers

## Known Issues

### `GroupType` / `Origin` redeclared after regeneration

**Status: resolved** — `compatibility.always-prefix-enum-values: true` in [`config.yaml`](config.yaml) makes oapi-codegen emit `GroupByCriterionGroupType` / `GroupByCriterionOrigin` directly, so the clash no longer occurs.

If you see this error after regeneration it means `always-prefix-enum-values` was removed from the config:

```
api/generated/client.gen.go: GroupType redeclared in this block
api/generated/client.gen.go: Origin redeclared in this block
```

Re-add `always-prefix-enum-values: true` under `compatibility` in `api/config.yaml` and regenerate.

### Duplicate types for client response wrapper types

**Status: resolved** — `response-type-suffix: Resp` in [`config.yaml`](config.yaml) renames the generated HTTP client response wrapper types from `*Response` to `*Resp`, keeping them distinct from any schema types of the same name.

If the error reappears after regeneration:

```
./client.gen.go: UpdateClientResponse redeclared in this block
```

Check that `response-type-suffix` is still set in `api/config.yaml`.
