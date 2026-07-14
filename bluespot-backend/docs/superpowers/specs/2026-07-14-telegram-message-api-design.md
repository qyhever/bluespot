# Telegram 消息发送接口设计

## 目标

为 bluespot 后端新增一个受 JWT 保护的 Telegram 消息发送接口。调用方只提交消息文本，Bot Token 和目标 Chat ID 统一从服务端配置读取，避免客户端接触 Telegram 凭据或控制发送目标。

## API 契约

- 方法与路径：`POST /api/telegram`
- 认证：沿用现有 `JWTAuthMiddleware`
- Content-Type：`application/json`
- 请求体：`{"text":"hello from bluespot"}`
- `text` 为必填字符串；缺失或为空时返回现有 `CodeInvalidParam` 响应。
- 发送成功时调用 `ResponseSuccess(c, nil)`。
- 配置缺失、网络异常、Telegram 非成功 HTTP 状态或 Telegram 响应 `ok: false` 时，记录服务端日志并返回现有 `CodeServerBusy` 响应，不向客户端泄露 Bot Token 或上游响应细节。

## 配置

在根配置结构中增加 `TG TelegramConfig`，字段如下：

- `bot_token`：Telegram Bot Token。
- `chat_id`：固定目标 Chat ID，使用字符串以兼容群组或频道的负数 ID。

保留 `internal/config/app.yml` 已有的 `tg` 配置段，并绑定以下环境变量：

- `BLUESPOT_TG_BOT_TOKEN`
- `BLUESPOT_TG_CHAT_ID`

## 代码结构

- `internal/model`：定义只包含 `text` 的请求模型，并通过 Gin binding 声明必填。
- `internal/controller`：解析、校验请求，调用 service，并使用统一响应函数返回结果。
- `internal/service`：提供 Telegram 发送业务入口，依赖一个可替换的发送客户端，便于单元测试。
- `internal/pkg/telegram`：封装 `sendMessage` HTTP 调用、请求序列化、超时、状态码检查及 Telegram `ok` 字段检查。
- `internal/api/router.go`：构造 service/controller，注册受 JWT 保护的 `POST /api/telegram`。
- `rest/index.http`：增加可直接调用的请求示例。

## 数据流

请求进入路由后先经过 JWT 中间件。Controller 解析 `text`，Service 使用服务端配置的 Bot Token 和 Chat ID 调用 Telegram 客户端。客户端向 `https://api.telegram.org/bot{token}/sendMessage` 提交 JSON，并检查 HTTP 与业务响应；结果最终转换为项目现有统一响应格式。

## 安全与可靠性

- Bot Token 只参与服务端 URL 构造，不写入日志或 API 响应。
- 使用显式 HTTP 超时，避免上游连接长期占用请求。
- 发送前校验 Token、Chat ID 和消息文本非空。
- 始终关闭 Telegram 响应体。

## 测试与验收

- 配置测试覆盖 `tg` 字段和环境变量绑定。
- Telegram 客户端使用 `httptest.Server` 覆盖正确 JSON、成功响应、非 2xx 响应及 `ok: false`。
- Controller 或路由测试覆盖空文本校验与成功调用。
- 执行 `go test ./...` 和 `go build ./...`。
- Swagger 注释描述接口，并重新生成 `docs` 下的 Swagger 产物。

## 非目标

- 不允许请求覆盖 Chat ID。
- 不支持图片、文件、Markdown 模式或批量发送。
- 不实现重试、队列或发送历史持久化。
