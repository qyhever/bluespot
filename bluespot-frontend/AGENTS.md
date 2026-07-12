# bluespot-frontend

本文件是 PC web 端协作入口。进入本子项目修改前，先阅读本文件，再按涉及模块继续查看对应源码。

## 技术栈

- Vue 3 + Vite
- TypeScript
- Vue Router
- Pinia
- Tailwind CSS 4（通过 `@tailwindcss/vite` 接入）
- lint/format：ESLint + oxlint + oxfmt

## 目录结构

- `src/main.ts`：应用入口，创建 Vue 应用并挂载 Pinia、Router。
- `src/App.vue`：根组件，目前只渲染 `router-view`。
- `src/router/index.ts`：路由配置，当前包含 `/` 与 `/about`。
- `src/views/`：页面级组件。
- `src/components/`：通用组件目录，当前为空。
- `src/stores/`：Pinia 状态管理，当前只有模板示例 `counter.ts`。
- `src/utils/fetch.ts`：底层请求封装，支持 fetch 与带进度事件的 XHR。
- `src/utils/request.ts`：业务请求封装，包含统一 `baseURL`、token 注入、错误提示、401 刷新与重试逻辑。
- `src/utils/version-checker.ts`：前端版本检测，读取构建产物中的 `meta.json`。
- `src/assets/`：全局样式与静态资源。
- `types/`：全局类型声明，包括接口响应、分页类型、环境变量和构建常量。
- `build/meta.ts`：Vite 构建插件，生成构建 hash 与 `meta.json`。
- `public/`：Vite 公共资源目录。

## 常用命令

本项目通过 `preinstall` 限制使用 pnpm，请优先使用 pnpm：

```sh
pnpm install
pnpm dev
pnpm build
pnpm lint
pnpm format
```

命令说明：

- `pnpm dev`：启动 Vite 开发服务，默认端口 `5175`。
- `pnpm build`：先执行 `vue-tsc --build` 类型检查，再执行 Vite 构建。
- `pnpm lint`：依次运行 oxlint 与 ESLint，并自动修复可修复问题。
- `pnpm format`：使用 oxfmt 格式化 `src/`。
- `pnpm preview`：预览构建产物。

## 本地联调

开发服务配置在 `vite.config.ts`：

- 前端端口：`5175`
- API 代理前缀：`/bluespot/api`
- 代理目标：`http://localhost:6300`
- 代理 rewrite：`/bluespot/api/...` 会转发为后端 `/api/...`

业务请求默认 `baseURL` 在 `src/utils/request.ts` 中配置为：

```ts
const defaultOptions = {
  baseURL: '/bluespot/api',
}
```

新增接口调用优先复用 `src/utils/request.ts` 暴露的 `get`、`post`、`put`、`del`、`patch`。

## 构建与版本信息

`vite.config.ts` 会在非 `dev` 模式下定义：

- `LOCAL_BUILD_HASH`
- `LOCAL_BUILD_TIME`

`build/meta.ts` 会在构建时根据 `src/**/*` 与 `public/**/*` 计算 hash，并写出 `dist/meta.json`。`src/utils/version-checker.ts` 通过拉取 `meta.json` 检测线上版本更新。

如果调整构建产物目录、`base`、静态资源路径或部署方式，需要同步检查 `build/meta.ts` 与 `version-checker.ts`。

## 代码约定

- 使用 `@` 指向 `src`，配置位于 `vite.config.ts`。
- 新页面放在 `src/views/`，并在 `src/router/index.ts` 注册路由。
- 可复用 UI 放在 `src/components/`，页面私有组件可放在对应页面附近，但要保持边界清晰。
- 业务状态放在 `src/stores/`，优先使用 Pinia composition store 写法。
- 接口响应类型优先复用 `types/global.d.ts` 中的 `IApiResponse<T>`、`IPaginationResponse<T>` 等全局类型。
- 请求层不要绕过统一封装，除非是 `meta.json`、静态资源或第三方直连请求。
- 页面组件不要直接写业务请求。优先将接口类型和请求函数抽到请求层：跨页面或模块级接口放在 `src/api/*.ts`，页面私有接口可放在当前页面同级 `service.ts`。请求函数内部复用 `src/utils/request.ts` 暴露的 `get`、`post`、`put`、`del`、`patch`。
- 不要把后端地址硬编码到业务组件中，开发代理由 `vite.config.ts` 统一维护。

## 当前注意事项

当前代码仍有模板和迁移中的痕迹，修改前请先确认依赖与文件是否已补齐：

- `src/main.ts` 引入的是 `./assets/main.css`，但当前目录中只看到 `src/assets/base.css`。
- `src/utils/request.ts` 引用了 `antd`、`@/stores/user`、`./navigate`，当前 `package.json` 与 `src/` 中未看到对应依赖或文件。
- `eslint.config.ts` 读取 `.oxlintrc.json`，当前前端目录未看到该文件。
- `src/stores/counter.ts`、`src/views/HomeView.vue`、`src/views/AboutView.vue` 仍接近 Vue 模板示例。

遇到上述问题时，不要直接删除相关逻辑；先判断是缺文件、待迁移代码，还是计划中的业务模块。
