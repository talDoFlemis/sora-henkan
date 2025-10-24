# Sora Henkan (空変換)

Sora Henkan is a powerful and scalable image processing service designed to handle on-the-fly image transformations. It provides a simple API to submit image URLs, apply transformations like scaling, and retrieve the processed images. The system is built with a decoupled architecture using a message queue, ensuring high throughput and resilience.

## Features

- **Asynchronous Image Processing:** Jobs are queued and processed in the background by workers, preventing API blocking.
- **Image Scaling:** Resize images by specifying target width and height.
- **Real-time Updates:** Subscribe to real-time progress updates for image processing jobs via Server-Sent Events (SSE).
- **Cloud-Native:** Designed to run on the cloud with infrastructure-as-code for AWS.
- **Local Development:** A complete Docker Compose setup for easy local development and testing.
- **Observability:** Integrated OpenTelemetry for tracing and metrics.

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

### AWS Implementation (`shitty-aws-architecture`)

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

1.  **API Server (`api`):** A Go server built with the Echo framework. It exposes a RESTful API to accept image processing requests. When a new request is received, it saves the job details to the PostgreSQL database and publishes a message to an SQS queue.
2.  **Worker (`worker`):** A Go application that listens for messages on the SQS queue. When a message is received, it fetches the image from the source URL, performs the required transformations (e.g., scaling) using `libvips`, and uploads the original and transformed images to an S3-compatible object storage.
3.  **Frontend:** A React application (built with Vite) that provides a user interface to interact with the API.
4.  **PostgreSQL:** The primary database for storing information about image processing jobs.
5.  **MinIO/S3:** An object storage service used to store the original and processed images.
6.  **LocalStack/AWS SQS:** A message queue used to decouple the API server from the worker, enabling asynchronous processing.
7.  **Terraform:** Infrastructure as Code (IaC) to provision the necessary AWS resources, including VPC, EC2, RDS, S3, and SQS.

## Getting Started

### Prerequisites

-   Go (version 1.21 or higher)
-   Node.js and pnpm
-   Docker and Docker Compose
-   [Task](https://taskfile.dev/installation/)

### Local Development

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/taldoflemis/sora-henkan.git
    cd sora-henkan
    ```

2.  **Start all services:**
    This command will start the API server, worker, frontend, PostgreSQL, MinIO, and LocalStack using Docker Compose.
    ```sh
    task all:up
    ```

3.  **Create the MinIO bucket:**
    You need to create the bucket where images will be stored.
    ```sh
    task minio:create-bucket
    ```

4.  **Run database migrations:**
    This will set up the necessary tables in the PostgreSQL database.
    ```sh
    task migrate:up
    ```

The services will be available at:
-   **Frontend:** `http://localhost:8080`
-   **API Server:** `http://localhost:42069`
-   **MinIO Console:** `http://localhost:9001`

## Technology Stack

-   **Backend:** Go
    -   **Web Framework:** [Echo](https://echo.labstack.com/)
    -   **Messaging:** [Watermill](https://watermill.io/) (with AWS SQS)
    -   **Image Processing:** [govips](https://github.com/davidbyttow/govips)
    -   **Database:** PostgreSQL with [pgx](https://github.com/jackc/pgx)
    -   **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
-   **Frontend:**
    -   **Framework:** React with Vite
    -   **Language:** TypeScript
-   **Infrastructure:**
    -   **Containerization:** Docker, Docker Compose
    -   **IaC:** Terraform
    -   **Cloud:** AWS (EC2, RDS, S3, SQS, ALB)
    -   **CI/CD:** GitHub Actions (not yet implemented)
-   **Observability:** OpenTelemetry

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

-   `POST /v1/images`: Create a new image processing request.
-   `GET /v1/images`: List all image processing jobs.
-   `GET /v1/images/:id`: Get details for a specific image processing job.
-   `GET /v1/images/sse`: Get real-time updates for all jobs.
-   `GET /v1/images/:id/sse`: Get real-time updates for a specific job.

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

## Deployment

The infrastructure for this project is defined using Terraform in the `terraform/shitty-aws-architecture/` directory. The Terraform code provisions all the necessary resources on AWS to run the application in a scalable and resilient manner.
