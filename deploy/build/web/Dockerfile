FROM nginx:1.25.3-alpine
# Dockerfile for building the image from local static assets

# Copy Nginx config templates into container
COPY default.conf.template /etc/nginx/conf.d/default.conf.template
COPY default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template

# Copy static assets into container
ADD dist_blog/ /usr/share/nginx/html/
ADD dist_admin/ /usr/share/nginx/html/admin

# Init script: render Nginx config from templates and env vars
COPY ./run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

# Run on container start
CMD [ "nginx", "-g", "daemon off;" ]

EXPOSE 80
EXPOSE 443