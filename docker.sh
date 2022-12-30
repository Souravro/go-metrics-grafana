#!/usr/bin/env bash
function setup() {
  PROJECT_ROOT_DIR=$(pwd)
  PROJECT_CONTAINER_DIR=$PROJECT_ROOT_DIR/containers

  export PROJECT_ROOT_DIR="$PROJECT_ROOT_DIR"
  export PROJECT_CONTAINER_DIR="$PROJECT_CONTAINER_DIR"
}

function createContainers() {
  #  setup
  #  container_dir=$(echo $PROJECT_CONTAINER_DIR)
  container_dir="containers"
  # run zookeeper
  docker-compose -f $container_dir/zookeeper.yml up -d
  sleep 0.5

  # run kafka broker 1
  docker-compose -f $container_dir/broker_1.yml up -d
  sleep 0.5

  # run kafka broker 2
  docker-compose -f $container_dir/broker_2.yml up -d
  sleep 0.5

  # run kafka drop
  docker-compose -f $container_dir/kafkadrop.yml up -d
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