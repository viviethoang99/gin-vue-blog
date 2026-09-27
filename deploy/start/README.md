It is recommended to run the `bootstrap.sh` script under the `deploy` directory. It also cleans old containers for you.

Alternatively, enter the `start` directory and run:

```bash
docker compose up -d --build
# 旧版 Docker 没有 compose 插件时用: docker-compose up -d --build
```
