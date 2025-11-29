#!/bin/bash
# Clean Docker containers built by this project

stop_and_remove_container() {
  container_name=$1
  if docker ps -a --format '{{.Names}}' | grep -q "$container_name"; then
    docker stop $container_name
    docker rm $container_name
    echo "Stopped and removed container $container_name"
  else
    echo "Container $container_name not found"
  fi
}

stop_and_remove_container "gvb-web"
stop_and_remove_container "gvb-server"
stop_and_remove_container "gvb-redis"
stop_and_remove_container "gvb-mysql"

echo "Cleanup of project Docker containers completed"