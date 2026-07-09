# 后端分片上传接口方案

## Summary

在 bluespot-backend 实现前端 bluespot-frontend/src/api/chunk.ts 约定的 3 个接口：

- POST /api/upload/verify
- POST /api/upload/chunk
- POST /api/upload/merge

所有接口走 JWTAuthMiddleware，从 Gin context 的 middleware.ContextUserIDKey 获取 userId。
rootMd5 按当前前端字段理解为 fileMd5，uploadId 使用：

MD5(fileMd5 + userId + chunk_dir_salt)

该值同时作为 chunks 目录名。

## Key Changes

### 配置变量

在 bluespot-backend/internal/config/app.yml 的 attach 下增加：

attach:
  view_attach_base_url: "http://localhost:6306/public/uploads"
  upload_dir_path: "./public/uploads"
  view_large_file_base_url: "http://localhost:6306/public/larges"
  upload_large_file_path: "./public/larges"
  chunk_dir_path: "./public/chunks"
  chunk_dir_salt: "bluespot_chunk_upload_salt"

对应新增配置字段和环境变量绑定：

- attach.view_large_file_base_url -> BLUESPOT_ATTACH_VIEW_LARGE_FILE_BASE_PATH
- attach.upload_large_file_path -> BLUESPOT_ATTACH_UPLOAD_LARGE_FILE_PATH
- attach.chunk_dir_path -> BLUESPOT_ATTACH_CHUNK_DIR_PATH
- attach.chunk_dir_salt -> BLUESPOT_ATTACH_CHUNK_DIR_SALT

### 接口行为

POST /api/upload/verify

请求体复用前端：

{
  "fileMd5": "xxx",
  "fileName": "demo.mp4",
  "fileSize": 123456
}

返回：

{
  "isExists": false,
  "url": "",
  "uploadId": "md5(fileMd5+userId+salt)",
  "uploadedChunks": [0, 1, 2]
}

规则：

- 校验 fileMd5/fileName/fileSize 必填且合法。
- 如果最终文件已存在，返回 isExists=true 和访问 URL。
- 否则扫描 chunk_dir_path/{uploadId} 下已有数字分片文件，返回 uploadedChunks。

POST /api/upload/chunk

表单字段复用前端：

- chunk: file
- uploadId: string
- fileMd5: string
- chunkIndex: number
- fileName: string

规则：

- 重新用当前登录用户计算期望 uploadId，必须和请求传入一致。
- 分片保存到 chunk_dir_path/{uploadId}/{chunkIndex}。
- chunkIndex 只允许非负整数，文件名不信任前端路径。
- 成功返回 null data，匹配前端 post<void>。

POST /api/upload/merge

请求体复用前端：

{
  "uploadId": "xxx",
  "fileMd5": "xxx",
  "chunkLength": 10
}

返回：

{
  "url": "http://localhost:6306/public/uploads/xxx.mp4",
  "msg": "合并成功"
}

规则：

- 重新校验 uploadId。
- 检查 0..chunkLength-1 所有分片都存在。
- 按索引顺序合并到 attach.upload_large_file_path/{fileMd5}{原扩展名}。
- 合并完成后删除 chunk_dir_path/{uploadId}。
- 若最终文件已存在，直接返回 前缀 view_large_file_base_url 的 URL，保证重复 merge 幂等。

### 后端结构

按现有分层新增：

- internal/model：新增 chunk upload request/response model。
- internal/controller：新增 UploadController，只负责参数解析、用户 ID 提取、响应封装。
- internal/service：新增 UploadService，负责 uploadId、校验、秒传、分片状态、合并逻辑。
- internal/repository 和 persistence：新增本地文件分片存储实现。
- internal/api/router.go：注册 /api/upload 路由组，并使用 JWTAuthMiddleware()。

- 单元测试 UploadService：
    - 相同 fileMd5 + userId + salt 生成稳定 uploadId。
    - 不同用户生成不同 uploadId。
    - verify 能识别最终文件已存在。
    - verify 能扫描并排序已有分片。
    - chunk 保存拒绝不匹配的 uploadId。
    - merge 缺少任一分片时失败。
    - merge 成功后生成最终文件并清理 chunks 目录。

- 接口测试或手工 HTTP 验证：
    - 登录后调用 verify -> chunk -> merge 全流程。
    - 重复 verify 返回已上传分片。
    - 重复 merge 返回已有最终文件 URL。

- 前端联调：
    - 临时移除 chunk-upload.ts 中 if (flag) return 的调试中断，否则不会真正上传。

## Assumptions

- rootMd5 就是当前前端字段 fileMd5。
- 最终文件名使用 {fileMd5}{原文件扩展名}，便于秒传检测和避免原始文件名冲突。

- 本方案不新增数据库表，上传状态全部来自文件系统。