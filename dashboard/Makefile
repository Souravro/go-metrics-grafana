run:
	make build
	docker run -d -p 2121:2121 --name dashboard_service --network=communication_bridge dashboard_service:latest
build:
	docker build --tag dashboard_service .
