<p align="center">
<a href="https://szluyu99.github.io/gin-vue-blog/">
<img src="./images/头像.jpeg" width="140" height="140" alt="gin-vue-blog" style="border-radius: 50%">
</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-blue"/>
  <img src="https://img.shields.io/badge/Gin-v1.12-blue"/>
  <img src="https://img.shields.io/badge/GORM-v1.31-blue"/>
  <img src="https://img.shields.io/badge/Jaspr-v0.23-orange"/>
  <img src="https://img.shields.io/badge/Dart-v3.10-blue"/>
  <img src="https://img.shields.io/badge/Vue-v3.5-green"/>
</p>

Dự án blog tách biệt frontend và backend: backend Go + Gin, giao diện blog dùng Jaspr và giao diện quản trị dùng Vue 3.

- GitHub: [szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)
- Gitee: [szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)
- Nhóm trao đổi QQ: 777260310

Hoan nghênh bạn Star và gửi PR.

## Bản demo trực tuyến

Phiên bản Vue cũ dùng dữ liệu giả lập hoàn toàn ở frontend (chỉ giữ lại để tham khảo khi chuyển đổi, cần kích hoạt GitHub Actions thủ công để triển khai):

- Giao diện blog: <https://szluyu99.github.io/gin-vue-blog/>
- Giao diện quản trị: <https://szluyu99.github.io/gin-vue-blog/admin/> (có thể đăng nhập bằng tài khoản và mật khẩu bất kỳ)
- Tài liệu API: <https://szluyu99.github.io/gin-vue-blog/api-docs/> (được tạo từ chú thích Swagger; bản lưu trữ tĩnh không hỗ trợ gọi API trực tuyến)

> Dữ liệu demo là dữ liệu giả được tích hợp trong frontend (`src/mock`). Các thay đổi chỉ tồn tại trên trang hiện tại và sẽ được đặt lại khi tải lại trang.

![Trang chủ blog](./images/前台首页.png)

![Danh sách bài viết trên blog](./images/前台文章列表.png)

![Danh sách bài viết trong trang quản trị](./images/后台文章列表.png)

## Tính năng

Giao diện blog (`irohasu_jaspr`):

- Giao diện tham khảo chủ đề Butterfly của Hexo, tự điều chỉnh để phù hợp với thiết bị di động
- Chế độ tối: lần truy cập đầu tiên theo tùy chọn của hệ thống, cho phép chuyển đổi thủ công và ghi nhớ lựa chọn, không nhấp nháy khi tải lại trang
- Trang chi tiết bài viết hỗ trợ mục lục có liên kết đến từng phần và bài viết đề xuất
- Trạng thái ngắn: nội dung ngắn gọn, hỗ trợ ghim, công khai/riêng tư và bình luận; phần đầu trang chủ luân phiên hiển thị các trạng thái mới nhất
- Trải nghiệm đọc: thanh tiến độ đọc ở đầu trang, sao chép khối mã bằng một lần nhấp, cảnh báo nội dung có thể đã lỗi thời đối với bài viết hơn 90 ngày
- Phân trang danh sách bài viết (áp dụng cả khi lọc theo danh mục/nhãn); danh sách bài viết ở trang chủ tải thêm khi cuộn và tải trước trang tiếp theo
- Trang lưu trữ hiển thị dòng thời gian theo năm, tháng; kích thước chữ trong đám mây nhãn thay đổi theo số bài viết
- Bình luận và trả lời, tường lời nhắn chạy ngang màn hình
- Thông báo trong trang: biểu tượng chuông ở đầu trang báo khi bình luận được trả lời hoặc bài viết của mình có bình luận, kèm trang thông báo riêng; nhấp thông báo sẽ chuyển thẳng đến bình luận tương ứng trong bài viết và tô nổi bật (nếu bật duyệt bình luận, thông báo chỉ được gửi khi quản trị viên nhấn «Duyệt»)
- Lượt thích và thống kê khách truy cập (cộng dồn trong Redis, sao lưu vào cơ sở dữ liệu mỗi 10 phút và trước khi thoát, khôi phục các khóa còn thiếu khi khởi động)
- Đăng ký người dùng: mặc định tạo tài khoản trực tiếp; có thể bật đăng ký bằng mã xác minh email (`Captcha.SendEmail` trong `config.yml`)

Giao diện quản trị (`gin-blog-admin`):

- Xác thực JWT và phân quyền theo vai trò; quyền truy cập menu và API đều có thể cấu hình động trong trang quản trị
- Bảng tổng quan: bình luận/lời nhắn chờ duyệt, bài viết mới nhất, phân bố danh mục, khu vực của khách truy cập, các lần đăng nhập gần đây và xu hướng truy cập trong 14 ngày; gọi API song song, một API lỗi không làm hỏng toàn bộ trang
- Menu frontend do backend cung cấp (định tuyến động)
- Soạn thảo bài viết bằng Markdown, hỗ trợ nhập `.md` / `.markdown` và xuất `.md`
- Nhật ký thao tác, nhật ký đăng nhập (bao gồm các lần thất bại), theo dõi người dùng trực tuyến và buộc đăng xuất
- Nhật ký lỗi frontend: các ngoại lệ chưa được xử lý ở cả giao diện blog và quản trị tự động được gửi về danh sách lỗi trong trang quản trị; các lỗi cùng loại được gộp theo dấu vân tay và chỉ tăng số lần xuất hiện (không dùng Sentry)
- Có thể vô hiệu hóa người dùng và vai trò: sau khi vô hiệu hóa, đăng nhập bị từ chối, token đã cấp lập tức mất hiệu lực và các quyền chỉ được cấp qua vai trò đó bị thu hồi
- Tải tệp lên bộ nhớ cục bộ hoặc Qiniu Cloud
- Các thao tác CRUD được đóng gói thành Hook dùng chung

## Công nghệ sử dụng

Backend: Go / Gin / GORM / SQLite (mặc định, có thể đổi sang MySQL) / Redis / Viper / `log/slog`

Giao diện blog: Dart / Jaspr / Jaspr Router / Dio, sử dụng kết xuất phía máy chủ; giao diện quản trị: Vue 3 / Vite / Pinia / Naive UI / Axios

Khác: triển khai bằng một lệnh với Docker Compose, Nginx phục vụ tài nguyên tĩnh và làm proxy ngược, lưu trữ đối tượng Qiniu Cloud, xác minh người dùng của Tencent Cloud (tùy chọn, mặc định tắt)

## Cấu trúc thư mục

```bash
gin-vue-blog
├── gin-blog-admin      # Frontend quản trị blog
├── irohasu_jaspr       # Giao diện blog dùng Jaspr
├── gin-blog-front      # Giao diện blog Vue cũ (giữ lại để tham khảo khi chuyển đổi)
├── gin-blog-server     # Backend của blog
└── deploy              # Triển khai bằng Docker
```

Backend:

```bash
gin-blog-server
├── cmd                 # Điểm khởi chạy chương trình, script khởi tạo dữ liệu
├── internal
│   ├── handle          # Xử lý API
│   ├── model           # Mô hình dữ liệu và thao tác cơ sở dữ liệu
│   ├── middleware      # Middleware
│   ├── global          # Cấu hình toàn cục, mã lỗi
│   └── utils           # Các hàm tiện ích
├── docs                # Tài liệu Swagger
├── assets              # Tệp tài nguyên
├── config.yml          # Tệp cấu hình
└── Dockerfile
```

Frontend (hai dự án Vue có cấu trúc tương tự, giao diện quản trị có thêm `layout` và `composables`):

```bash
├── src
│   ├── api.js          # API
│   ├── assets          # Tài nguyên tĩnh (chỉ có ở giao diện quản trị)
│   ├── components      # Các thành phần giao diện
│   ├── mock            # Dữ liệu giả lập
│   ├── router          # Định tuyến (giao diện blog dùng một tệp router.js)
│   ├── store           # Quản lý trạng thái
│   ├── utils           # Các hàm tiện ích
│   └── views           # Các trang
├── .env.*              # Biến môi trường (development / production / mock)
├── uno.config.js       # Cấu hình UnoCSS
└── vite.config.js      # Cấu hình Vite
```

## Bắt đầu nhanh

**Xem hướng dẫn phát triển cục bộ tại [quick_start.md](./quick_start.md)**, bao gồm hai cách chạy: chế độ giả lập (không khởi động backend) và khởi động đầy đủ, cùng địa chỉ truy cập, tài khoản mặc định và các vấn đề thường gặp.

Hướng dẫn tiếng Việt đầy đủ về chạy local, deploy Docker/VPS, HTTPS, backup/restore, cập nhật và xử lý sự cố: [HUONG_DAN_CHAY_VA_DEPLOY.md](./HUONG_DAN_CHAY_VA_DEPLOY.md).

Dữ liệu khởi tạo chỉ gồm dữ liệu hệ thống cơ bản (menu / API / vai trò / cấu hình), không có bài viết. Để tạo nội dung dùng thử (danh mục, nhãn, 15 bài viết, bình luận và trả lời, lời nhắn, liên kết bạn bè), dùng `./dev.sh fresh --demo` hoặc `cd gin-blog-server/cmd/generate-data && go run main.go -t demo`.

Nếu chỉ muốn xem thử, có thể chạy bằng một lệnh với Docker Compose (cần Docker + Docker Compose; trên Windows hãy dùng GitBash):

```bash
git clone https://github.com/szluyu99/gin-vue-blog
cd gin-vue-blog/deploy
./bootstrap.sh
```

Giao diện blog: [localhost](http://localhost/), giao diện quản trị: [localhost/admin](http://localhost/admin), tài khoản mặc định: `admin / 123456`.

Xem tài liệu triển khai chi tiết tại [deploy/README.md](./deploy/README.md).

> Trên Windows, nên chạy `git config --global core.autocrlf false` trước khi clone. Dự án sử dụng ký tự xuống dòng LF; CRLF có thể gây lỗi khi dựng Docker.
>
> Nếu bật đăng ký xác minh email khi triển khai bằng Docker (`Captcha.SendEmail: true`), cần thay localhost trong hàm `GetEmailVerifyURL` tại `gin-blog-server/internal/utils/email.go` bằng tên miền của bạn.

## Kiểm thử và CI

```bash
cd gin-blog-server && go test ./...   # Backend: model / handle / middleware / global / utils / đăng ký tuyến, 8 gói
cd irohasu_jaspr && flutter test      # Giao diện blog Jaspr
cd gin-blog-admin  && pnpm test       # Giao diện quản trị: vitest, 27 tệp, 165 ca kiểm thử
```

Kiểm thử frontend chủ yếu là kiểm thử hồi quy: mỗi ca đều đã được xác nhận sẽ thất bại nếu hoàn tác bản sửa lỗi tương ứng. Các lỗi cụ thể được kiểm thử nằm trong mục «Kiểm thử thành phần» của `code_audit.md`.

`.github/workflows/ci.yml` chạy bốn nhóm tác vụ khi push và tạo PR:

- Server: `gofmt` / `go vet` / `go test`, đồng thời kiểm tra tính nhất quán giữa `docs/` và chú thích Swagger
- Frontend: giao diện Jaspr chạy `dart analyze` / `flutter test` / `jaspr build`; giao diện quản trị Vue chạy `pnpm lint` / `pnpm test` / `pnpm build`
- Docker: dựng image web và server, sau đó khởi động để kiểm tra tình trạng hoạt động. Image web sử dụng trực tiếp sản phẩm đầu ra của tác vụ Frontend cùng `deploy/build/web/Dockerfile` (không dựng lại frontend trong container); image server dùng buildx và bộ nhớ đệm GitHub Actions. Khi nhánh main có commit mới, sau khi vượt qua kiểm thử nhanh, hai image được đẩy lên GHCR với các thẻ `latest` và `main-<SHA rút gọn>`
- Deploy: khởi động toàn bộ hệ thống Compose, kiểm tra các điều kiện bất biến của dữ liệu phân quyền ban đầu (admin có thể sửa cấu hình, guest không thể, người chưa đăng nhập bị chặn, chạy lại không tạo thêm tài nguyên)

## Lộ trình phát triển

Kế hoạch đầy đủ và lý do không triển khai một số hạng mục nằm trong [roadmap.md](./roadmap.md). Khi tiếp quản công việc, hãy đọc mục «Trạng thái hiện tại» trong đó trước. Dưới đây là phần tổng quan.

Tính năng:

- Cần thiết kế lại các khối thông tin ở thanh bên của giao diện blog (`AuthorInfo` / `Announcement` / `WebsiteInfo`): toàn bộ cột bị ẩn trên màn hình hẹp, khiến người dùng điện thoại không thấy thông báo và thông tin website
- Siêu dữ liệu khi chia sẻ bài viết: mới chỉ thiết lập `document.title`, chưa có `meta description` và `og:*`
- Trang «Album ảnh» hiện chỉ là trang giữ chỗ «Truy cập bị cấm» (`views/entertainment/album`), cả hai liên kết trong menu điều hướng vẫn còn
- Hỗ trợ hình ảnh cho trạng thái ngắn (phiên bản đầu chỉ có văn bản và bình luận; việc bổ sung hình ảnh cần giải quyết cả lưu trữ và ảnh thu nhỏ)
- Hiện frontend phải tải lần lượt từng trang để tìm bình luận đích từ thông báo; khi có nhiều bình luận, backend cần bổ sung API tính số trang theo ID bình luận

Kỹ thuật:

- Liên kết đăng ký chứa mật khẩu dạng văn bản thuần, cùng các vấn đề về nhật ký và TLS của mô-đun email (S6, S7 trong `code_audit.md`) — chỉ phát sinh khi bật đăng ký qua email; cấu hình hiện tại chưa bật nên tạm hoãn
- Chuyển cách triển khai sang tải image từ GHCR (CI đã phát hành image, nhưng Compose vẫn dùng `build:` cục bộ)
- Tách `gin-blog-front` và `gin-blog-admin` thành các kho mã độc lập
- Hoàn thiện tài liệu API

Các hạng mục xác định không triển khai (xem lý do trong [roadmap.md](./roadmap.md)): tìm kiếm bằng ElasticSearch, SSR / SSG, đăng nhập qua bên thứ ba, quốc tế hóa, RSS / sitemap, xoay vòng nhật ký backend, thay `@iconify/vue` bằng các lớp biểu tượng UnoCSS.
