# Versioning Guide

This document explains how versioning works for the Terraform Provider for OpenShift.

## Versioning Strategy

We follow [Semantic Versioning (SemVer)](https://semver.org/) with the format: `MAJOR.MINOR.PATCH`

- **MAJOR** (x.0.0): Breaking changes that require user action
- **MINOR** (0.x.0): New features, backward compatible
- **PATCH** (0.0.x): Bug fixes, backward compatible

### Current Version

The current version is **0.1.0** (initial release).

## Version Sources

The provider version is determined from multiple sources:

1. **Git Tags** (primary): Version is extracted from git tags (e.g., `v0.1.0`)
2. **Build-time variables**: Set via ldflags during compilation
3. **Metadata file**: `provider/.terraform-registry/metadata.json`

## How to Version a Release

### 1. Update Version References

Before creating a release, update version references:

```bash
# Update metadata.json
vim provider/.terraform-registry/metadata.json
# Change "version": "0.1.0" to new version

# Update CHANGELOG.md
vim CHANGELOG.md
# Add new version section with changes
```

### 2. Create Git Tag

Create a git tag with the version (prefixed with `v`):

```bash
# For a new release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# For a patch release
git tag -a v0.1.1 -m "Release v0.1.1 - Bug fixes"
git push origin v0.1.1
```

### 3. GitHub Release

The GitHub Actions workflow will automatically:
- Build binaries for all platforms
- Create checksums
- Sign with GPG
- Create a GitHub release

Or create manually:

```bash
# Create release via GitHub CLI
gh release create v0.1.0 \
  --title "v0.1.0" \
  --notes "$(cat CHANGELOG.md | sed -n '/## \[0.1.0\]/,/## \[/p')" \
  provider/dist/*
```

### 4. Terraform Registry

The Terraform Registry will automatically detect the GitHub release and publish it.

## Version in Code

The version is embedded in the binary at build time:

```go
// provider/version/version.go
var Version = "dev"  // Set via ldflags
```

Build with version:

```bash
make build VERSION=v0.1.0
```

Or manually:

```bash
go build -ldflags "-X github.com/redhat/terraform-provider-openshift-operator/provider/version.Version=v0.1.0" \
  -o terraform-provider-openshift
```

## Version Constraints

Users specify version constraints in their Terraform configuration:

```hcl
terraform {
  required_providers {
    openshift = {
      source  = "registry.terraform.io/rh-mobb/openshift"
      version = "~> 0.1"  # Allows 0.1.x, not 0.2.0
    }
  }
}
```

Common version constraints:
- `~> 0.1` - Allows patch updates (0.1.x)
- `~> 0.1.0` - Exact version (0.1.0 only)
- `>= 0.1.0, < 0.2.0` - Range
- `>= 0.1.0` - Minimum version

## Release Checklist

- [ ] Update `CHANGELOG.md` with new version section
- [ ] Update `provider/.terraform-registry/metadata.json` version
- [ ] Update version in examples if needed
- [ ] Run tests: `make test`
- [ ] Build locally: `make build`
- [ ] Create git tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Verify GitHub release was created
- [ ] Verify Terraform Registry published the release
- [ ] Update documentation if needed

## Pre-1.0.0 Versions

During the 0.x.x phase:
- **0.1.0** → **0.2.0**: New features (minor)
- **0.1.0** → **0.1.1**: Bug fixes (patch)
- Breaking changes are allowed but should be documented

## Post-1.0.0 Versions

After 1.0.0:
- **1.0.0** → **2.0.0**: Breaking changes (major)
- **1.0.0** → **1.1.0**: New features (minor)
- **1.0.0** → **1.0.1**: Bug fixes (patch)
- Breaking changes require major version bump

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## References

- [Semantic Versioning](https://semver.org/)
- [Terraform Provider Versioning](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [Go Module Versioning](https://go.dev/ref/mod#versions)
