#!/bin/bash

# Install the frontend CLIs if not present.
if ! command -v pnpm &> /dev/null; then
  npm install -g pnpm
else 
  echo "pnpm is installed"
fi
command -v flutter >/dev/null 2>&1 || { echo "flutter is required" >&2; exit 1; }
if ! command -v jaspr >/dev/null 2>&1; then
  dart pub global activate jaspr_cli 0.23.5
fi

# Build the Jaspr public site and move the complete server bundle.
cd ../irohasu_jaspr
flutter pub get
docker_arch="$(docker info --format '{{.Architecture}}' 2>/dev/null || uname -m)"
case "$docker_arch" in
  arm64|aarch64) jaspr_arch=arm64 ;;
  amd64|x86_64) jaspr_arch=x64 ;;
  *) echo "unsupported Docker architecture: $docker_arch" >&2; exit 1 ;;
esac
jaspr build \
  --target-os=linux \
  --target-arch="$jaspr_arch" \
  --dart-define=API_BASE_URL=http://gvb-server:8765/ \
  --dart-define=PUBLIC_API_BASE_URL=/
rm -rf ../deploy/build/web/dist_blog
mv ./build/jaspr ../deploy/build/web/dist_blog

# Build gin-blog-admin and move to Docker deploy directory
cd ../gin-blog-admin
pnpm install
pnpm build
rm -rf ../deploy/build/web/dist_admin
mv ./dist ../deploy/build/web/dist_admin
