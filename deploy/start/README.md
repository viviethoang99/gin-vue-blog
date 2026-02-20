It is recommended to run the `bootstrap.sh` script under the `deploy` directory. It also cleans old containers for you.

Alternatively, enter the `start` directory and run:

```bash
docker-compose up -d
```

This setup now includes a `gvb-caddy` service as the public entrypoint:

- `http://localhost` (or your `CADDY_SITE_ADDRESS`) -> frontend (`gvb-web`)
- `/api/*` -> backend (`gvb-server`)
