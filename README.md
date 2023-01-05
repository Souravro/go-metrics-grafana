## go-metrics-grafana
This project implements 2 simple go applications to produce and consume messages respectively via kafka cluster.
Also another go application to interact with the consumer to get the consumed data. Relevant metrics in the process will be pushed to prometheus data store and will finally be visualised in grafana.

### Prerequisites
1. Docker and Docker Compose
2. Java runtime
3. Golang installed

### Getting Started
- Clone the project
- Use the following command to start kafka cluster, application containers and prometheus and grafana instances.
  ```bash
  $ sh docker.sh create
   ```
- Verify if the containers are successfully created or not. The following command should show 2 docker containers for zookeeper and docker respectively -
   ```bash
  $ docker ps
   ```
- Check container logs for each container to see if its running correctly
  ```bash
  $ docker logs <container_name>
  ```

### Modifications to application services
- Make required changes in the go applications the migrate to the application service directory and run the following command to rebuild and run it.
  ```bash
  $ make run
  ```
- Verify with `docker logs <container_name>` if the service is successfully restarted and running as expected.

### Visualise metrics
- Open `http://localhost:3000` i.e. Grafana UI and configure `http://prometheus:9090` as the data source.
- Then create a new dashboard.
- In the dashboard, create a new panel and create the query as required to visualise the graph.
- Few example queries are listed below -
  ```
   - Average response time in ms for Dashboard Service: avg(rate(dashboard_id_api_latency_sum[1m]) / rate(dashboard_id_api_latency_count[1m])) * 1000
   - Requests per Second (RPS) for dashboard service: sum(rate(dashboard_id_api_latency_count[5m]))
   - Avg no. of messages produced by Producer per sec: rate(producer_message_produced[$__rate_interval])
  ```