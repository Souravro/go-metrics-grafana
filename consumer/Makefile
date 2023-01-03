run:
	make build
	docker run -d -p 8080:8080 --name consumer_service --network=communication_bridge consumer_service:latest
build:
	docker build --tag consumer_service .
