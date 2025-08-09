# Titanic Passenger Data Service

This project provides a Go-based web service for querying passenger data from the Titanic dataset. It supports multiple data sources (CSV and SQLite), is fully containerized with Docker, and can be deployed to different Kubernetes environments (Kind or Docker Desktop) using a Helm chart.

## Features

- **RESTful API**: Exposes API endpoints to query passenger data.
- **Multiple Data Sources**: Can be configured to read data from either a CSV file or a SQLite database at deployment time.
- **API Documentation**: Automatically generates interactive API documentation using Swagger (OpenAPI).
- **Containerized**: Fully containerized using Docker for both the application and its data, following a clean separation of concerns.
- **Kubernetes Ready**: Includes a flexible Helm chart for easy deployment to a Kubernetes cluster.
- **Versatile Local Development**: Streamlined local development workflow using `make` that supports both `kind` and `docker-desktop` as a Kubernetes environment.

## Prerequisites

Before you begin, ensure you have the following tools installed:

- **Go**: Version 1.24 or later
- **Docker**: For building and running containers
- **Kind**: (Optional) For running a local Kubernetes cluster with `kind`
- **Helm**: For managing Kubernetes deployments
- **Make**: For running the streamlined commands in the `Makefile`

## Project Structure

```

/titanic-go-service
|-- cmd/server/main.go       # Main application entry point
|-- cmd/seed/main.go         # Script to seed the SQLite DB
|-- internal/                # Internal application logic (handlers, data, models)
|-- docs/                    # Auto-generated Swagger documentation
|-- helm/titanic-chart/      # Helm chart for Kubernetes deployment
|-- test/                    # Unit and functional tests
|-- Dockerfile               # Dockerfile for the Go service
|-- Dockerfile.data          # Dockerfile for the data container
|-- Makefile                 # Command shortcuts for development and deployment
|-- README.md                # This file
|-- config.yaml              # Local configuration (not used in Docker/K8s)
|-- go.mod & go.sum          # Go module files
|-- titanic.csv              # The raw dataset
|-- titanic.db               # The SQLite database file

````

## Quickstart: Local Kubernetes Deployment

This is the fastest way to get the entire application running on a local Kubernetes cluster.

1.  **Clone the repository.**
    ```bash
    git clone git@github.com:dhope-nagesh/titanic-go-service.git
    cd titanic-go-service
    ```
2.  **Choose your Kubernetes environment and run the installation command.**

    **For Docker Desktop(Default):**
    Ensure Kubernetes is enabled in Docker Desktop settings, then run:
    CSV as a data source:
    ```bash
    make install
    ```
    SQLite as a data source:
    ```bash
    make install DATA_SOURCE=sqlite
    ```

    **For Kind:**
    This single command will build the Docker images, create a local `kind` cluster, load the images into it, and deploy the application using Helm.
    CSV as a data source:
    ```bash
    make install K8S_ENV=kind
    ```
    SQLite as a data source:
    ```bash
    make install K8S_ENV=kind DATA_SOURCE=sqlite
    ```


3.  **Access the Service:**
    Once the pods are running, you can expose the service using `kubectl port-forward` (see section below).

4.  **Access Swagger Docs:**
    Navigate to the following URL in your browser (assuming you are port-forwarding to `8080`):
    `http://127.0.0.1:8080/swagger/index.html`

## Local Testing with Port-Forwarding

The most reliable way to connect to your service running in the cluster is to forward a local port to the Kubernetes **Service**. This automatically routes traffic to a healthy pod.

1.  **Get the name of the service:**
    The service name is determined by your Helm release name. If you used the `Makefile`, it will be `titanic-release-titanic-go-service`. You can confirm by running:
    ```bash
    kubectl get services
    ```

2.  **Run the `port-forward` command:**
    Open a **new terminal window** and run the following command. This will create a persistent connection from your local machine to the service.
    ```bash
    # The format is 'svc/<service-name>'
    kubectl port-forward svc/titanic-release-titanic-go-service 8080:8080
    ```
    Note: We forward local port `8080` to the service's port `80`, which then routes to the container's port `8080`. The command will appear to hang, which means the connection is active.

3.  **Test the API locally:**
    While the port-forward is running, open **another terminal** and you can now test the API using `curl` against `localhost:8080`.
    ```bash
    # Get all passengers
    curl http://127.0.0.1:8080/api/v1/passengers

    # Get a specific passenger
    curl http://127.0.0.1:8080/api/v1/passengers/5

    # Get specific attributes for a passenger
    curl -G http://127.0.0.1:8080/api/v1/passengers/2/attributes \
    --data-urlencode "attributes=Name" \
    --data-urlencode "attributes=Age" \
    --data-urlencode "attributes=Fare"
    ```

## Development Workflow using `make`

The `Makefile` provides several commands to streamline development.

-   **`make help`**
    Displays a list of all available commands and their descriptions.
-   **`make build`**
    Builds both the Go service and the data container Docker images. You can override the tag: `make build TAG=v1.1.0`.
-   **`make test`**
    Runs all unit and functional tests in the project.
-   **`make install`**
    The main command to deploy everything to your chosen Kubernetes environment. You can configure the deployment by overriding variables:
    ```bash
    # Deploy to Docker Desktop with a specific tag and use the SQLite data source
    make install K8S_ENV=docker-desktop TAG=v1.1.0 DATA_SOURCE=sqlite
    ```
-   **`make uninstall`**
    Removes the Helm release from the Kubernetes cluster.
-   **`make clean`**
    Deletes the local `kind` cluster if `K8S_ENV=kind`.

## API Endpoints

The service exposes the following endpoints under the base path `/api/v1`.

| Method | Path                                   | Description                                                  |
| :----- | :------------------------------------- | :----------------------------------------------------------- |
| `GET`  | `/passengers`                          | Returns a list of all passengers.                            |
| `GET`  | `/passengers/{id}`                     | Returns all data for a single passenger by their ID.         |
| `GET`  | `/passengers/{id}/attributes`          | Returns specific attributes for a passenger. (e.g., `?attributes=Name&attributes=Age`) |
| `GET`  | `/stats/fare_histogram`                | Returns data for a histogram of fare prices by percentile.   |
