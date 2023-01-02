run:
	make build
	docker run -d --name producer_service --network=communication_bridge producer_service:latest
build:
	docker build --tag producer_service .