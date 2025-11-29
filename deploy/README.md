# Latest: available scripts

- `build_web.sh`: Build web projects locally (requires Node) and move built assets into the container build context
- `clean_docker.sh`: Clean old Docker containers related to this project
- `bootstrap.sh`: Use a Docker Node image to build frontend static assets
- `bootstrap.sh dev`: Use your local `pnpm` to build frontend assets

In most cases, just run `bootstrap.sh`. It will clean old containers, build the latest code, and rebuild the images each time. The first run can be slow, later runs use cache and are much faster.

---

One-click run consists of two steps:

1. Environment setup: Ensure Docker and Docker Compose are installed
2. Start: Execute the `bootstrap.sh` script
3. Troubleshooting: If it fails, check the notes below

## 1. Environment Setup

## Windows / macOS

Most people use Windows daily and can run this project with Docker on Windows to preview it quickly.

If you plan to deploy to a cloud server, Linux is recommended.

Install [Docker Desktop](https://www.docker.com/products/docker-desktop/). It includes both Docker and Docker Compose.

## Linux

If you use a desktop Linux, you can also install [Docker Desktop](https://www.docker.com/products/docker-desktop/).

This guide assumes **Ubuntu**. Other distributions may differ slightly. Refer to official docs when in doubt.

#### 1) Install required dependencies

```bash
sudo apt update && sudo apt install -y vim curl git
```

#### 2) Install Docker

Always refer to the official docs for the latest steps: https://docs.docker.com/engine/install/ubuntu/

1. Update and install prerequisites

```bash
sudo apt-get update

sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
```

2. Add Docker’s official GPG key

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```

3. Set up the repository

```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

4. Install Docker Engine

```bash
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

5. Verify installation

```bash
sudo docker run hello-world
```

#### 3) Install Docker Compose (standalone)

Refer to: https://docs.docker.com/compose/install/other/

1. Download binary

```bash
sudo curl -SL https://github.com/docker/compose/releases/download/v2.14.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
# If slow, try a mirror:
# sudo curl -SL https://get.daocloud.io/docker/compose/releases/download/v2.14.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
```

2. Make it executable

```bash
sudo chmod +x /usr/local/bin/docker-compose
```

## 2. Run

## Quick preview without changing source code

```bash
# Clone project
git clone https://github.com/szluyu99/gin-vue-blog
cd gin-vue-blog

cd deploy
./bootstrap.sh
```

If you just want to preview, use the defaults above.

For production, edit important values (like DB passwords) in `start/.env`.

## After modifying source code

Backend: Modify `gin-blog-server` directly; the backend image is built from its Dockerfile.

Admin frontend (`gin-blog-admin`): After building, copy the `dist` output to `build/web/dist_admin`.

Blog frontend (`gin-blog-front`): After building, copy the `dist` output to `build/web/dist_blog`.

Then rebuild and run under `start/`:

```bash
docker-compose up -d --build
```

The above is automated in `build_web.sh`. Just run the script.

## Production notes

Production deployment is simply running this project on a server. Two places are recommended to adjust:

1. Before `docker-compose up -d`, edit the environment in `.env`:

For security, change at least:

- `REDIS_PASSWORD`
- `MYSQL_ROOT_PASSWORD`

Other values are optional per your needs.

2. The backend image builds directly from the `gin-blog-server` source and loads `config/config.docker.toml`.

Review and adjust values there as needed (see file comments).

---

Default admin user is `admin / 123456`. Important: After startup, log into the admin and change the password.

## Troubleshooting

## gvb-mysql and gvb-server fail to start

If you already have a local MySQL using port 3306, `gvb-mysql` may fail to start. Since `gvb-server` depends on it, it will fail too.

- Option A: Stop your local MySQL service
- Option B: Change `MYSQL_PORT` in `.env` to a different port

Note: The changed port is the host-exposed port. If you set it to `33069`, connect via `127.0.0.1:33069` in Navicat. Inside the container MySQL still runs on 3306. Containers are isolated from the host and from each other, so this is fine.

## gvb-web and gvb-server fail to start on Windows

You may see an error similar to:

```
'：No such file or directory
```

On Windows, run the following first to avoid line-ending issues, or download the ZIP instead of cloning with Git. Linux/macOS do not need this.

Reason: The project uses LF line endings. Windows Git may convert to CRLF on checkout, which breaks the build.

```bash
git config --global core.autocrlf false
```

## Changed DB password in .env after a previous run

Delete old data under `start/gvb` and re-run.