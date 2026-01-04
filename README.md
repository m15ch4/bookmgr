# Book API - Docker Deployment

## Application Port
The application listens on **port 8080** by default. This can be changed using the `SERVER_PORT` environment variable.

## Quick Start with Docker Compose

### Prerequisites
- Docker
- Docker Compose

### Run the Application

1. **Start all services (App + MySQL):**
```bash
docker-compose up -d
```

2. **View logs:**
```bash
docker-compose logs -f app
```

3. **Access the application:**
   - Web UI: http://localhost:8080
   - API: http://localhost:8080/api/books

4. **Stop the application:**
```bash
docker-compose down
```

5. **Clean up (remove volumes):**
```bash
docker-compose down -v
```

## Build and Run Manually

### Build the Docker Image
```bash
docker build -t bookapi:latest .
```

### Run with Existing MySQL
```bash
docker run -d \
  --name bookapi-app \
  -p 8080:8080 \
  -e DATABASE_HOST=mysql-host \
  -e DATABASE_PORT_NUMBER=3306 \
  -e DATABASE_NAME=bookdb \
  -e DATABASE_USER=bookuser \
  -e DATABASE_PASSWORD=bookpassword \
  -e SKIP_BOOTSTRAP=false \
  bookapi:latest
```

### Run with Custom Port
```bash
docker run -d \
  --name bookapi-app \
  -p 9000:9000 \
  -e DATABASE_HOST=mysql-host \
  -e DATABASE_PORT_NUMBER=3306 \
  -e DATABASE_NAME=bookdb \
  -e DATABASE_USER=bookuser \
  -e DATABASE_PASSWORD=bookpassword \
  -e SERVER_PORT=9000 \
  -e SKIP_BOOTSTRAP=false \
  bookapi:latest
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | MySQL host | localhost |
| `DATABASE_PORT_NUMBER` | MySQL port | 3306 |
| `DATABASE_NAME` | Database name | bookdb |
| `DATABASE_USER` | Database user | root |
| `DATABASE_PASSWORD` | Database password | - |
| `SKIP_BOOTSTRAP` | Skip database bootstrap | false |
| `SERVER_PORT` | Application port | 8080 |

## Docker Commands Cheatsheet
```bash
# Build image
docker build -t bookapi:latest .

# Run container
docker run -d -p 8080:8080 --name bookapi bookapi:latest

# View logs
docker logs -f bookapi

# Stop container
docker stop bookapi

# Remove container
docker rm bookapi

# Remove image
docker rmi bookapi:latest

# Execute shell in container
docker exec -it bookapi sh

# Using docker-compose
docker-compose up -d          # Start services
docker-compose down           # Stop services
docker-compose logs -f app    # View logs
docker-compose ps             # List services
docker-compose restart app    # Restart app service
```

## Multi-stage Build Benefits

The Dockerfile uses a multi-stage build which:
- Reduces final image size (Alpine-based ~20MB vs full Go image ~800MB)
- Separates build and runtime dependencies
- Improves security by running as non-root user
- Only includes necessary runtime files

## Production Deployment

For production, consider:

1. **Use specific version tags:**
```bash
docker build -t bookapi:1.0.0 .
```

2. **Enable SKIP_BOOTSTRAP after initial setup:**
```yaml
environment:
  SKIP_BOOTSTRAP: "true"
```

3. **Use Docker secrets for passwords:**
```yaml
secrets:
  - db_password

environment:
  DATABASE_PASSWORD_FILE: /run/secrets/db_password
```

4. **Add health checks:**
```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/api/books"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

## Kubernetes Deployment (Optional)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bookapi
spec:
  replicas: 3
  selector:
    matchLabels:
      app: bookapi
  template:
    metadata:
      labels:
        app: bookapi
    spec:
      containers:
      - name: bookapi
        image: bookapi:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_HOST
          value: "mysql-service"
        - name: DATABASE_PORT_NUMBER
          value: "3306"
        - name: DATABASE_NAME
          value: "bookdb"
        - name: DATABASE_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        - name: SKIP_BOOTSTRAP
          value: "true"
---
apiVersion: v1
kind: Service
metadata:
  name: bookapi-service
spec:
  selector:
    app: bookapi
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```