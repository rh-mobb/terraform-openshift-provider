# AI Agent Guide for Terraform Provider OpenShift

This guide helps AI agents understand the repository structure, conventions, and best practices for contributing to this Terraform provider.

## Repository Overview

This is a **Terraform Provider** for managing OpenShift operators. The provider name is `openshift` (registry: `rh-mobb/openshift`).

### Key Facts
- **Provider Name**: `openshift` (not `openshift-operator`)
- **Registry Path**: `registry.terraform.io/rh-mobb/openshift`
- **Current Version**: 0.1.0
- **Language**: Go 1.21+
- **Terraform SDK**: Plugin SDK v2

## Directory Structure

```
rh-mobb_openshift-provider/
├── provider/              # Main provider code
│   ├── internal/         # Internal packages (not exported)
│   │   ├── client/       # Kubernetes client wrapper
│   │   ├── provider/     # Provider schema and configuration
│   │   ├── resources/    # Resource implementations
│   │   └── version/      # Version information
│   ├── docs/             # Documentation
│   │   ├── resources/    # Resource documentation
│   │   ├── data-sources/# Data source documentation
│   │   └── guides/      # Usage guides
│   ├── main.go           # Provider entry point
│   ├── Makefile          # Build commands
│   └── go.mod            # Go dependencies
├── examples/             # Example Terraform configurations
├── website/              # Terraform Registry website docs
├── temp/                 # Temporary files (NOT committed)
├── .github/              # GitHub workflows and templates
└── [root docs]          # README, CHANGELOG, etc.
```

## Important Conventions

### 1. Provider vs Resource Names
- **Provider name**: `openshift` (used in `provider "openshift" {}`)
- **Resource name**: `openshift_operator` (used in `resource "openshift_operator" {}`)
- **Data source name**: `openshift_operator` (used in `data "openshift_operator" {}`)

### 2. Registry References
- Always use: `registry.terraform.io/rh-mobb/openshift`
- Never use: `registry.terraform.io/redhat/openshift-operator` (old name)

### 3. Version References
- Git tags: `v0.1.0` (with `v` prefix)
- Metadata: `0.1.0` (without `v` prefix)
- Version constraints: `~> 0.1` (allows 0.1.x)

### 4. Temp Directory
- **`./temp/`** is for temporary files that should NOT be committed
- Use for: drafts, scratch files, temporary documentation
- Files here are ignored by git (via .gitignore)

## Versioning Process

### Current Version: 0.1.0

### Versioning Rules (SemVer)
- **MAJOR** (x.0.0): Breaking changes
- **MINOR** (0.x.0): New features, backward compatible
- **PATCH** (0.0.x): Bug fixes, backward compatible

### Files to Update for Releases
1. **CHANGELOG.md** - Move items from `[Unreleased]` to new version section
2. **provider/.terraform-registry/metadata.json** - Update `version` field
3. **Git tag** - Create annotated tag `v0.1.0`

### Release Process
1. Update CHANGELOG.md with release date
2. Update metadata.json version
3. Commit changes
4. Create git tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
5. Push tag: `git push origin v0.1.0`
6. GitHub Actions automatically builds and releases

See `RELEASE.md` for detailed release instructions.
See `VERSIONING.md` for detailed versioning guide.

## Code Style & Patterns

### Go Code
- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Add comments for exported functions/types
- Use descriptive error messages

### Resource Implementation Pattern
```go
// Schema definition in operator.go
func ResourceOperator() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceOperatorCreate,
        ReadContext:   resourceOperatorRead,
        UpdateContext: resourceOperatorUpdate,
        DeleteContext: resourceOperatorDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: map[string]*schema.Schema{...},
    }
}

// Implementation in operator_impl.go
func resourceOperatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // Implementation
}
```

### Error Handling
- Use `diag.Diagnostics` for errors
- Provide detailed error messages with context
- Include resource names and namespaces in errors
- Use `diag.FromErr()` for simple errors
- Use `diag.Diagnostics` with `diag.Diagnostic` for complex errors

### Testing
- Unit tests: `*_test.go` files
- Acceptance tests: `*_acceptance_test.go` files (require `TF_ACC=1`)
- Test naming: `TestResourceName` or `TestFunctionName`
- Run tests: `make test` or `go test ./...`

## Common Tasks

### Adding a New Resource
1. Create `provider/internal/resources/resource_name.go` (schema)
2. Create `provider/internal/resources/resource_name_impl.go` (CRUD)
3. Create `provider/internal/resources/resource_name_test.go` (tests)
4. Register in `provider/internal/provider/provider.go` ResourcesMap
5. Add documentation in `provider/docs/resources/resource_name.md`
6. Add website docs in `website/docs/r/resource_name.html.markdown`
7. Add example in `examples/`

### Adding a New Data Source
1. Create `provider/internal/resources/resource_name_data.go`
2. Create `provider/internal/resources/resource_name_data_test.go`
3. Register in `provider/internal/provider/provider.go` DataSourcesMap
4. Add documentation in `provider/docs/data-sources/resource_name.md`
5. Add website docs in `website/docs/d/resource_name.html.markdown`

### Updating Documentation
- **Resource docs**: `provider/docs/resources/`
- **Data source docs**: `provider/docs/data-sources/`
- **Guides**: `provider/docs/guides/`
- **Website docs**: `website/docs/` (for Terraform Registry)
- **Examples**: `examples/` directory

### Building & Testing
```bash
cd provider
make build          # Build provider binary
make test           # Run tests
make install        # Install to local Terraform plugins
make fmt            # Format code
make lint           # Run linter
make version        # Show version info
```

## Best Practices Checklist

When making changes, ensure:

- [ ] Code follows Go best practices
- [ ] Tests added/updated for new functionality
- [ ] Documentation updated (resource/data source docs)
- [ ] Examples updated if needed
- [ ] CHANGELOG.md updated for user-facing changes
- [ ] No hardcoded credentials or sensitive data
- [ ] Error messages are user-friendly and informative
- [ ] Version references are correct (rh-mobb/openshift)
- [ ] Provider name is `openshift` (not `openshift-operator`)

## File Locations Reference

### Core Code
- Provider entry: `provider/main.go`
- Provider schema: `provider/internal/provider/provider.go`
- Resources: `provider/internal/resources/`
- Client: `provider/internal/client/client.go`
- Version: `provider/internal/version/version.go`

### Documentation
- Main README: `README.md`
- Provider README: `provider/README.md`
- Resource docs: `provider/docs/resources/`
- Data source docs: `provider/docs/data-sources/`
- Troubleshooting: `provider/docs/guides/troubleshooting.md`
- Website docs: `website/docs/`

### Configuration
- Go modules: `provider/go.mod`
- Build config: `provider/Makefile`
- Registry metadata: `provider/.terraform-registry/metadata.json`
- CI/CD: `.github/workflows/`
- Pre-commit: `.pre-commit-config.yaml`

### Versioning
- Changelog: `CHANGELOG.md`
- Versioning guide: `VERSIONING.md`
- Best practices: `BEST_PRACTICES.md`

## Important Notes

1. **Never commit to `temp/` directory** - It's for temporary files only
2. **Always use `rh-mobb/openshift`** - Never reference old `redhat/openshift-operator`
3. **Provider name is `openshift`** - Resource name is `openshift_operator`
4. **Version tags use `v` prefix** - `v0.1.0` not `0.1.0`
5. **Update CHANGELOG.md** - For any user-facing changes
6. **Run tests before committing** - Use `make test`
7. **Follow SemVer** - For version bumps

## Testing Requirements

- Unit tests for all new functionality
- Acceptance tests for resources (if possible)
- Test error cases and edge cases
- Mock external dependencies where appropriate
- Clean up test resources

## Release Checklist

When preparing a release:
1. Update CHANGELOG.md (move Unreleased to version)
2. Update metadata.json version
3. Run all tests: `make test`
4. Build and verify: `make build`
5. Create git tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
6. Push tag: `git push origin v0.1.0`
7. Verify GitHub Actions workflow
8. Check Terraform Registry (5-10 min delay)

## Common Mistakes to Avoid

1. ❌ Using `openshift-operator` instead of `openshift`
2. ❌ Using `redhat` namespace instead of `rh-mobb`
3. ❌ Forgetting to update CHANGELOG.md
4. ❌ Not running tests before committing
5. ❌ Committing files to `temp/` directory
6. ❌ Using wrong version format (missing `v` prefix in tags)
7. ❌ Not updating both resource docs and website docs

## Getting Help

- Check `BEST_PRACTICES.md` for best practices checklist
- Check `VERSIONING.md` for versioning details
- Check `CONTRIBUTING.md` for contribution guidelines
- Check `PLAN.md` for original design document

## Quick Reference

```bash
# Build
cd provider && make build

# Test
cd provider && make test

# Install locally
cd provider && make install

# Show version
cd provider && make version

# Format code
cd provider && make fmt

# Lint code
cd provider && make lint
```

---

**Remember**: This is a Terraform provider for OpenShift operators. The provider name is `openshift`, and it's published to `registry.terraform.io/rh-mobb/openshift`.
