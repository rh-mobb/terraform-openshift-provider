# TODO

This document contains tasks requested for implementation. Each task includes detailed requirements, prerequisites, and acceptance criteria.

## Task: Implement Static Terraform Provider Registry via GitHub Pages

**Description:** Set up a self-hosted, static Terraform registry for our custom provider(s) to allow seamless installation via terraform init. This will follow a "chart-releaser" pattern where GitHub Releases store the binaries and GitHub Pages hosts the metadata API.

### Pre-requisites

- [ ] GPG Key Pair: Generate a dedicated GPG key for signing provider releases.
- [ ] GoReleaser Config: Ensure the provider repo has a `.goreleaser.yml` that generates binaries for multiple platforms and produces a signed SHA256SUMS file.

### Detailed Todo List

#### 1. Registry Infrastructure (GitHub Pages)

- [ ] Create a new repository (or use the existing provider repo) and enable GitHub Pages.
- [ ] Create a `gh-pages` branch to host the static JSON files.
- [ ] Configure a custom domain (optional) or use the default `username.github.io/repo` URL.

#### 2. Service Discovery Setup

- [ ] Add `.well-known/terraform.json` to the root of the Pages branch:

```json
{
  "providers.v1": "/v1/providers/"
}
```

#### 3. Automation Workflow (CI/CD)

- [ ] Secret Management: Add `GPG_PRIVATE_KEY`, `GPG_PASSPHRASE`, and `GPG_PUBLIC_KEY` to GitHub Actions Secrets.
- [ ] Release Pipeline: Update the `.github/workflows/release.yml` to:
  - Run goreleaser to build and upload assets to GitHub Releases.
  - Run a "Registry Releaser" step (using `tf-static-registry` or a custom script) to:
    - Generate `v1/providers/.../versions` listing the new version.
    - Generate the specific platform download JSONs pointing to the GitHub Release URLs.
    - Embed the GPG Public Key (ASCII armor) in the download JSON.
    - Commit and push these JSON updates to the `gh-pages` branch.

#### 4. Validation & Testing

- [ ] Verify Metadata: Ensure `https://<domain>/v1/providers/<org>/<name>/versions` returns valid JSON.
- [ ] End-to-End Test: Create a test Terraform configuration using the new source:

```hcl
required_providers {
  custom = {
    source  = "<domain>/<org>/<name>"
    version = "1.0.0"
  }
}
```

- [ ] Run `terraform init` and confirm the provider downloads and verifies the signature correctly.

### Acceptance Criteria

- New provider versions are automatically added to the registry upon GitHub Release.
- The registry requires no running server (100% static files).
- Terraform can successfully verify provider authenticity using the GPG public key hosted in the registry metadata.
