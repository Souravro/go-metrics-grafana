#!/usr/bin/env bash

function createContainers() {
  # TODO: add prometheus and grafana containers
  # this function creates zookeeper, brokers and kafdrop containers and assigns a custom made network to each
  container_dir="containers"
  local_bridge="communication_bridge"

  # create a local bridge network bridge for all subsequent containers
  docker network create -d bridge $local_bridge || true

  # run zookeeper
  docker-compose -f $container_dir/zookeeper.yml up -d
  sleep 0.5
  docker network connect $local_bridge zookeeper

  # run kafka broker 1
  docker-compose -f $container_dir/broker_1.yml up -d
  sleep 0.5
  docker network connect $local_bridge broker_1

  # run kafka broker 2
  docker-compose -f $container_dir/broker_2.yml up -d
  sleep 0.5
  docker network connect $local_bridge broker_2

  # run kafka drop
  docker-compose -f $container_dir/kafdrop.yml up -d
  sleep 0.5
  docker network connect $local_bridge kafdrop

  # start application containers: consumer, producer, dashboard - in that order
  make -C consumer/
  make -C producer/
  make -C dashboard/
}

function removeContainers() {
  # remove zookeeper
  docker rm -f zookeeper

  # remove kafka broker
  docker rm -f broker_1

  # remove kafka broker
  docker rm -f broker_2

  # remove kafkadrop
  docker rm -f kafdrop

  # remove producer_service
  docker rm -f producer_service

  # remove dashboard_service
  docker rm -f dashboard_service

  # remove consumer_service
  docker rm -f consumer_service
}

function createTopic() {
   echo "Enter topic name:"
   read topic_name
   if [ $1 == 1 ]; then
      partition_count=1
      replication_factor=1
   elif [ $1 == 2 ]; then
      partition_count=1
      replication_factor=2
   elif [ $1 == 3 ]; then
      partition_count=2
      replication_factor=1
   elif [ $1 == 4 ]; then
      partition_count=2
      replication_factor=2
   fi

   echo "Creating topic with name: $topic_name, partition: $partition_count, replication_factor: $replication_factor"

   # exec into a broker and create the topic.
   docker exec broker_1 \
   kafka-topics --bootstrap-server 0.0.0.0:8097 \
               --create \
               --topic $topic_name \
               --partitions $partition_count \
               --replication-factor $replication_factor
}

function createMetricsContainer() {
    container_dir="containers"
    # run metrics
    docker-compose -f $container_dir/metrics/metrics.yml up -d
    sleep 0.5
}

if [ $1 == "create" ]; then
   createContainers
elif [ $1 == "remove" ]; then
   removeContainers
elif [ $1 == "create-topic" ]; then
   if [ -z $2 ]; then
      echo "Required argument 'type'."
      echo "type '1' : 1 partition 1 replica."
      echo "type '2' : 1 partition 2 replica."
      echo "type '3' : 2 partition 1 replica."
      echo "type '4' : 2 partition 2 replica."
      exit 0
   elif [ $2 -gt 4 ] || [ $2 -lt 1 ]; then
      echo "Invalid type. Accepted range: [1, 4]"
      exit 0
   fi
   createTopic $2
elif [ $1 == "metrics" ]; then
  createMetricsContainer
else
   echo "Invalid argument. Expected values are 'create', 'remove', 'create-topic' or 'metrics'."
fi