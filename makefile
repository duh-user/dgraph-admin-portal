DGRAPH := dgraph/standalone:latest
RATEL  := dgraph/ratel:latest
CERTS_DIR := ./certs

# dev-docker: pull required docker images
dev-pull:
	docker pull $(DGRAPH)
	docker pull $(RATEL)


#dev-docker: build the docker containers
build-dev-dgraph:
	@echo "Starting dev dgraph-local"
	docker run --name dgraph-local \
		-it -p 5080:5080 -p 6080:6080 -p 9999:8080 -p 9080:9080 \
		-v ./store:/dgraph $(DGRAPH)
	@echo "Starting dev dgraph ratel ui"
	docker run --name ratel-local \
        --platform linux/amd64 -d -p "8000:8000" $(RATEL)


#start-dev-dgraph: run the docker containers
start-dev-dgraph:
	@echo "Strating dev dgraph-local"
	docker start dgraph-local
	@echo "Starting dev dgraph ratel ui"
	docker start ratel-local


# generate SSL certificates for local dev
generate-certs:
	@echo "Generating self-signed certs"
	openssl genpkey -algorithm RSA -out $(CERTS_DIR)/app.key
	openssl req -new -key $(CERTS_DIR)/app.key -out $(CERTS_DIR)/app.csr \
		-subj "/C=US/ST=Alabama/L=McMullen/O=Local Corp/OU=LC/CN=local.corp/emailAddress=admin@local.corp"
	openssl x509 -req -in $(CERTS_DIR)/app.csr -signkey $(CERTS_DIR)/app.key -out $(CERTS_DIR)/app.crt -days 365
	@echo "If any issues with certificates check $(CERTS_DIR)/"


# build the dev app
dev-build:
	@echo "Build dev"
	cd ./app/ && go build -o ../dgraph-client 
	

# run go dev app
dev-start-api:
	make dev-build
	@echo "starting dev api server"
	./dgraph-client api start


help:
	@echo "dev-pull          -  pull docker images"
	@echo "build-dev-dgraph  -  build dgraph dev containers"
	@echo "start-dev-dgraph  -  start dgraph dev containers"
	@echo "generate-certs    -  generate self-signed certs"
	@echo "dev-start-api     -  start the dev api with default values"

