# Docker variables #THIS CONFIGURATION IS FOR POSTGRES IT MAY BE DIFFERENT FOR OTHER DATABASES
CONTAINER_NAME=postgres-latest
IMAGE_NAME=postgres:latest
# Database variables
USER=root
OWNER=root
PASSWORD=@case123!
DATABASE_NAME=my_bank
SSL_MODE=disable
PORT=5432

# Start PostgreSQL container
postgres:
	docker run --name $(CONTAINER_NAME) -p $(PORT):$(PORT) -e POSTGRES_USER=$(USER) -e POSTGRES_PASSWORD=$(PASSWORD) -d $(IMAGE_NAME)

# Create database
createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=$(USER) --owner=$(USER) $(DATABASE_NAME)

# Drop database
dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb $(DATABASE_NAME)

#Create migration
migrateup:
	migrate -path db/migration -database "postgres://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DATABASE_NAME)?sslmode=$(SSL_MODE)" -verbose up

#Remove migration
migratedown:
	migrate -path db/migration -database "postgres://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DATABASE_NAME)?sslmode=$(SSL_MODE)" -verbose down

#Run SQLC
sqlc:
	sqlc generate

# Stop PostgreSQL container
stop:
	docker stop $(CONTAINER_NAME)

# Restart PostgreSQL container
restart:
	docker restart $(CONTAINER_NAME)

# Ensure the commands are not mistaken for files
.PHONY: postgres createdb dropdb stop restart
