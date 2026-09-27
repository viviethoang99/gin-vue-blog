# Development Guide

The MySQL image mainly serves to **initialize database data**, with the core being the `gvb.sql` file.

This SQL file runs automatically when the MySQL container starts.

To change initial data, edit `gvb.sql`.

> If the project has already run once, delete existing data under `start/gvb` first (remember to back up).

Then rerun the one-click script: `./bootstrap.sh`

## 注意：权限数据不在这里

`gvb.sql` 只负责**示例内容**（文章、分类、标签、用户、博客配置、页面）。

角色、接口资源、菜单以及它们的关联关系（`role`、`resource`、`menu`、`role_resource`、
`role_menu`、`user_auth_role`）**以代码为准**，由后端镜像启动时执行的
`cmd/generate-data` 生成，资源列表定义在 `internal/model/seed_resource.go`。

这么做是因为接口会随代码变动：手写的 SQL 一旦落后，新接口就不在资源表里，
而 `middleware.JWTAuth` 对未登记的接口会跳过权限校验，等于任何登录用户都能调用。
新增后台接口时只需要改 `seed_resource.go`（`internal/route_resource_test.go` 会校验
它和路由一一对应），重启后端容器即生效。