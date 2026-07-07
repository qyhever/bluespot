# bluespot-backend Agent 指南

## 项目结构

- 这是一个 Go + Gin 后端服务。
- 可执行程序入口是 `cmd/main.go`。
- 业务代码应保留在 `internal/` 中，并保持现有分层：controller -> service -> repository/persistence。

## 常用命令

- 安装并整理依赖：`make deps`
- 运行测试：`make test`
- 构建二进制文件：`make build`
- 构建并运行一次：`make dev`
- 使用热重载运行：`make hot`
- 清理构建产物：`make clean`

## 架构规则

- 在 `internal/api/router.go` 中添加 HTTP 路由。
- Controller 应保持轻量：解析请求、调用 service、写入响应。
- 业务逻辑放在 `internal/service/`。
- 数据访问放在 `internal/repository/` 和 `internal/repository/persistence/`。
- 优先在路由 setup 中按现有 controller 和 service 的构造方式完成依赖装配。

## 响应约定

- API 响应复用 `internal/controller/response.go`。
- 成功响应应通过 `ResponseSuccess` 返回。
- 业务失败应使用 `ResponseFailed` 或 `ResponseFailedWithMsg`。
- 除非任务明确要求新的响应契约，否则不要引入临时的 JSON 响应结构。

## 配置与运行说明

- 配置在 `internal/config` 中初始化，且早于 logger 和 router 启动。
- 修改启动行为前，先检查 `internal/config/app.yml`、`dev.yml`、`dev.local.yml` 和 `prod.yml`。
- 服务还会通过 `/public` URL 暴露 `/public` 下的静态文件；编辑 router setup 时不要破坏该路径。
- `rest/index.http` 是查看或扩展示例 API 请求最快的位置。

## Swagger 文档

- `internal/api/router.go` 已注册 Swagger UI：`/api/swagger/index.html`。
- `internal/api/router.go` 还注册了 knife4go UI：`/api/k4/index.html`。
- Swagger 基础信息维护在 `cmd/main.go` 的 swag 注释中，当前 `BasePath` 为 `/api`。
- Swagger 生成产物在 `docs/`：`docs.go`、`swagger.json`、`swagger.yaml`。
- 重新生成文档可执行：`go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g ./cmd/main.go -o ./docs`。
- `scripts/deploy.sh` 会在部署前执行同一条 Swagger 生成命令，并同步 `docs/` 到服务器。
- 不要删除 `internal/api/router.go` 中的 `_ "bluespot/docs"`，否则 Swagger 规格不会注册到运行时。

## 工作约定

- 优先在现有分层边界内做小改动，不要绕过 service 直接从 controller 调用 repository。
- 导出名称保持符合 Go 习惯，避免不必要的抽象。
- 添加 endpoint 时，应同步更新 router、controller、service 和 repository 层。
- 如果变更影响示例 payload 或本地 stub 数据，请检查 `public/meta.json`。

## 参考文件

- 运行入口：`cmd/main.go`
- Router setup：`internal/api/router.go`
- 共享 API 响应工具：`internal/controller/response.go`
- 配置文件：`internal/config/`
- 示例 HTTP 请求：`rest/index.http`
