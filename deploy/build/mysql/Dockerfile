FROM mysql:8.0

# Working directory for initialization assets
ENV WORK_PATH /usr/local/work
# Directory auto-executed by the container on init
ENV AUTO_RUN_DIR /docker-entrypoint-initdb.d
# Init script filename
ENV RUN_SHELL run.sh

COPY ./mysql.cnf /etc/mysql/conf.d/
# Copy initial data SQL into working directory
COPY ./gvb.sql ${WORK_PATH}/
# Place init script into the auto-run directory
COPY ./${RUN_SHELL} ${AUTO_RUN_DIR}/

# Ensure init script is executable
RUN chmod a+x ${AUTO_RUN_DIR}/${RUN_SHELL}