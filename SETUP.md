# Setup Guide for Inventory Service CI/CD

This guide walks you through setting up all required secrets and configurations for the CI/CD pipeline.

## Prerequisites

- GitHub account
- Docker Hub account
- Access to a Kubernetes cluster
- GitHub repository for inventory-service

## Step 1: Create Docker Hub Access Token

## Step 1: Create Docker Hub Access Token

The access token is used to authenticate with Docker Hub for pushing Docker images.

### Instructions:

1. Log in to Docker Hub at https://hub.docker.com/
2. Go to **Account Settings** → **Security**
3. Click **New Access Token**
4. Configure the token:
   - **Description**: `inventory-service-ci`
   - **Access permissions**: **Read, Write, Delete**
5. Click **Generate**
6. **IMPORTANT**: Copy the token immediately (you won't see it again)

## Step 2: Prepare Kubernetes Configuration

### Get your kubeconfig

```bash
# View your current kubeconfig
cat ~/.kube/config

# Or if using a specific cluster
kubectl config view --flatten --minify
```

### Encode kubeconfig for GitHub Secret

```bash
# For Linux (no line wrapping)
cat ~/.kube/config | base64 -w 0

# For macOS (base64 doesn't have -w flag)
cat ~/.kube/config | base64 | tr -d '\n'

# Copy the output - this is your KUBECONFIG_STAGING value
```

**Example output:**
```
YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF...
```

## Step 3: Add Secrets to GitHub Repository

### Navigate to Repository Secrets

1. Go to your repository: `https://github.com/YOUR_USERNAME/inventory-service`
2. Click **Settings** tab
3. In the left sidebar, click **Secrets and variables** → **Actions**
4. Click **New repository secret**

### Add Each Secret

#### Secret 1: REGISTRY_USERNAME

- **Name**: `REDOCKER_USERNAME

- **Name**: `DOCKER_USERNAME`
- **Value**: Your Docker Hub username (e.g., `yourusername`)
- Click **Add secret**

#### Secret 2: REGISTRY_USERNAME

- **Name**: `REGISTRY_USERNAME`
- **Value**: Your Docker Hub username (same as DOCKER_USERNAME)
- Click **Add secret**

#### Secret 3: REGISTRY_PASSWORD

- **Name**: `REGISTRY_PASSWORD`
- **Value**: The Docker Hub access token you created in Step 1
- Click **Add secret**

#### Secret 4KUBECONFIG_STAGING`
- **Value**: The base64-encoded kubeconfig from Step 2
- Click **Add secret**

### Verify Secrets

After adding all secrets, you should see:

```
DOCKER_USERNAME          Updated X seconds ago
REGISTRY_USERNAME        Updated X seconds ago
REGISTRY_PASSWORD        Updated X seconds ago
KUBECONFIG_STAGING       Updated X seconds ago
```

## Step 4: Configure Docker Hub Repository
 Step 4: Configure Docker Hub Repository

### Create Repository (if not auto-created)

1. Log in to Docker Hub
2. Click **Create Repository**
3. Repository name: `inventory-service`
4. Visibility: **Public** or **Private** (as needed)
5. Click **Create**

The repository will be auto-created on first push if it doesn't exist.-service` repository
4. Click **Connect**

## Step 5: Configure Branch Protection for `main`

### Enable Branch Protection

1. In repository, go to **Settings** → **Branches**
2. Click **Add rule** or **Add branch protection rule**
3. Configure as follows:

**Branch name pattern**: `main`

**Protect matching branches**:
- ✅ Require a pull request before merging
  - ✅ Require approvals: **1**
  - ✅ Dismiss stale pull request approvals when new commits are pushed
  - ✅ Require review from Code Owners (optional)
  Require a pull request before merging
  - Require approvals: **1**
  - Dismiss stale pull request approvals when new commits are pushed
  - Require review from Code Owners (optional)
  
- Require status checks to pass before merging
  - Require branches to be up to date before merging
  - **Status checks**: Add `test` (this appears after first workflow run)
  
- Require conversation resolution before merging

-Step 6: Test the Pipeline

### Create a Feature Branch

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/inventory-service.git
cd inventory-service

# Create a feature branch
git checkout -b feature/test-pipeline

# Make a small change
echo "# Test" >> README.md

# Commit and push
git add README.md
git commit -m "Test: Verify CI/CD pipeline"
git push origin feature/test-pipeline
```

### Create a Pull Request

1. Go to your repository on GitHub
2. Click **Pull requests** → **New pull request**
3. Set base: `main` ← compare: `feature/test-pipeline`
4. Click **Create pull request**
5. Observe the CI checks running:
   - Matrix tests (6 jobs)
   - Docker build (skipped - not on main)
   - Deploy (skipped - not on main)

### Merge to Main

1. After tests pass and review is approved
2. Click **Merge pull request**
3. Observe the full pipeline:
   - Matrix tests
   - Docker build and push
   - Deploy to staging

## Step 7: Verify Deployment

### Check Docker Image

```bash
# Pull the image (make sure it's public or you're authenticated)
docker pull YOUR_DOCKERHUB_USERNAME/inventory-service:latest

# Run locally
docker run -p 8080:8080 YOUR_DOCKERHUB_USERNAME/inventory-service:latest

# Test
curl http://localhost:8080/status
```

### Check Kubernetes Deployment

```bash
# Check namespace
kubectl get namespace inventory-staging

# Check deployment
kubectl get deployment -n inventory-staging

# Check pods
kubectl get pods -n inventory-staging

# Check service
kubectl get svc -n inventory-staging

# Get service URL (if LoadBalancer)
kubectl get svc inventory-service -n inventory-staging -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# Test the service
curl http://<EXTERNAL-IP>/status
curl http://<EXTERNAL-IP>/items
```

## Step 8: Create a Release (Bonus Feature)

```bash
# Tag a release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# This triggers:
# - GoReleaser workflow
# - Multi-arch binaries
# - GitHub Release with artifacts
```

## Troubleshooting

### Pipeline Fails with "Error: Invalid workflow file"

- Check YAML syntax using: https://www.yamllint.com/
- Ensure all required secrets are set
- Verify workflow files have correct indentation

### Docker Push Fails with "authentication required"

- Verify `REGISTRY_USERNAME` is correct (your GitHub username)
- Verify `REGISTRY_PASSWORD` is a valid PAT with `write:packages` scope
- Check PAT hasn't expired

### Kubernetes Deployment Fails

- Verify `KUBECONFIG_STAGING` is base64-encoded correctly:
  ```bash
  echo "$KUBECONFIG_STAGING" | base64 -d | kubectl --kubeconfig=/dev/stdin get nodes
  ```
- Ensure the kubeconfig has permissions for the target cluster
- Check cluster is accessible from GitHub Actions runners

### Matrix Tests Fail on Windows

- Windows may have different line endings (CRLF vs LF)
- Add `.gitattributes`:
  ```
  * text=auto
  *.go text eol=lf
  ```

## Security Best Practices

1. Never commit secrets to the repository
2. Rotate access tokens regularly (every 90 days)
3. Use least-privilege for kubeconfig (namespace-scoped ServiceAccount)
4. Enable audit logging on Kubernetes cluster
5. Review workflow runs regularly for suspicious activity
6. Use environment protection rules for production deployments

## Next Steps

1. Set up secrets
2. Enable branch protection
3. Test with a PR
4. Merge to main
5. Verify deployment
6. Add monitoring and alerting
7. Set up staging to production promotion workflow
8. Add Slack/Discord notifications

## Need Help?

- **GitHub Actions Documentation**: https://docs.github.com/en/actions
- **Kubernetes Documentation**: https://kubernetes.io/docs/
- **Docker Hub Documentation**: https://docs.docker.com/docker-hub/

---

**Last Updated**: January 4, 2026
