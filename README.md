# Inventory Service

A lightweight Golang microservice for the BlueSky Logistics Platform that provides inventory management APIs.

## Overview

This service exposes two REST endpoints:
- `GET /items` - Returns a list of inventory items
- `GET /status` - Health check endpoint

## Architecture

### Technology Stack
- **Language**: Go 1.23
- **Container Runtime**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
- **Registry**: Docker Hub

## CI/CD Pipeline

The project uses a comprehensive CI/CD pipeline with three stages:

### Pipeline Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Pull Request / Push                  │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│  Stage 1: Matrix Tests (Reusable Workflow)              │
│  ------------------------------------------------        │
│  • Go Versions: 1.22, 1.23                              │
│  • OS: ubuntu-latest, macos-latest, windows-latest      │
│  • Tests: Unit tests with coverage                      │
│  • Coverage: Uploaded to Codecov                        │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│  Stage 2: Docker Build & Push (main branch only)        │
│  ------------------------------------------------        │
│  • Multi-arch: linux/amd64, linux/arm64                 │
│  • Registry: Docker Hub                                 │
│  • Tags: SHA, branch, latest                            │
│  • Cache: GitHub Actions cache                          │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│  Stage 3: Deploy to Kubernetes (main branch only)       │
│  ------------------------------------------------        │
│  • Environment: inventory-staging namespace             │
│  • Deployment: Rolling update                           │
│  • Health checks: Liveness & Readiness probes           │
└─────────────────────────────────────────────────────────┘
```

### Reusable Workflows

The pipeline leverages three reusable workflows:

#### 1. Go Matrix Test (`.github/workflows/go-matrix-test.yml`)

**Purpose**: Run tests across multiple Go versions and operating systems

**Matrix Configuration**:
```yaml
os: [ubuntu-latest, macos-latest, windows-latest]
go-version: ['1.22', '1.23']
```

**Features**:
- Parallel test execution (6 jobs: 2 Go versions × 3 OS)
- Go module caching for faster builds
- Race condition detection (`-race` flag)
- Code coverage generation
- Coverage upload to Codecov (ubuntu + Go 1.23 only)

**Inputs**:
- `working-directory`: Directory containing Go code (default: `.`)

#### 2. Build and Push Docker Image (`.github/workflows/build-and-push-image.yml`)

**Purpose**: Build multi-architecture Docker images and push to Docker Hub

**Features**:
- Multi-platform builds (amd64, arm64)
- Docker layer caching via GitHub Actions cache
- Automatic tagging strategy:
  - `main-<sha>` for main branch commits
  - `<branch>-<sha>` for feature branches
  - `latest` for main branch
- Metadata extraction for labels

**Inputs**:
- `image-name`: Full image name (required)
- `dockerfile`: Path to Dockerfile (default: `./Dockerfile`)
- `context`: Build context (default: `.`)
- `platforms`: Target platforms (default: `linux/amd64`)
- `push`: Whether to push image (default: `true`)

**Secrets Required**:
- `registry-username`: Container registry username
- `registry-password`: Container registry password/token

#### 3. Deploy to Kubernetes (`.github/workflows/deploy-k8s.yml`)

**Purpose**: Deploy application to Kubernetes cluster

**Features**:
- Namespace creation (if not exists)
- Dynamic image update in manifests
- Rollout status verification (5-minute timeout)
- Environment protection (staging)

**Inputs**:
- `image`: Full image URL with tag (required)
- `namespace`: Target Kubernetes namespace (required)
- `deployment`: Deployment name (required)
- `manifest-path`: Path to K8s manifests (default: `./k8s`)

**Secrets Required**:
- `kubeconfig-content`: Base64-encoded kubeconfig file

### How Matrix Tests Work

Matrix testing allows running the same test suite across different combinations of parameters:

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go-version: ['1.22', '1.23']
```

This creates a **test matrix** with 6 combinations:
1. ubuntu-latest + Go 1.22
2. ubuntu-latest + Go 1.23
3. macos-latest + Go 1.22
4. macos-latest + Go 1.23
5. windows-latest + Go 1.22
6. windows-latest + Go 1.23

**Benefits**:
- Ensures cross-platform compatibility
- Catches platform-specific bugs early
- Validates against multiple Go versions
- Parallel execution reduces total test time

## Secrets Management

The pipeline requires the following secrets to be configured in the repository:

| Secret Name | Purpose | Where to Get It |
|------------|---------|-----------------|
| `DOCKER_USERNAME` | Docker Hub username | Your Docker Hub username |
| `REGISTRY_USERNAME` | Docker Hub username | Your Docker Hub username |
| `REGISTRY_PASSWORD` | Docker Hub password | Docker Hub Access Token |
| `KUBECONFIG_STAGING` | Kubernetes cluster configuration | Base64-encoded kubeconfig: `cat ~/.kube/config \| base64 -w 0` |

### Setting Up Secrets

1. Navigate to your repository on GitHub
2. Go to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with its corresponding value

### Creating a Docker Hub Access Token

1. Log in to Docker Hub at https://hub.docker.com/
2. Go to **Account Settings** → **Security** → **Access Tokens**
3. Click **New Access Token**
4. Give it a description (e.g., "inventory-service-ci")
5. Select permissions: **Read, Write, Delete**
6. Click **Generate**
7. Copy the token and save it as `REGISTRY_PASSWORD` secret

### Preparing Kubeconfig Secret

```bash
# Encode your kubeconfig
cat ~/.kube/config | base64 -w 0

# Or for macOS
cat ~/.kube/config | base64

# Copy the output and save as KUBECONFIG_STAGING secret
```

## Branching Strategy

The project follows **GitHub Flow** with branch protection:

### Branch Types

- **`main`**: Production-ready code
  - Protected branch
  - Requires PR reviews
  - Requires CI checks to pass
  - No direct pushes allowed
  - Triggers Docker build and staging deployment

- **`feature/*`**: Feature development branches
  - Example: `feature/add-authentication`
  - Runs tests only
  - No deployment

### Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes & Commit**
   ```bash
   git add .
   git commit -m "Add your feature"
   git push origin feature/your-feature-name
   ```

3. **Open Pull Request**
   - Target: `main` branch
   - CI automatically runs matrix tests
   - Review required before merge

4. **Merge to Main**
   - After PR approval and passing tests
   - Automatically triggers:
     - Docker image build
     - Push to GHCR
     - Deployment to staging

### Branch Protection Rules

Configure in **Settings** → **Branches** → **Branch protection rules** for `main`:

- Require a pull request before merging
- Require approvals (1+)
- Require status checks to pass before merging
  - `test` (matrix test workflow)
- Require branches to be up to date before merging
- Do not allow bypassing the above settings

## Deployment Strategy

### Environments

- **Staging**: `inventory-staging` namespace
  - Auto-deployed on merge to `main`
  - Used for pre-production testing
  - Mirrors production configuration

### Kubernetes Resources

The deployment includes:

**Namespace** (`k8s/namespace.yaml`):
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: inventory-staging
```

**Deployment** (`k8s/deployment.yaml`):
- 2 replicas for high availability
- Resource limits: 500m CPU, 256Mi memory
- Liveness probe on `/status`
- Readiness probe on `/status`
- Rolling update strategy

**Service** (`k8s/service.yaml`):
- Type: LoadBalancer
- Port: 80 → 8080
- Exposes the application externally

### Deployment Process

1. **Trigger**: Push to `main` branch
2. **Image Build**: New image tagged with commit SHA
3. **Manifest Update**: Deployment image updated to new SHA
4. **Apply**: Manifests applied to cluster
5. **Rollout**: Rolling update with health checks
6. **Verification**: Wait for deployment to be ready (5-min timeout)

### Accessing the Service

After deployment:

```bash
# Get service external IP
kubectl get svc -n inventory-staging

# Test endpoints
curl http://<EXTERNAL-IP>/status
curl http://<EXTERNAL-IP>/items
```

## Local Development

### Prerequisites

- Go 1.22 or higher
- Docker (for local container testing)

### Running Locally

```bash
# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Run the server
go run cmd/server/main.go

# Test endpoints
curl http://localhost:8080/status
curl http://localhost:8080/items
```

### Building Docker Image

```bash
# Build
docker build -t inventory-service:local .

# Run
docker run -p 8080:8080 inventory-service:local

# Test
curl http://localhost:8080/status
```

## API Documentation

### GET /status

Health check endpoint

**Response**:
```json
{
  "status": "healthy"
}
```

### GET /items

Get inventory items

**Response**:
```json
[
  {
    "id": "1",
    "name": "Laptop",
    "quantity": 10
  },
  {
    "id": "2",
    "name": "Keyboard",
    "quantity": 25
  },
  {
    "id": "3",
    "name": "Phone",
    "quantity": 15
  }
]
```

## Project Structure

```
inventory-service/
├── .github/
│   └── workflows/
│       ├── pipeline.yml              # Main CI/CD orchestrator
│       ├── go-matrix-test.yml        # Reusable: Matrix tests
│       ├── build-and-push-image.yml  # Reusable: Docker build
│       └── deploy-k8s.yml            # Reusable: K8s deployment
├── cmd/
│   └── server/
│       └── main.go                   # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── items.go                  # Items endpoint handler
│   │   ├── status.go                 # Status endpoint handler
│   │   └── handlers_test.go          # Handler tests
│   └── models/
│       └── item.go                   # Item model
├── k8s/
│   ├── namespace.yaml                # Staging namespace
│   ├── deployment.yaml               # K8s deployment
│   └── service.yaml                  # K8s service
├── Dockerfile                        # Multi-stage Docker build
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksums
└── README.md                         # This file
```

## Bonus Features Implemented

- Unit Tests: Comprehensive test coverage for handlers
- Code Coverage: Automated coverage reporting to Codecov
- Multi-arch Builds: Support for amd64 and arm64
- Docker Layer Caching: Faster builds using GitHub Actions cache
- Health Probes: Kubernetes liveness and readiness checks
- Resource Limits: CPU and memory constraints
- Professional Documentation: Comprehensive README

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

BlueSky Logistics Platform Engineering Team

---

**Last Updated**: January 4, 2026
