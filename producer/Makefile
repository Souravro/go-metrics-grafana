run:
	make build
	docker run -d -p 8181:8181 --name producer_service --network=communication_bridge producer_service:latest
build:
	docker build --tag producer_service .