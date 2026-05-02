# tech-muyi-shenji

面向**未来开发平台**中「**应用中心**」的后端**雏形**：集中管理平台上各应用的注册信息、导航菜单与页面配置，并为运行态提供稳定的只读发布物。

当前实现基于 [Beego v2](https://github.com/beego/beego) 的 HTTP 服务，用本地 JSON 文件承载「应用—菜单—页面」元数据，并完成开发态、版本快照与发布态的流转；后续可在此基础上接入统一认证、存储与多租户等平台能力。

## 功能概览

- **应用（App）**：创建、更新、删除；发布、上线、回滚；按 `appCode` 或 `appCode.版本` 查询开发/历史版本信息；对外提供已发布应用的 `app.json`。
- **菜单（Menu）**：增删改、排序、从树中移除；按应用查询菜单树；对外提供已发布应用的 `menu.json`。
- **页面（Page）**：按 `appCode` + `menuCode` 保存与读取页面 JSON；对外提供已发布路径下的 `pages/{menuCode}.json`。

数据全部落在仓库下的 `data/` 目录（开发态 `data/dev/`、版本快照 `data/version/`、发布态 `data/app/`、临时 `data/temp/`），无独立数据库。

## 技术栈

- Go（模块见 `go.mod`，当前声明 `go 1.22.12`）
- `github.com/beego/beego/v2` Web 与路由

## 快速开始

```bash
go mod download
go run .
```

默认配置见 `conf/app.conf`：应用名 `tech-muyi-shenji`，HTTP 端口 **10011**，`runmode` 为 `local`；各环境片段通过 `include` 引入（`app-dev.conf`、`app-local.conf` 等）。

## 主要 API 前缀

| 前缀 | 说明 |
|------|------|
| `/api/shenji/app/...` | 应用生命周期与查询 |
| `/api/shenji/menu/...` | 菜单维护与查询 |
| `/api/shenji/pages/...` | 页面读写 |
| `/api/shenji/:appCode/app.json` | 已发布应用信息（只读） |
| `/api/shenji/:appCode/menu.json` | 已发布菜单（只读） |
| `/api/shenji/:appCode/pages/:menuCode.json` | 已发布页面（只读） |

具体方法与 HTTP 动词见 `routers/router.go`。

## 目录结构（简要）

| 路径 | 说明 |
|------|------|
| `main.go` | 入口：控制台日志 + `beego.Run()` |
| `routers/` | 路由注册 |
| `controllers/` | 控制器（应用、菜单、页面及基类） |
| `app/` | 应用/菜单/路径与业务逻辑 |
| `common/` | 通用查询、结果封装等 |
| `utils/` | JSON 文件读写、版本工具等 |
| `conf/` | Beego 配置与各环境片段 |
| `data/` | 运行时数据目录（建议勿提交敏感内容，按 `.gitignore` 处理） |
| `Dockerfile` / `start.sh` | 容器构建与启动脚本（镜像内默认暴露 **10011**） |

## 核心数据模型（摘录）

- **App**：`appId`、`appCode`、`appName`、`version`、`description`（见 `app/appDTO.go`）。
- **Menu**：树形结构，含 `menuCode`、`menuLabel`、`menuPath`、`menuType`（如 group / imperative / lowCode / link）等。

## 容器运行

```bash
docker build -t tech-muyi-shenji .
docker run -p 10011:10011 tech-muyi-shenji
```

说明：`Dockerfile` 基础镜像中的 Go 版本与本地 `go.mod` 可能不一致，若构建失败可统一 Go 版本或调整基础镜像。
