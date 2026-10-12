# Build the public Jaspr application as a server executable. Server-side
# rendering keeps the article list in sync with Gin instead of baking it into
# the image at build time.
FROM ghcr.io/cirruslabs/flutter:3.44.0 AS BLOG_BUILD

ARG API_BASE_URL=http://gvb-server:8765/
ARG PUBLIC_API_BASE_URL=/

ENV PATH="/root/.pub-cache/bin:${PATH}"
RUN dart pub global activate jaspr_cli 0.23.5

WORKDIR /app/blog
COPY irohasu_jaspr/pubspec.yaml irohasu_jaspr/pubspec.lock ./
RUN flutter pub get
COPY irohasu_jaspr .
RUN jaspr build \
    --dart-define="API_BASE_URL=${API_BASE_URL}" \
    --dart-define="PUBLIC_API_BASE_URL=${PUBLIC_API_BASE_URL}"

# The admin remains a Vue application under /admin.
FROM node:22-alpine AS ADMIN_BUILD
RUN npm config set registry https://registry.npmmirror.com \
    && npm install -g pnpm@11

WORKDIR /app/admin
COPY gin-blog-admin/package.json gin-blog-admin/pnpm-lock.yaml gin-blog-admin/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY gin-blog-admin .
RUN pnpm build

# Debian-based Nginx is used because the Dart executable is linked against
# glibc. Nginx proxies public pages to Jaspr and keeps /admin and /api stable.
FROM nginx:1.24.0

COPY --from=BLOG_BUILD /app/blog/build/jaspr /opt/irohasu
COPY --from=ADMIN_BUILD /app/admin/dist /usr/share/nginx/html/admin

COPY deploy/build/web/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY deploy/build/web/default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template
COPY deploy/build/web/run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh /opt/irohasu/app
ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]

EXPOSE 80
EXPOSE 443
