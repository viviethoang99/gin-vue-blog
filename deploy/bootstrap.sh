#!/bin/bash

if [ "$1" == "dev" ]; then
  # Development: use local pnpm to build (relative to start/Dockerfile)
  export WEB_BUILD_CONTEXT="../build/web"
  ./build_web.sh
else
  # Production: use Docker container's Node to build
  export WEB_BUILD_CONTEXT="../.."
fi

# Clean old containers
./clean_docker.sh

# Start new containers
cd start
docker-compose up -d --build