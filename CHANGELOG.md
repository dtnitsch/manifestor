# Changelog

All notable changes to manifestor will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-01-05

### Added
- **YAML output format support** - New default output format optimized for LLM consumption
- **CLI flag support** using urfave/cli:
  - `--root` / `-r`: Override root directory to scan
  - `--format` / `-f`: Override output format (yaml or json)
  - `--output` / `-o`: Override output file path
  - `--config`: Specify alternate config file
  - `--version`: Display version information
  - `--help`: Display help text
- **Token optimization** - YAML format saves ~26% tokens compared to JSON
- **omitempty struct tags** - Reduces manifest size by omitting default/null values
- **Query examples** - Comprehensive yq and jq query examples in `docs/examples.md`
- **Documentation** - Updated README with CLI usage, output formats, and token savings

### Changed
- **Default output format** changed from JSON to YAML (can be overridden via config or CLI)
- **Manifest version** bumped from 0.2 to 0.3
- **Config file structure** - Updated example configs to reflect YAML as default
- **CLI behavior** - Flags now properly override config values (previously `--root` was ignored)

### Fixed
- CLI flag handling - `--root` flag now correctly overrides config file root directory
- Empty field serialization - Default values (like `is_dir: false`) are now omitted

### Maintained
- **Full backward compatibility** - JSON output still fully supported
- **All existing queries work** - No changes needed to jq/yq query syntax
- **Config-first design** - CLI flags are optional overrides, not required

### Performance
- **26.3% token reduction** measured on real manifests (YAML + omitempty)
- **28.6% line reduction** in generated output files
- **No build time impact** - Compilation speed unchanged

### Migration Notes

**From v0.2 to v0.3:**

No breaking changes. Your existing setup will continue to work.

**To adopt YAML output (recommended):**
```yaml
# manifestor-config.yaml
output:
  format: "yaml"
  file: "manifest.yaml"
```

**To keep JSON output:**
```yaml
# manifestor-config.yaml
output:
  format: "json"
  file: "manifest.json"
```

**To use CLI overrides:**
```bash
# Scan different directory
./manifestor --root /path/to/project

# Force JSON output
./manifestor --format json --output manifest.json

# Both
./manifestor --root ~/projects/myapp --format yaml
```

---

## [0.2.0] - 2026-01-02

### Added
- Rollup statistics (directory-level aggregations)
- Size statistics (min, max, mean, median)
- Extension counting
- Capability-driven validation
- Manifest metadata and versioning

### Changed
- Improved scanner performance with concurrent workers
- Enhanced filter rules (block/allow patterns)

### Fixed
- Directory traversal edge cases
- Inode collection on different filesystems

---

## [0.1.0] - 2025-12-30

### Added
- Initial release
- Basic filesystem scanning
- JSON manifest output
- Configuration file support
- Filter rules (block/allow)
- Directory and file metadata collection

---

## Upgrade Guide

### v0.2 → v0.3

**No action required** - v0.3 is fully backward compatible.

**Recommended actions:**
1. Switch to YAML output for better LLM performance
2. Use CLI flags for quick overrides
3. Review query examples in `docs/examples.md`

**If you encounter issues:**
- JSON output still works exactly as before
- Set `format: "json"` in config to keep old behavior
- Report issues at https://github.com/dtnitsch/manifestor/issues

---

[0.3.0]: https://github.com/dtnitsch/manifestor/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/dtnitsch/manifestor/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/dtnitsch/manifestor/releases/tag/v0.1.0
