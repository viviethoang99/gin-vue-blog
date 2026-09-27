#!/usr/bin/env sh

# Render template conf with env vars into the target file
set -eu 
if [ "$USE_HTTPS" = "true" ]; then
	envsubst '${SERVER_NAME} ${BACKEND_HOST} ${SERVER_PORT}' < /etc/nginx/conf.d/default.conf.ssl.template > /etc/nginx/conf.d/default.conf
else
	envsubst '${SERVER_NAME} ${BACKEND_HOST} ${SERVER_PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
fi
# rm /etc/nginx/conf.d/default.conf.template
# rm /etc/nginx/conf.d/default.conf.ssl.template
exec "$@"