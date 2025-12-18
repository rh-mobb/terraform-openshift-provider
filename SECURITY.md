# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please do **not** open a public issue. Instead, please report it via one of the following methods:

1. **Email**: Send an email to the maintainers (if available)
2. **Private Security Advisory**: Use GitHub's private security advisory feature
3. **Direct Contact**: Contact the repository maintainers directly

Please include the following information:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if available)

## Security Best Practices

When using this provider:

1. **Credentials**: Never commit Kubernetes credentials or kubeconfig files to version control
2. **RBAC**: Use minimal required permissions for the service account
3. **Network**: Ensure secure network communication with the Kubernetes API
4. **Updates**: Keep the provider updated to the latest version
5. **Audit**: Regularly audit operator installations and permissions

## Security Considerations

- This provider requires access to Kubernetes/OpenShift API
- Ensure proper RBAC permissions are configured
- Use secure authentication methods (tokens, certificates)
- Consider network policies for API access
- Review operator permissions before installation

## Disclosure Policy

- We will acknowledge receipt of your vulnerability report within 48 hours
- We will provide an initial assessment within 7 days
- We will keep you informed of the progress
- We will coordinate public disclosure after a fix is available

Thank you for helping keep this project secure!
