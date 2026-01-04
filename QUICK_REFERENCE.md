# Quick Reference - Secrets Setup

## Required Secrets

Set these in GitHub Repository Settings → Secrets and variables → Actions:

### 1. DOCKER_USERNAME
```
Your Docker Hub username
```

### 2. REGISTRY_USERNAME
```
Your Docker Hub username (same as DOCKER_USERNAME)
```

### 3. REGISTRY_PASSWORD
```bash
# Create Docker Hub Access Token:
# 1. Log in to Docker Hub at https://hub.docker.com/
# 2. Account Settings → Security → Access Tokens
# 3. Click "New Access Token"
# 4. Description: "inventory-service-ci"
# 5. Permissions: Read, Write, Delete
# 6. Generate and copy the token
# 7. Use as REGISTRY_PASSWORD
```

### 4. KUBECONFIG_STAGING
```bash
# Encode your kubeconfig:

# Linux:
cat ~/.kube/config | base64 -w 0

# macOS:
cat ~/.kube/config | base64 | tr -d '\n'

# Use the output as KUBECONFIG_STAGING secret value
```

## Quick Test Commands

### Test Locally
```bash
# Run tests
go test -v ./...

# Run server
go run cmd/server/main.go

# Test endpoints
curl http://localhost:8080/status
curl http://localhost:8080/items
```

### Test Docker Build
```bash
docker build -t inventory-service:test .
docker run -p 8080:8080 inventory-service:test
```

### Test Kubernetes Manifests
```bash
# Validate manifests
kubectl apply --dry-run=client -f k8s/

# Deploy locally (if you have local k8s)
kubectl apply -f k8s/
kubectl get all -n inventory-staging
```

## Workflow Trigger Conditions

| Workflow | Trigger | Notes |
|----------|---------|-------|
| Matrix Tests | All branches | 6 jobs (2 Go × 3 OS) |
| Docker Build | `main` only | After tests pass |
| K8s Deploy | `main` only | After Docker build |
| Release | Tags `v*` | Creates GitHub Release |

## Common Commands

### Branch Management
```bash
# Create feature branch
git checkout -b feature/my-feature

# Push and create PR
git push origin feature/my-feature
# Then create PR on GitHub

# After merge, update local main
git checkout main
git pull origin main
```

### Check Pipeline Status
```bash
# View in browser
open https://github.com/YOUR_USERNAME/inventory-service/actions

# Or use GitHub CLI
gh run list
gh run watch
```

### Tag a Release
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Troubleshooting Quick Fixes

### "Unexpected value 'go-version-file'"
Fixed - Removed from workflow inputs

### "authentication required" for Docker Hub
- Check REGISTRY_USERNAME is your Docker Hub username
- Check REGISTRY_PASSWORD is valid access token
- Ensure token hasn't been revoked
- Verify DOCKER_USERNAME secret is set

### Tests fail on Windows
- Add `.gitattributes` with `*.go text eol=lf`

### Deployment fails
- Verify KUBECONFIG_STAGING is base64-encoded correctly
- Test: `echo "$SECRET" | base64 -d | kubectl --kubeconfig=/dev/stdin get nodes`

## File Structure Summary

```
inventory-service/
├── .github/workflows/
│   ├── pipeline.yml                # Main orchestrator
│   ├── go-matrix-test.yml          # Matrix testing
│   ├── build-and-push-image.yml    # Docker build
│   ├── deploy-k8s.yml              # K8s deployment
│   └── release.yml                 # Semantic versioning
├── cmd/server/main.go              # Application entry
├── internal/handlers/              # HTTP handlers + tests
├── k8s/                            # K8s manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   └── namespace.yaml
├── Dockerfile                      # Multi-stage build
├── README.md                       # Main documentation
├── SETUP.md                        # Setup instructions
└── SUBMISSION_CHECKLIST.md         # Submission guide
```

---

**Last Updated**: January 4, 2026
