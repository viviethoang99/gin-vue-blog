# Hướng dẫn chạy và triển khai gin-vue-blog

Tài liệu này áp dụng cho cấu trúc hiện tại của repository, gồm ba ứng dụng:

| Thành phần | Thư mục | Công nghệ | Cổng local |
| --- | --- | --- | --- |
| Backend API | `gin-blog-server` | Go 1.26, Gin, GORM | `8765` |
| Blog frontend | `gin-blog-front` | Vue 3, Vite | `8888` |
| Admin frontend | `gin-blog-admin` | Vue 3, Vite, Naive UI | `8889` |

Các địa chỉ mặc định khi chạy local:

- Blog: <http://localhost:8888>
- Admin: <http://localhost:8889>
- API: <http://localhost:8765/api>
- Swagger: <http://localhost:8765/swagger/index.html>
- Tài khoản quản trị: `admin / 123456`

## 1. Nên chọn cách nào?

| Nhu cầu | Cách nên dùng |
| --- | --- |
| Phát triển trên Linux/WSL | `./dev.sh start --demo` |
| Phát triển trên macOS hoặc máy không có `setsid` | Chạy thủ công theo mục 4 |
| Chỉ xem giao diện, không chạy backend | Frontend Mock mode |
| Muốn hiểu và chạy từng thành phần | Chạy thủ công |
| Deploy lên VPS/server | `cd deploy && ./bootstrap.sh` |

## 2. Chạy local nhanh nhất

### 2.1. Yêu cầu

- Git
- Go `1.26+`
- Node.js `>=20.19`; Node 22 khuyến nghị
- pnpm (`11` là phiên bản dùng trong Dockerfile)
- Redis 7, hoặc Docker để script tự chạy Redis
- Lệnh `setsid` nếu muốn dùng `dev.sh`

Cài pnpm bằng một trong hai cách:

```bash
corepack enable
corepack install --global pnpm@11
```

Hoặc:

```bash
npm install -g pnpm@11
```

Kiểm tra môi trường:

```bash
go version
node --version
pnpm --version
```

### 2.2. Khởi động toàn bộ hệ thống

`dev.sh` hiện dùng `setsid` và một số cơ chế quản lý process kiểu Linux. Trước tiên kiểm tra:

```bash
command -v setsid
```

Nếu lệnh không có output (macOS mặc định thường không có), dùng cách chạy thủ công ở mục 4. Docker Compose ở mục 6 vẫn chạy bình thường trên Docker Desktop.

Từ thư mục gốc repository:

```bash
./dev.sh start --demo
```

Lệnh này sẽ:

1. Dùng Redis đang chạy ở `127.0.0.1:6379`; nếu chưa có thì ưu tiên `redis-server`, sau đó mới dùng Docker.
2. Dùng SQLite tại `gin-blog-server/cmd/gvb.db`, không cần MySQL.
3. Tự tạo bảng và dữ liệu hệ thống khi chưa có database.
4. Nạp thêm dữ liệu demo do có cờ `--demo`.
5. Build và chạy backend.
6. Chạy blog frontend và admin frontend với API thật, tự tắt Mock mode.

Nếu không cần bài viết và bình luận demo:

```bash
./dev.sh start
```

Các lệnh quản lý:

```bash
./dev.sh status
./dev.sh logs server
./dev.sh logs front
./dev.sh logs admin
./dev.sh restart
./dev.sh stop
```

Nếu PID cũ bị mất nhưng cổng vẫn bị chiếm:

```bash
./dev.sh stop --force
```

### 2.3. Tạo lại database local

```bash
./dev.sh fresh --demo
```

`fresh` không xóa ngay database cũ mà đổi tên thành `gin-blog-server/cmd/gvb.db.bak`. Nó cũng xóa Redis DB 7 nếu có thể; các file upload không bị xóa.

Muốn khôi phục database cũ:

```bash
./dev.sh stop
mv gin-blog-server/cmd/gvb.db gin-blog-server/cmd/gvb.db.new
mv gin-blog-server/cmd/gvb.db.bak gin-blog-server/cmd/gvb.db
./dev.sh start
```

## 3. Chạy frontend bằng Mock mode

Cách này không cần Go, Redis hay database. Dữ liệu nằm trong `src/mock`; các thay đổi sẽ mất khi reload trang.

Blog frontend:

```bash
cd gin-blog-front
pnpm install
pnpm dev
```

Admin frontend:

```bash
cd gin-blog-admin
pnpm install
pnpm dev
```

Hai file `.env.development` hiện đặt `VITE_USE_MOCK=true`, nên chạy trực tiếp `pnpm dev` sẽ dùng Mock mode. Admin Mock mode chấp nhận bất kỳ tài khoản/mật khẩu nào.

## 4. Chạy thủ công với backend thật

### 4.1. Redis

Dùng Redis local:

```bash
redis-server --port 6379
```

Hoặc Docker:

```bash
docker run -d --name gvb-local-redis -p 6379:6379 redis:7.0-alpine
```

Backend local dùng Redis DB `7` theo `gin-blog-server/config.yml`.

### 4.2. Khởi tạo dữ liệu

```bash
cd gin-blog-server
go mod download
cd cmd
sh generate_data.sh
```

Script tạo bảng SQLite và dữ liệu nền: cấu hình website, trang, menu, resource, role, tài khoản `admin` và `guest`. Script có thể chạy lại; các bản ghi đã có sẽ được bỏ qua.

Muốn nạp thêm dữ liệu demo:

```bash
cd gin-blog-server/cmd/generate-data
go run main.go -t demo
```

### 4.3. Chạy backend

Backend phải chạy với working directory là `gin-blog-server/cmd`, vì đường dẫn config, SQLite và upload đang là đường dẫn tương đối:

```bash
cd gin-blog-server/cmd
go run main.go
```

### 4.4. Chạy hai frontend với API thật

Trong cả hai file sau, đổi `VITE_USE_MOCK` thành `false`:

- `gin-blog-front/.env.development`
- `gin-blog-admin/.env.development`

Sau đó chạy ở hai terminal:

```bash
cd gin-blog-front
pnpm install
pnpm dev
```

```bash
cd gin-blog-admin
pnpm install
pnpm dev
```

Vite proxy `/api` và `/public` về `http://localhost:8765`, nên local không cần cấu hình CORS bổ sung. Vite chỉ đọc `.env*` lúc khởi động; thay đổi biến môi trường xong phải restart `pnpm dev`.

## 5. Cấu hình backend

Local dùng `gin-blog-server/config.yml`; Docker dùng `gin-blog-server/config.docker.yml`.

| Mục | Ý nghĩa |
| --- | --- |
| `Server.Mode` | `debug` cho local, `release` cho production |
| `Server.DbType` | `sqlite` hoặc `mysql` |
| `Server.DbAutoMigrate` | Tự tạo/cập nhật schema khi chạy |
| `Server.trusted-proxies` | Các proxy được phép gửi IP thật |
| `Server.allowed-origins` | Danh sách origin được CORS cho phép |
| `JWT.Secret` | Khóa ký JWT; production không được dùng giá trị mẫu |
| `Session.Salt` | Khóa session; production không được dùng giá trị mẫu |
| `Session.Secure` | Nên đặt `true` khi chạy HTTPS |
| `Redis.DB` | Local mặc định dùng DB 7 |
| `Upload.OssType` | `local` hoặc `qiniu` |
| `Captcha.SendEmail` | Bật/tắt đăng ký qua email xác thực |

Các biến môi trường backend được hỗ trợ:

```text
SERVER_PORT
SERVER_DBTYPE
JWT_SECRET
SESSION_SALT
SQLITE_DSN
MYSQL_HOST
MYSQL_PORT
MYSQL_DBNAME
MYSQL_USERNAME
MYSQL_PASSWORD
REDIS_ADDR
REDIS_PASSWORD
```

Ở `release` mode, backend sẽ từ chối chạy nếu `JWT_SECRET` hoặc `SESSION_SALT` rỗng/vẫn là giá trị mẫu.

## 6. Deploy bằng Docker Compose

Đây là cách nên dùng trên VPS. Stack production gồm Nginx, Go backend, MySQL 8 và Redis 7.

### 6.1. Yêu cầu server

- Linux VPS khuyến nghị; macOS/Windows có thể dùng Docker Desktop để chạy thử
- Docker Engine và Docker Compose
- Tối thiểu khoảng 2 GB RAM; lần build frontend cần nhiều RAM hơn lúc chạy
- Mở cổng `80` và `443` nếu dùng domain công khai

Kiểm tra:

```bash
docker --version
docker compose version
```

### 6.2. Clone và chuẩn bị cấu hình

```bash
git clone https://github.com/viviethoang99/gin-vue-blog.git
cd gin-vue-blog/deploy
```

Sửa `start/.env` trước khi chạy lần đầu. Ít nhất cần đổi:

```dotenv
REDIS_PASSWORD=<mat-khau-redis-manh>
MYSQL_ROOT_PASSWORD=<mat-khau-mysql-manh>
SERVER_NAME=blog.example.com
```

Các giá trị đáng chú ý khác:

```dotenv
DATA_DIRECTORY=./gvb
REDIS_PORT=63799
MYSQL_PORT=33066
SERVER_PORT=8765
USE_HTTPS=false
```

`REDIS_PORT`, `MYSQL_PORT` và `SERVER_PORT` là cổng expose ra host. Các container vẫn giao tiếp bằng cổng nội bộ cố định.

### 6.3. Deploy HTTP lần đầu

Tiếp tục từ thư mục `gin-vue-blog/deploy` của bước trên:

```bash
./bootstrap.sh
```

Nếu file chưa có quyền execute:

```bash
chmod +x bootstrap.sh clean_docker.sh build_web.sh
./bootstrap.sh
```

`bootstrap.sh` sẽ:

1. Dừng và xóa các container cũ của riêng dự án.
2. Tạo `deploy/start/.env.secrets` chứa `JWT_SECRET` và `SESSION_SALT` ngẫu nhiên.
3. Build hai frontend và backend trong Docker.
4. Khởi động MySQL, Redis, backend và Nginx.
5. Chờ MySQL/Redis healthy trước khi chạy backend.
6. Lần đầu, MySQL import dữ liệu mẫu từ `deploy/build/mysql/gvb.sql`.
7. Chạy seed dữ liệu quyền/menu/resource mỗi lần backend khởi động.

`bootstrap.sh` không phải zero-downtime deployment: container cũ bị dừng trước khi image mới build xong. Nếu uptime là yêu cầu bắt buộc, cần thêm registry và chiến lược rolling/blue-green bên ngoài script hiện tại.

Sau khi hoàn tất:

- Blog: `http://<server-ip>/`
- Admin: `http://<server-ip>/admin/`
- API: `http://<server-ip>/api`
- Tài khoản: `admin / 123456`

Đăng nhập admin và đổi mật khẩu mặc định ngay sau lần chạy đầu tiên.

### 6.4. Xem trạng thái và log

```bash
cd deploy/start
docker compose ps
docker compose logs -f gvb-server
docker compose logs -f gvb-web
docker compose logs -f gvb-mysql
docker compose logs -f gvb-redis
```

Kiểm tra nhanh:

```bash
curl -I http://localhost/
curl -I http://localhost/admin/
```

Swagger được backend phục vụ nhưng Nginx hiện chỉ proxy `/api`, vì vậy truy cập trực tiếp:

```text
http://<server-ip>:8765/swagger/index.html
```

Nếu không muốn expose backend ra Internet, chỉ dùng Swagger qua SSH tunnel hoặc chỉnh Nginx thêm route `/swagger`.

### 6.5. Deploy HTTPS bằng certificate có sẵn

1. Trỏ DNS A/AAAA của domain về server.
2. Đặt certificate và private key hợp lệ tại:

```text
deploy/start/server.crt
deploy/start/server.key
```

3. Sửa `deploy/start/.env`:

```dotenv
USE_HTTPS=true
SERVER_NAME=blog.example.com
```

4. Sửa `gin-blog-server/config.docker.yml`:

```yaml
Session:
  Secure: true
```

5. Build và chạy lại:

```bash
cd /duong-dan/toi/gin-vue-blog/deploy
./bootstrap.sh
```

Nginx sẽ redirect HTTP sang HTTPS. Không dùng certificate mẫu cho production; phải thay bằng certificate đúng domain.

Nếu dùng Caddy, Traefik hoặc Nginx Proxy Manager ở ngoài stack để tự cấp Let's Encrypt, có thể giữ `USE_HTTPS=false` và để proxy ngoài terminate TLS. Khi proxy không nằm trong mạng private mặc định, bổ sung IP/CIDR của proxy vào `Server.trusted-proxies`.

### 6.6. Build frontend bằng pnpm trên host

Mặc định frontend được build trong Docker. Nếu server đã có Node và pnpm:

```bash
cd /duong-dan/toi/gin-vue-blog/deploy
./bootstrap.sh dev
```

Lệnh này tạo `deploy/build/web/dist_blog` và `deploy/build/web/dist_admin`, sau đó Nginx image copy static assets đã build sẵn.

## 7. Cập nhật phiên bản đang deploy

Backup trước, sau đó:

```bash
cd /duong-dan/toi/gin-vue-blog
git pull
cd deploy
./bootstrap.sh
```

`clean_docker.sh` chỉ dừng và xóa bốn container của dự án; không xóa dữ liệu bind mount trong `deploy/start/gvb`.

Chỉ restart container, không build source:

```bash
cd deploy/start
docker compose restart
```

Build lại image từ source hiện có:

```bash
cd deploy/start
docker compose up -d --build
```

Lệnh này yêu cầu `deploy/start/.env.secrets` đã được tạo bởi lần chạy `bootstrap.sh` trước đó.

## 8. Dữ liệu, backup và restore

Dữ liệu production mặc định nằm dưới:

```text
deploy/start/gvb/data/mysql     # MySQL
deploy/start/gvb/data/redis     # Redis AOF
deploy/start/gvb/file/uploaded  # File upload local
```

`deploy/start/.env.secrets` cũng phải được giữ an toàn. Nếu đổi `JWT_SECRET` hoặc `SESSION_SALT`, toàn bộ token/session hiện có sẽ mất hiệu lực.

### 8.1. Backup toàn bộ, có downtime ngắn

```bash
cd deploy/start
docker compose down
tar -czf "gvb-backup-$(date +%Y%m%d-%H%M%S).tar.gz" gvb .env .env.secrets server.crt server.key
docker compose up -d
```

Chép file backup sang máy/storage khác; không để bản backup duy nhất trên cùng VPS.

### 8.2. Restore

```bash
cd deploy/start
docker compose down
mv gvb gvb.before-restore
tar -xzf gvb-backup-YYYYMMDD-HHMMSS.tar.gz
docker compose up -d --build
```

Chỉ xóa `gvb.before-restore` sau khi đã kiểm tra dữ liệu khôi phục thành công.

### 8.3. Đổi mật khẩu MySQL sau khi đã chạy

Chỉ sửa `MYSQL_ROOT_PASSWORD` trong `.env` không làm mật khẩu bên trong database cũ tự đổi. Với môi trường thử nghiệm không cần dữ liệu: backup, `docker compose down`, đổi tên thư mục `gvb`, rồi chạy lại `../bootstrap.sh` để MySQL khởi tạo mới.

Production có dữ liệu thật thì phải đổi password trong MySQL trước, sau đó mới cập nhật `.env`.

## 9. Lưu ý bảo mật production

1. Đổi ngay `MYSQL_ROOT_PASSWORD`, `REDIS_PASSWORD` và mật khẩu admin mặc định.
2. Giữ kín `deploy/start/.env.secrets`; file này đã được `.gitignore`.
3. Bật HTTPS và đặt `Session.Secure: true`.
4. Compose hiện publish MySQL, Redis và backend ra host. Dùng firewall chỉ cho phép cổng `80/443` từ Internet. Tốt hơn nữa, bind các cổng quản trị vào `127.0.0.1` hoặc bỏ `ports` nếu không cần truy cập từ host.
5. Không mở cổng MySQL `33066`, Redis `63799` hay backend `8765` công khai.
6. Backup định kỳ cả MySQL, Redis, uploaded files và `.env.secrets`.
7. Không commit SMTP password, Qiniu secret, certificate private key hoặc file `.env.secrets`.
8. Nếu frontend/backend khác domain, cấu hình chính xác `Server.allowed-origins`; không dùng wildcard tùy tiện.

## 10. Email xác thực và upload

### 10.1. Email xác thực

Mặc định `Captcha.SendEmail: false`: đăng ký tạo tài khoản ngay và không cần SMTP.

Muốn bật email:

1. Điền `Email.Host`, `Email.Port`, `Email.From`, `Email.SmtpPass`, `Email.SmtpUser`.
2. Đặt `Captcha.SendEmail: true`.
3. Sửa `GetEmailVerifyURL` trong `gin-blog-server/internal/utils/email.go` để link xác thực dùng domain HTTPS thật. Code hiện tại tự tạo link `localhost:<port>`.
4. Build lại backend.

### 10.2. Upload local

Docker mount upload vào `deploy/start/gvb/file/uploaded`. Nginx proxy `/public/uploaded` về backend, nên frontend và API dùng cùng domain.

### 10.3. Qiniu

Đổi `Upload.OssType` thành `qiniu`, rồi điền đầy đủ phần `Qiniu` trong config tương ứng. Không commit `AccessKey` và `SecretKey`; nên dùng config production riêng hoặc secret management của hạ tầng.

## 11. Test trước khi deploy

```bash
cd gin-blog-server
go test ./...

cd ../gin-blog-front
pnpm install --frozen-lockfile
pnpm test
pnpm build

cd ../gin-blog-admin
pnpm install --frozen-lockfile
pnpm test
pnpm build
```

## 12. Xử lý sự cố

### `pnpm` không tồn tại

```bash
corepack enable
corepack install --global pnpm@11
```

Hoặc dùng `corepack pnpm <command>` thay cho `pnpm <command>`.

### Frontend vẫn dùng Mock/API cũ

- Đặt `VITE_USE_MOCK=false` khi chạy với backend thật.
- Restart Vite sau khi sửa `.env`.
- Kiểm tra `.env.development` của cả hai frontend.

### Backend không chạy vì Redis

```bash
redis-cli -p 6379 ping
./dev.sh logs server
```

Local backend cần Redis ngay cả khi dùng SQLite.

### Cổng bị chiếm

```bash
./dev.sh status
lsof -iTCP:8765 -sTCP:LISTEN
lsof -iTCP:8888 -sTCP:LISTEN
lsof -iTCP:8889 -sTCP:LISTEN
```

Dừng tiến trình cũ hoặc dùng `./dev.sh stop --force` nếu đó là tiến trình của lần chạy trước.

### Docker MySQL không healthy

```bash
cd deploy/start
docker compose ps
docker compose logs gvb-mysql
```

Lần chạy đầu phải import SQL nên có thể lâu. Nếu đổi password sau khi database đã được tạo, xem mục backup/đổi password.

### Container backend restart liên tục

```bash
docker compose logs --tail=200 gvb-server
```

Nguyên nhân thường gặp:

- Thiếu `start/.env.secrets` vì chạy `docker compose up` trước `bootstrap.sh`.
- MySQL/Redis chưa healthy hoặc password không khớp.
- `config.docker.yml` ở release mode nhưng secret vẫn rỗng/giá trị mẫu.

### Upload báo `413 Request Entity Too Large`

Nginx hiện đặt `client_max_body_size 12M`, còn backend giới hạn file khoảng 10 MB. Muốn tăng phải sửa đồng bộ Nginx và giới hạn backend, sau đó build lại.

### Sửa database nhưng API vẫn trả dữ liệu cũ

`page` và `config` được cache trong Redis khoảng 10 phút. Local có thể xóa cache bằng:

```bash
redis-cli -n 7 del page config
```

Docker có password:

```bash
docker exec -it gvb-redis redis-cli -a '<REDIS_PASSWORD>' -n 7 del page config
```

### Windows báo `No such file or directory`

Repository dùng LF. Trước khi clone:

```bash
git config --global core.autocrlf false
```

Clone lại hoặc chuyển shell script về LF, rồi chạy bằng Git Bash/WSL.

## 13. Checklist production

- [ ] DNS đã trỏ về VPS.
- [ ] Chỉ cổng 80/443 được mở công khai.
- [ ] Đã đổi MySQL/Redis password.
- [ ] `.env.secrets` đã sinh và được backup kín.
- [ ] Đã bật HTTPS và `Session.Secure: true`.
- [ ] Certificate đúng domain và còn hạn.
- [ ] Đã đổi mật khẩu `admin / 123456`.
- [ ] `docker compose ps` cho thấy các service hoạt động.
- [ ] Blog, admin, API và upload đều truy cập được.
- [ ] Backup/restore đã được thử ít nhất một lần.
- [ ] Có lịch backup định kỳ sang máy/storage khác.
