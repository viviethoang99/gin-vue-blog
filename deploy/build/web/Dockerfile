# Jaspr's prebuilt server executable requires glibc.
FROM nginx:1.24.0
# Dockerfile for building the image from local static assets

# Copy Nginx config templates into container
COPY default.conf.template /etc/nginx/conf.d/default.conf.template
COPY default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template

# Copy the Jaspr server executable and its web assets.
ADD dist_blog/ /opt/irohasu/
ADD dist_admin/ /usr/share/nginx/html/admin

# Init script: render Nginx config from templates and env vars
COPY ./run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh /opt/irohasu/app
ENTRYPOINT ["/docker-entrypoint.sh"]

# Run on container start
CMD [ "nginx", "-g", "daemon off;" ]

EXPOSE 80
EXPOSE 443
