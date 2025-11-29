#!/bin/bash

# Install pnpm if not present
if ! command -v pnpm &> /dev/null; then
  npm install -g pnpm
else 
  echo "pnpm is installed"
fi

# Build gin-blog-front and move to Docker deploy directory
cd ../gin-blog-front
pnpm install
pnpm build
rm -rf ../deploy/build/web/dist_blog
mv ./dist ../deploy/build/web/dist_blog

# Build gin-blog-admin and move to Docker deploy directory
cd ../gin-blog-admin
pnpm install
pnpm build
rm -rf ../deploy/build/web/dist_admin
mv ./dist ../deploy/build/web/dist_admin