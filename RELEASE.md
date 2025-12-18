# Release Process

This document describes the process for creating a new release of the Terraform Provider for OpenShift.

## Prerequisites

- Write access to the repository
- GitHub CLI (`gh`) installed (optional, for easier releases)

## Release Steps

### 1. Prepare Release

```bash
# Ensure you're on main branch and up to date
git checkout main
git pull origin main

# Run tests
cd provider
make test

# Build and verify
make build
make version
```

### 2. Update Version Files

```bash
# Update metadata.json
vim provider/.terraform-registry/metadata.json
# Change version to new version (e.g., "0.1.0")

# Update CHANGELOG.md
vim CHANGELOG.md
# Add new version section with all changes
```

### 3. Commit Changes

```bash
git add provider/.terraform-registry/metadata.json CHANGELOG.md
git commit -m "chore: prepare release v0.1.0"
git push origin main
```

### 4. Create Git Tag

```bash
# Create annotated tag
git tag -a v0.1.0 -m "Release v0.1.0

See CHANGELOG.md for details."

# Push tag (triggers release workflow)
git push origin v0.1.0
```

### 5. Verify Release

1. Check GitHub Actions workflow completed successfully
2. Verify GitHub release was created with binaries
3. Check Terraform Registry published the release (may take a few minutes)

### 6. Announce Release (Optional)

- Update release notes on GitHub
- Announce in relevant channels/forums
- Update documentation if needed

## Quick Release Script

```bash
#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 0.1.0"
  exit 1
fi

# Update metadata.json
sed -i.bak "s/\"version\": \".*\"/\"version\": \"${VERSION}\"/" provider/.terraform-registry/metadata.json
rm provider/.terraform-registry/metadata.json.bak

# Commit changes
git add provider/.terraform-registry/metadata.json
git commit -m "chore: prepare release v${VERSION}"

# Create and push tag
git tag -a "v${VERSION}" -m "Release v${VERSION}"
git push origin main
git push origin "v${VERSION}"

echo "Release v${VERSION} triggered. Check GitHub Actions for progress."
```

## Version Numbering

Follow Semantic Versioning:
- **0.1.0** → **0.2.0**: New features
- **0.1.0** → **0.1.1**: Bug fixes
- **1.0.0** → **2.0.0**: Breaking changes

See [VERSIONING.md](VERSIONING.md) for details.

## Troubleshooting

### Release workflow failed
- Check GitHub Actions logs
- Ensure tag format is correct (`v0.1.0`)

### Terraform Registry not updating
- Wait a few minutes (can take up to 10 minutes)
- Verify release has binaries attached
- Check metadata.json version matches tag

### Version mismatch
- Ensure metadata.json version matches git tag (without `v` prefix)
- Verify CHANGELOG.md has correct version
