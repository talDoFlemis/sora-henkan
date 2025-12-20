# Sora Henkan (空変換)

Sora Henkan is a powerful and scalable image processing service designed to handle on-the-fly image transformations. It provides a simple API to submit image URLs, apply transformations like scaling, and retrieve the processed images. The system is built with a decoupled architecture using a message queue, ensuring high throughput and resilience.

## Features

- **Asynchronous Image Processing:** Jobs are queued and processed in the background by workers, preventing API blocking.
- **Image Scaling:** Resize images by specifying target width and height.
- **Real-time Updates:** Subscribe to real-time progress updates for image processing jobs via Server-Sent Events (SSE).
- **Cloud-Native:** Designed to run on the cloud with infrastructure-as-code for AWS.
- **Local Development:** A complete Docker Compose setup for easy local development and testing.
- **Observability:** Integrated OpenTelemetry for tracing and metrics with Jaeger for distributed tracing visualization.
- **API Documentation:** Swagger/OpenAPI 3.0 documentation for easy API exploration and testing.
- **Health Monitoring:** Comprehensive health check endpoints for both API and worker services.
- **AMI-Based Deployment:** Fast instance launches using pre-built Amazon Machine Images.
- **Load Testing:** Built-in load testing capabilities for performance validation.
- **CloudWatch Integration:** AWS CloudWatch agent integration for detailed monitoring and logging.
- **CDN Support:** Cloudflare client IP extraction for accurate request tracking behind CDN.

## Architecture

Sora Henkan follows a distributed and decoupled architecture:

### Conceptual Architecture

```mermaid
graph TD
    subgraph "User"
        A[Client]
    end

    subgraph "Sora Henkan System"
        B[API Server]
        C[Message Queue]
        D[Worker]
        E[RDBMS]
        F[Object Storage]
    end

    A -- "1. POST /images (Image URL)" --> B
    B -- "2. Save Job (pending)" --> E
    B -- "3. Publish Job" --> C
    C -- "4. Consume Job" --> D
    D -- "5. Fetch Image from URL" --> G[Internet]
    D -- "6. Upload Original Image" --> F
    D -- "7. Process Image (e.g., scale)" --> D
    D -- "8. Upload Transformed Image" --> F
    D -- "9. Update Job Status (e.g., completed)" --> E
```

### AWS Implementation (`simple-aws-architecture`)

This diagram shows how the conceptual architecture is implemented using AWS services as defined in the Terraform configuration.

```mermaid
graph TD
    subgraph "User"
        A[Client]
    end

    subgraph "AWS Cloud"
        subgraph "VPC"
            subgraph "Public Subnet"
                B[Application Load Balancer]
            end

            subgraph "Private Subnet"
                subgraph "API Auto Scaling Group"
                    C[EC2 Instances for API]
                end
                subgraph "Worker Auto Scaling Group"
                    D[EC2 Instances for Worker]
                end
                E[RDS PostgreSQL]
                F[S3 Bucket]
                G[SQS Queue]
            end
        end
    end

    A -- HTTPS --> B
    B -- Forwards traffic to --> C
    C -- "Writes Job & Publishes to" --> G
    C -- "Reads/Writes Job Info" --> E
    D -- "Consumes from" --> G
    D -- "Reads/Writes Job Info" --> E
    D -- "Uploads/Reads Images" --> F
```

1. **API Server (`api`):** A Go server built with the Echo framework. It exposes a RESTful API to accept image processing requests. When a new request is received, it saves the job details to the PostgreSQL database and publishes a message to an SQS queue.
2. **Worker (`worker`):** A Go application that listens for messages on the SQS queue. When a message is received, it fetches the image from the source URL, performs the required transformations (e.g., scaling) using `libvips`, and uploads the original and transformed images to an S3-compatible object storage.
3. **Frontend:** A React application (built with Vite) that provides a user interface to interact with the API.
4. **PostgreSQL:** The primary database for storing information about image processing jobs.
5. **MinIO/S3:** An object storage service used to store the original and processed images.
6. **LocalStack/AWS SQS:** A message queue used to decouple the API server from the worker, enabling asynchronous processing.
7. **Terraform:** Infrastructure as Code (IaC) to provision the necessary AWS resources, including VPC, EC2, RDS, S3, and SQS.
8. **Jaeger:** Distributed tracing system for monitoring and troubleshooting microservices-based architectures.
9. **CloudWatch Agent:** AWS monitoring and observability service for collecting metrics and logs from EC2 instances.

### Kubernetes Implementation

This diagram shows how the conceptual architecture is implemented in Kubernetes using the `task k8s:all` deployment.

```mermaid
graph TD
    subgraph "User"
        A[Client]
    end

    subgraph "k3d Cluster"
        subgraph "Istio Service Mesh"
            B[Gateway API]
        end

        subgraph "Application Namespace"
            subgraph "Backend Helm Chart"
                C[API Deployment]
                D[Worker Deployment]
                E[Migrate Job]
            end

            subgraph "Frontend Helm Chart"
                F[Frontend Deployment]
            end

            subgraph "Infrastructure Services"
                G[PostgreSQL StatefulSet]
                H[MinIO StatefulSet]
                I[RabbitMQ StatefulSet]
                J[LocalStack Deployment]
            end
        end

        subgraph "Observability"
            K[Kubernetes Dashboard]
        end
    end

    A -- HTTP/HTTPS --> B
    B -- Routes to --> C
    B -- Routes to --> F
    B -- Routes to --> H
    B -- Routes to --> I
    C -- Publishes Jobs --> I
    C -- Reads/Writes --> G
    D -- Consumes Jobs --> I
    D -- Reads/Writes --> G
    D -- Stores Images --> H
    E -- Runs Migrations --> G
    F -- API Calls --> C
```

#### Kubernetes Components

1. **k3d Cluster:** A lightweight Kubernetes cluster running locally via k3d (k3s in Docker), configured with port mappings for services (80, 443, 9001, 9000, 15672, 42069).

2. **Istio Service Mesh:**

   - **Gateway API:** Modern Kubernetes ingress using Gateway API resources (v1.4.0) for traffic routing
   - **Control Plane (istiod):** Manages service mesh configuration and routing rules
   - **HTTPRoute Resources:** Define routing rules for frontend, API, MinIO, and RabbitMQ management interfaces

3. **Backend Helm Chart (`helm/backend`):**

   - **API Deployment:** Go-based REST API server with configurable replicas and auto-scaling (HPA)
   - **Worker Deployment:** Image processing workers with auto-scaling capabilities
   - **Migrate Job:** Kubernetes Job that runs database migrations on deployment
   - **Service Accounts:** IAM for Kubernetes with proper RBAC configurations
   - **ConfigMaps & Secrets:** Environment-based configuration for databases, message queues, and object storage

4. **Frontend Helm Chart (`helm/frontend`):**

   - **Deployment:** React/Vite application served via Nginx
   - **Service:** ClusterIP service exposing the frontend internally
   - **HTTPRoute:** Routes traffic from the gateway to the frontend service

5. **Infrastructure Services:**

   - **PostgreSQL:** StatefulSet deployment for persistent database storage (job metadata, processing status)
   - **MinIO:** S3-compatible object storage deployed as StatefulSet with persistent volumes for images
   - **RabbitMQ:** Message broker using AMQP protocol for pub/sub messaging (alternative to AWS SQS)
   - **LocalStack:** AWS service emulator for local development (optional, provides SQS/S3 emulation)

6. **Message Queue Options:**

   - **RabbitMQ (AMQP):** Used in Kubernetes deployments for reliable message delivery with durable pub/sub configuration
   - **LocalStack SQS:** Alternative queue option for AWS-compatible workflows

7. **Observability:**
   - **Kubernetes Dashboard:** Web-based UI for cluster management and monitoring
   - **Service Mesh Telemetry:** Built-in observability from Istio for traffic monitoring

#### Key Features in Kubernetes Deployment

- **Declarative Infrastructure:** All resources defined in Helm charts for reproducible deployments
- **Auto-scaling:** Horizontal Pod Autoscalers (HPA) for API and Worker based on CPU utilization
- **Gateway API:** Modern, role-oriented API for traffic management replacing traditional Ingress
- **StatefulSets:** Persistent storage for databases and stateful services
- **Health Checks:** Liveness and readiness probes for all deployments
- **Resource Limits:** CPU and memory limits defined for efficient resource utilization
- **Secret Management:** Kubernetes Secrets for sensitive credentials (database, AMQP, object storage)
- **Service Mesh:** Istio provides traffic management, security, and observability out of the box

## Getting Started

### Prerequisites

- Go (version 1.21 or higher)
- Node.js and pnpm
- Docker and Docker Compose
- [Task](https://taskfile.dev/installation/)

### Running with Docker Compose

1. **Clone the repository:**

   ```sh
   git clone https://github.com/taldoflemis/sora-henkan.git
   cd sora-henkan
   ```

2. **Start all services:**
   This command will start the API server, worker, frontend, PostgreSQL, MinIO, and LocalStack using Docker Compose.

   ```sh
   task
   ```

The services will be available at:

- **Frontend:** `http://localhost:8080`
- **API Server:** `http://localhost:42069`
- **API Documentation (Swagger UI):** `http://localhost:42069/swagger/index.html`
- **MinIO Console:** `http://localhost:9001`

### Running with Kubernetes

You can also run the application on a local Kubernetes cluster using k3d. This simulates a production-like environment.

1. **Prerequisites:**

   - [k3d](https://k3d.io/)
   - [Helm](https://helm.sh/)
   - [kubectl](https://kubernetes.io/docs/tasks/tools/)

2. **Create cluster and deploy services:**

   ```sh
   task k8s:all
   ```

   This command creates a k3d cluster and deploys all services including dependencies (PostgreSQL, MinIO, RabbitMQ, LocalStack) and the application (Frontend, Backend).

3. **Access Services:**

   - **Frontend:** `http://localhost:80`
   - **API Server:** `http://localhost:42069`
   - **MinIO Console:** `http://localhost:9001`
   - **RabbitMQ Management:** `http://localhost:15672`

4. **Clean up:**

   ```sh
   task k8s:delete-cluster
   ```

## Technology Stack

- **Backend:** Go
  - **Web Framework:** [Echo](https://echo.labstack.com/)
  - **Messaging:** [Watermill](https://watermill.io/) (with AWS SQS)
  - **Image Processing:** [govips](https://github.com/davidbyttow/govips)
  - **Database:** PostgreSQL with [pgx](https://github.com/jackc/pgx)
  - **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Frontend:**
  - **Framework:** React with Vite
  - **Language:** TypeScript
- **Infrastructure:**
  - **Containerization:** Docker, Docker Compose
  - **IaC:** Terraform
  - **Cloud:** AWS (EC2, RDS, S3, SQS, ALB, CloudWatch)
  - **CI/CD:** GitHub Actions (not yet implemented)
- **Observability:** OpenTelemetry, Jaeger
- **Load Testing:** k6

## Project Structure

```
.
├── cmd/                # Main applications (api, worker, migrate)
├── deploy/             # Deployment scripts and configurations
├── docker/             # Dockerfiles for different services
├── frontend/           # React frontend application
├── internal/           # Core application logic
│   ├── core/
│   │   ├── application/ # Use cases and business logic
│   │   └── domain/      # Domain entities and types
│   └── infra/           # Infrastructure adapters (database, messaging, etc.)
├── settings/           # Application configuration
└── terraform/          # Terraform code for AWS infrastructure
```

## API Endpoints

- `POST /v1/images`: Create a new image processing request.
- `GET /v1/images`: List all image processing jobs.
- `GET /v1/images/:id`: Get details for a specific image processing job.
- `GET /v1/images/sse`: Get real-time updates for all jobs.
- `GET /v1/images/:id/sse`: Get real-time updates for a specific job.
- `GET /health`: API server health check endpoint.
- `GET /swagger/*`: Swagger UI for interactive API documentation (OpenAPI 3.0).

### API Documentation

The API is fully documented using OpenAPI 3.0 specification. You can explore and test the API interactively using Swagger UI:

- **Local development:** `http://localhost:42069/swagger/index.html`
- **Regenerate API documentation:** `task swagger:generate`

The OpenAPI specification file is located at `docs/openapi.yaml`.

### Example Request

```sh
curl -X POST http://localhost:42069/v1/images -H "Content-Type: application/json" -d '{
  "image_url": "https://www.nasa.gov/sites/default/files/styles/full_width_feature/public/thumbnails/image/j2m-shareable.jpg",
  "scale_transformation": {
    "enabled": true,
    "width": 200,
    "height": 200
  }
}'
```

## Configuration

The application is configured using YAML files located in the `settings/` directory. The `base.yaml` file contains default settings, which can be overridden by environment-specific files or environment variables.

For example, to change the database host for the API, you can set the `API_DATABASE_HOST` environment variable.

## Load Testing

The project includes a load testing script using k6 located in the `loadgenerator/` directory. This allows you to test the performance and scalability of the image processing service.

To run load tests:

```sh
cd loadgenerator
docker build -t sora-henkan-loadgen .
docker run --rm sora-henkan-loadgen
```

The load test script (`script.js`) includes a set of anime image URLs and simulates realistic user traffic patterns.

## Observability

### Distributed Tracing with Jaeger

The application integrates Jaeger for distributed tracing, allowing you to visualize request flows across the API and worker services. In AWS deployments, Jaeger UI is exposed through the Application Load Balancer.

### CloudWatch Integration

For AWS deployments, the CloudWatch agent collects metrics and logs from EC2 instances, providing comprehensive monitoring and alerting capabilities.

### OpenTelemetry

All services are instrumented with OpenTelemetry for standardized observability, including:

- Request tracing
- Performance metrics
- Custom application metrics

## Deployment

The infrastructure for this project is defined using Terraform in the `terraform/simple-aws-architecture/` directory. The Terraform code provisions all the necessary resources on AWS to run the application in a scalable and resilient manner.

### Key Deployment Features

- **AMI-Based Instances:** Uses pre-built Amazon Machine Images for faster instance launches and consistent deployments
- **Auto Scaling:** Automatic scaling for both API and worker instances based on load
- **Health Checks:** Comprehensive health monitoring for both API and worker services through ALB
- **Application Load Balancer:** Routes traffic to API instances and provides access to Jaeger UI
- **Monitoring:** Integrated CloudWatch agent for metrics and logs collection
- **Optimized Costs:** Instance type optimizations and resource management to reduce AWS costs

## Recent Changes (Since v1.7.0)

### New Features

- **Swagger/OpenAPI 3.0 Documentation** (v1.21.0, v1.22.0): Complete API documentation with interactive Swagger UI
- **Better Orchestration** (v1.22.0): Improved container orchestration and deployment workflows
- **AMI Image Builder** (v1.16.0): Pre-built Amazon Machine Images for faster instance launches and deployments
- **Worker Health Endpoints** (v1.17.0, v1.18.0): Comprehensive health check endpoints for worker services exposed on ALB
- **Load Testing Framework** (v1.15.0): Built-in k6 load testing scripts with anime image dataset
- **Jaeger Integration** (v1.10.0-v1.13.0): Distributed tracing with Jaeger UI accessible through ALB
- **CloudWatch Agent** (v1.13.0): AWS CloudWatch integration for enhanced monitoring and logging
- **Cloudflare Support** (v1.14.0): Client IP extraction for accurate request tracking behind Cloudflare CDN
- **MIME Type Handling** (v1.9.0): Proper MIME type detection and extension mapping for saved images
- **S3 Enhancements** (v1.10.0, v1.11.0): Public bucket access and CORS policy configuration

### Infrastructure Improvements

- **AMI Builder EC2** (v1.16.0): Automated AMI creation for consistent instance deployments
- **Container Resource Limits** (v1.19.0): Proper resource requests and limits for containerized services
- **Cloud-Init Integration** (v1.19.0): Always-run cloud-init scripts for instance configuration
- **Cost Optimizations** (v1.20.0): Reduced costs by removing unnecessary services and optimizing instance types
- **NAT Gateway Sequencing** (v1.20.0): Proper dependency management for OpenTelemetry collector instances

### Bug Fixes and Optimizations

- **OpenTelemetry Logging** (v1.20.0): Removed AWS EMF exporter to reduce costs
- **Load Test Script** (v1.16.2): Updated to use inline URLs for better reliability
- **Instance Scaling** (v1.17.0): Improved handling of worker instance scaling and health checks
- **Docker Compose Optimization** (v1.16.1): Streamlined AMI builder to use docker compose pull

For a complete changelog, see [CHANGELOG.md](CHANGELOG.md).
