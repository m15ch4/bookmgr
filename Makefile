.PHONY: build run stop clean rebuild logs

# Build the Docker image
build:
	docker build -t bookapi:latest .

# Run with docker-compose
run:
	docker-compose up -d

# Stop containers
stop:
	docker-compose down

# Clean up (remove containers, volumes, and images)
clean:
	docker-compose down -v
	docker rmi bookapi:latest || true

# Rebuild and run
rebuild: clean build run

# View logs
logs:
	docker-compose logs -f app

# Build and run standalone (without MySQL)
run-standalone:
	docker build -t bookapi:latest .
	docker run -d \
		--name bookapi-app \
		-p 8080:8080 \
		-e DATABASE_HOST=your-mysql-host \
		-e DATABASE_PORT_NUMBER=3306 \
		-e DATABASE_NAME=bookdb \
		-e DATABASE_USER=root \
		-e DATABASE_PASSWORD=yourpassword \
		-e SKIP_BOOTSTRAP=false \
		bookapi:latest

# Stop standalone
stop-standalone:
	docker stop bookapi-app
	docker rm bookapi-app