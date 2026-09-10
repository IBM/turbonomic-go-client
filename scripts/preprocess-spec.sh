#!/bin/bash
# Preprocess the upstream OpenAPI spec before passing it to oapi-codegen.
#
# Fixes applied:
#   1. Rename all snake_case path parameters to camelCase in both path template
#      keys and parameter name fields (e.g. {group_Uuid} → {groupUuid},
#      {market_Uuid} → {marketUuid}).  The one exception handled first is
#      /groups/{entity_Uuid}/tags → /groups/{group_Uuid}/tags since that
#      parameter represents a group UUID, not a generic entity UUID.
#   2. Remove literal `null` enum values (not valid Go identifiers).
#   3. Merge duplicate path definitions via a PyYAML round-trip.

set -e

SOURCE_SPEC="${1:-../api/com.vmturbo.api/target/swagger/openapi3/openapi3.yaml}"
OUT_SPEC="${2:-/tmp/openapi3-preprocessed.yaml}"

if [ ! -f "$SOURCE_SPEC" ]; then
    echo "Error: spec not found at $SOURCE_SPEC"
    exit 1
fi

echo "Preprocessing OpenAPI spec..."
echo "  Source : $SOURCE_SPEC"
echo "  Output : $OUT_SPEC"

cp "$SOURCE_SPEC" "$OUT_SPEC"

# 1. Rename path parameters + merge duplicates + remove null enums via Python.
#    Python handles all renames in one pass to avoid ordering issues with sed.
python3 - "$OUT_SPEC" <<'PYTHON_SCRIPT'
import sys
import re
import yaml

spec_path = sys.argv[1]

def to_camel(name):
    """Convert snake_case or mixed-case_Suffix to camelCase.
    e.g. group_Uuid -> groupUuid, entity_Uuid -> entityUuid, ad_group_Uuid -> adGroupUuid
    """
    parts = name.split("_")
    return parts[0] + "".join(p.capitalize() for p in parts[1:])

def rename_path_params(text):
    """Replace {snake_case} placeholders in path strings with {camelCase}."""
    def replacer(m):
        inner = m.group(1)
        if "_" in inner:
            return "{" + to_camel(inner) + "}"
        return m.group(0)
    return re.sub(r"\{([^}]+)\}", replacer, text)

try:
    with open(spec_path) as f:
        content = f.read()

    # Step 1a: special-case rename before any other entity_Uuid processing.
    # /groups/{entity_Uuid}/tags → /groups/{group_Uuid}/tags
    # (group_Uuid will be renamed to groupUuid in the general pass below)
    content = content.replace(
        "/groups/{entity_Uuid}/tags",
        "/groups/{group_Uuid}/tags"
    )

    spec = yaml.safe_load(content)

    # Step 1b: rename all snake_case path template keys and parameter name fields.
    if "paths" in spec:
        new_paths = {}
        for path_key, path_item in spec["paths"].items():
            new_key = rename_path_params(path_key)
            if path_item is None:
                new_paths[new_key] = path_item
                continue

            # Rename parameter name fields at path level
            for param in path_item.get("parameters", []):
                if isinstance(param, dict) and param.get("in") == "path":
                    old_name = param.get("name", "")
                    if "_" in old_name:
                        param["name"] = to_camel(old_name)

            # Rename parameter name fields at operation level
            for method in ("get", "post", "put", "patch", "delete", "head", "options", "trace"):
                op = path_item.get(method)
                if not isinstance(op, dict):
                    continue
                for param in op.get("parameters", []):
                    if isinstance(param, dict) and param.get("in") == "path":
                        old_name = param.get("name", "")
                        if "_" in old_name:
                            param["name"] = to_camel(old_name)

            new_paths[new_key] = path_item

        spec["paths"] = new_paths

    # Step 2: remove null from enum lists.
    def remove_null_enums(obj):
        if isinstance(obj, dict):
            if "enum" in obj and isinstance(obj["enum"], list):
                obj["enum"] = [v for v in obj["enum"] if v is not None]
            for v in obj.values():
                remove_null_enums(v)
        elif isinstance(obj, list):
            for item in obj:
                remove_null_enums(item)

    remove_null_enums(spec)

    # Step 3: write back (PyYAML merges duplicate keys automatically on load).
    with open(spec_path, "w") as f:
        yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

    print("  ✓ Renamed snake_case path parameters to camelCase")
    print("  ✓ Removed null enum values")
    print("  ✓ Merged duplicate path definitions")

except Exception as e:
    print(f"Error during preprocessing: {e}", file=sys.stderr)
    sys.exit(1)
PYTHON_SCRIPT

echo "✓ Preprocessed spec written to $OUT_SPEC"
