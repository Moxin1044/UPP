<div align="center">

# UPP - Uptime Probe Platform

**现代化网站监控平台 · 深色科技风 · 多项目多用户**

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-6.0-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-7-47A248?style=flat&logo=mongodb)](https://www.mongodb.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

</div>

---

UPP (Uptime Probe Platform) 是一个类似 [Uptime Kuma](https://github.com/louislam/uptime-kuma) 的自托管网站监控平台，采用深蓝色科技感暗色主题，支持多种监控类型、多项目管理、多用户隔离、实时状态推送和通知告警。

![](attachments/image-20260526012537912.png)

![](attachments/image-20260526010335375.png)

![](attachments/image-20260526011903490.png)

![](attachments/image-20260526011921411.png)

![](attachments/image-20260526012321792.png)

![](attachments/image-20260526012359437.png)

![](attachments/image-20260526012445196.png)

![](attachments/image-20260526012458154.png)

![](attachments/image-20260526012521479.png)

---

## ✨ 功能特性

### 监控能力

| 功能 | 说明 |
|------|------|
| **HTTP/HTTPS 监控** | 支持自定义请求方法、请求头、请求体，可校验状态码与响应体关键字 |
| **TCP 端口监控** | 检测目标主机端口是否可达，支持自定义超时时间 |
| **Ping 监控** | 基于 ICMP 协议的 Ping 检测，实时获取延迟与丢包率 |
| **灵活调度** | 每个监控独立配置检测间隔（秒级），Goroutine + Worker Pool 并发调度 |
| **启停控制** | 单个监控可独立启用/暂停，不影响其他监控运行 |

### 项目管理

| 功能 | 说明 |
|------|------|
| **多项目分组** | 按业务线、团队、环境等维度创建项目，将监控归类管理 |
| **项目仪表盘** | 每个项目独立仪表盘，展示监控数量、可用率、平均延迟等统计 |
| **数据隔离** | 项目间监控数据隔离，用户只能查看自己参与的项目的监控 |

### 用户与权限

| 功能 | 说明 |
|------|------|
| **多用户注册** | 支持用户自主注册，管理员也可在后台创建用户 |
| **角色隔离** | 管理员（admin）可管理所有用户和系统设置，普通用户仅管理自己的项目与监控 |
| **JWT 认证** | 基于 JWT 的无状态认证，Token 有效期可配置 |

### 实时与通知

| 功能 | 说明 |
|------|------|
| **WebSocket 推送** | 监控状态变更实时推送到前端，无需刷新页面 |
| **飞书通知** | 通过飞书 Webhook 机器人推送告警消息 |
| **钉钉通知** | 通过钉钉 Webhook 机器人推送告警消息 |
| **WebHook 通知** | 自定义 WebHook 地址，状态变更时 POST 推送 JSON 数据 |
| **邮箱通知** | SMTP 邮件通知，支持自定义发件人与邮件模板 |
| **通知测试** | 创建通知后可一键发送测试消息，验证配置是否正确 |

### 界面与体验

| 功能 | 说明 |
|------|------|
| **深蓝科技风** | 全局深蓝色暗色主题（`#0a1628` / `#0e1729`），毛玻璃卡片、渐变装饰、粒子动画 |
| **仪表盘** | 全局概览与项目级仪表盘，ECharts 可视化延迟趋势图 |
| **中英文切换** | 基于 vue-i18n 的完整国际化，登录页即可切换语言 |
| **响应式布局** | 适配桌面端与移动端，登录/注册页自适应屏幕尺寸 |

---

## 🏗️ 技术架构

### 系统架构图

```
┌─────────────┐     ┌──────────────────────────────────────────┐
│   Browser   │────▶│              Nginx (前端)                 │
│             │◀────│  静态文件 + API 反向代理 + WebSocket 代理  │
└─────────────┘     └──────────┬───────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    Go Backend       │
                    │  (Gin + GORM)       │
                    │                     │
                    │  ┌───────────────┐  │     ┌─────────────┐
                    │  │  HTTP API     │  │     │  WebSocket   │
                    │  │  (Gin Router) │  │     │  Hub         │
                    │  └───────────────┘  │     └─────────────┘
                    │  ┌───────────────┐  │
                    │  │  Scheduler    │  │     ┌─────────────┐
                    │  │  (Worker Pool)│──┼────▶│  Probe       │
                    │  └───────────────┘  │     │  (HTTP/TCP/  │
                    │  ┌───────────────┐  │     │   Ping)      │
                    │  │  Notification │  │     └─────────────┘
                    │  │  Dispatcher   │  │
                    │  └───────────────┘  │
                    └──┬─────────────┬───┘
                       │             │
              ┌────────▼───┐  ┌──────▼──────┐
              │ PostgreSQL │  │   MongoDB   │
              │  (主数据)   │  │  (日志数据)  │
              └────────────┘  └─────────────┘
```

### 后端技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.25 | 高性能编译型语言 |
| Web 框架 | Gin | 轻量级 HTTP 框架，中间件生态丰富 |
| ORM | GORM | Go 最流行的 ORM，支持 PostgreSQL |
| 主数据库 | PostgreSQL 16 | 存储用户、项目、监控、通知等核心数据 |
| 日志数据库 | MongoDB 7 | 存储监控检测结果等高频写入的时序数据 |
| 认证 | JWT (golang-jwt/v5) | 无状态 Token 认证 |
| 实时通信 | Gorilla WebSocket | 双向实时推送监控状态变更 |
| 调度 | Goroutine + Worker Pool | 并发调度监控任务，可控制并发数 |
| Ping 探针 | pro-bing | 基于 ICMP 的 Ping 检测库 |
| 密码加密 | bcrypt (golang.org/x/crypto) | 安全的密码哈希存储 |

### 前端技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Vue 3.5 + TypeScript 6.0 | 组合式 API + 类型安全 |
| 构建 | Vite 8 | 极速开发与构建 |
| 状态管理 | Pinia 3 | 轻量级响应式状态管理 |
| 路由 | Vue Router 4 | SPA 路由，支持路由守卫 |
| UI 组件库 | TDesign Vue Next 1.20 | 腾讯出品，暗色主题支持 |
| HTTP 客户端 | Axios | 请求拦截、Token 注入、错误统一处理 |
| 图表 | ECharts 6 | 延迟趋势图、统计图表，暗色主题适配 |
| 国际化 | vue-i18n 9 | 完整中英文切换 |

### 部署架构

| 组件 | 技术 | 说明 |
|------|------|------|
| 容器化 | Docker + Docker Compose | 一键部署，服务编排 |
| 前端服务 | Nginx | 静态文件托管 + API 反向代理 + WebSocket 代理 |
| 后端服务 | Go Binary | 多阶段构建，Alpine 精简镜像 |

---

## 📁 项目结构

```
UPP/
├── backend/                          # 后端服务 (Go)
│   ├── cmd/
│   │   └── main.go                   # 程序入口，初始化与启动
│   ├── internal/
│   │   ├── api/                      # HTTP Handler 层
│   │   │   ├── auth.go               # 认证接口 (登录/注册/用户信息)
│   │   │   ├── dashboard.go          # 仪表盘接口
│   │   │   ├── monitor.go            # 监控 CRUD + 启停 + 结果查询
│   │   │   ├── notification.go       # 通知 CRUD + 测试
│   │   │   ├── project.go            # 项目 CRUD
│   │   │   ├── setting.go            # 系统设置
│   │   │   └── user.go               # 用户管理 (管理员)
│   │   ├── config/
│   │   │   └── config.go             # 配置加载 (YAML + 环境变量)
│   │   ├── middleware/
│   │   │   ├── auth.go               # JWT 认证中间件
│   │   │   └── cors.go               # 跨域中间件
│   │   ├── model/
│   │   │   ├── database.go           # 数据库初始化 (PostgreSQL + MongoDB)
│   │   │   ├── monitor.go            # 监控模型
│   │   │   ├── monitor_result.go     # 监控结果模型
│   │   │   ├── notification.go       # 通知模型
│   │   │   ├── project.go            # 项目模型
│   │   │   ├── setting.go            # 设置模型
│   │   │   └── user.go               # 用户模型
│   │   ├── notification/
│   │   │   ├── dispatcher.go         # 通知调度器
│   │   │   ├── email.go              # 邮箱通知
│   │   │   ├── feishu.go             # 飞书通知
│   │   │   ├── dingtalk.go           # 钉钉通知
│   │   │   └── webhook.go            # WebHook 通知
│   │   ├── probe/
│   │   │   ├── http.go               # HTTP/HTTPS 探针
│   │   │   ├── tcp.go                # TCP 端口探针
│   │   │   └── ping.go               # Ping 探针
│   │   ├── repository/
│   │   │   ├── monitor.go            # 监控数据访问
│   │   │   ├── monitor_result.go     # 监控结果数据访问
│   │   │   ├── notification.go       # 通知数据访问
│   │   │   ├── project.go            # 项目数据访问
│   │   │   ├── setting.go            # 设置数据访问 (UPSERT)
│   │   │   └── user.go               # 用户数据访问
│   │   ├── scheduler/
│   │   │   └── scheduler.go          # 调度引擎 (Worker Pool)
│   │   ├── service/
│   │   │   ├── auth.go               # 认证业务逻辑
│   │   │   ├── dashboard.go          # 仪表盘业务逻辑
│   │   │   ├── monitor.go            # 监控业务逻辑
│   │   │   ├── notification.go       # 通知业务逻辑
│   │   │   ├── project.go            # 项目业务逻辑
│   │   │   ├── setting.go            # 设置业务逻辑
│   │   │   └── user.go               # 用户业务逻辑
│   │   └── websocket/
│   │       └── hub.go                # WebSocket Hub (广播推送)
│   ├── config.yaml                   # 配置文件
│   ├── Dockerfile                    # 多阶段构建
│   └── go.mod
│
├── frontend/                         # 前端服务 (Vue 3 + TypeScript)
│   ├── src/
│   │   ├── api/                      # API 请求封装
│   │   │   ├── auth.ts               # 认证 API
│   │   │   ├── dashboard.ts          # 仪表盘 API
│   │   │   ├── monitor.ts            # 监控 API
│   │   │   ├── notification.ts       # 通知 API
│   │   │   ├── project.ts            # 项目 API
│   │   │   └── user.ts               # 用户 API
│   │   ├── layouts/
│   │   │   └── MainLayout.vue        # 主布局 (侧边栏 + 顶栏)
│   │   ├── locales/                  # 国际化
│   │   │   ├── zh/index.ts           # 中文翻译
│   │   │   └── en/index.ts           # 英文翻译
│   │   ├── pages/                    # 页面组件
│   │   │   ├── Login.vue             # 登录 (粒子动画 + 毛玻璃)
│   │   │   ├── Register.vue          # 注册 (粒子动画 + 毛玻璃)
│   │   │   ├── Dashboard.vue         # 仪表盘 (统计卡片 + ECharts)
│   │   │   ├── Projects.vue          # 项目列表
│   │   │   ├── ProjectDetail.vue     # 项目详情
│   │   │   ├── Monitors.vue          # 监控列表
│   │   │   ├── MonitorDetail.vue     # 监控详情 (延迟趋势图)
│   │   │   ├── Notifications.vue     # 通知管理
│   │   │   ├── Users.vue             # 用户管理 (管理员)
│   │   │   └── Settings.vue          # 系统设置 (管理员)
│   │   ├── router/
│   │   │   └── index.ts              # 路由配置 + 导航守卫
│   │   ├── stores/
│   │   │   └── user.ts               # 用户状态 (Pinia)
│   │   ├── styles/
│   │   │   └── dark-theme.css        # 全局暗色主题 + TDesign 变量覆盖
│   │   ├── utils/
│   │   │   └── request.ts            # Axios 封装 (拦截器)
│   │   ├── App.vue
│   │   └── main.ts                   # 入口
│   ├── index.html                    # HTML 模板 (theme-mode="dark")
│   ├── nginx.conf                    # Nginx 配置 (API 代理 + WS 代理)
│   ├── Dockerfile
│   └── package.json
│
└── docker-compose.yml                # 服务编排
```

---

## 🚀 快速开始

### 前置要求

- [Docker](https://docs.docker.com/get-docker/) >= 20.10
- [Docker Compose](https://docs.docker.com/compose/install/) >= 2.0

### 一键启动

```bash
# 克隆项目
git clone https://github.com/Moxin1044/UPP.git
cd UPP

# 构建并启动所有服务
docker-compose up -d --build
```

启动完成后访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端界面 | http://localhost | 深蓝科技风监控面板 |
| 后端 API | http://localhost/api/v1 | 通过 Nginx 反向代理 |
| PostgreSQL | localhost:5432 | 可通过客户端直连 |
| MongoDB | localhost:27017 | 可通过客户端直连 |

### 默认账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | admin123 | 管理员 | 可管理用户和系统设置 |

> ⚠️ 生产环境请务必修改默认密码和 JWT 密钥。

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 仅查看后端日志
docker-compose logs -f backend

# 仅查看前端日志
docker-compose logs -f frontend
```

### 停止服务

```bash
docker-compose down

# 停止并删除数据卷（清除所有数据）
docker-compose down -v
```

---

## 📡 API 文档

### 公开接口

| 方法 | 路径 | 说明 | 请求体 |
|------|------|------|--------|
| POST | /api/v1/auth/login | 用户登录 | `{ "username": "", "password": "" }` |
| POST | /api/v1/auth/register | 用户注册 | `{ "username": "", "password": "" }` |

### 认证接口

> 所有认证接口需在请求头携带 `Authorization: Bearer <token>`

#### 仪表盘

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/dashboard/stats | 全局仪表盘统计（总监控数、可用率、平均延迟等） |
| GET | /api/v1/dashboard/project/:id | 项目级仪表盘统计 |

#### 项目管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/projects | 获取当前用户的项目列表 |
| POST | /api/v1/projects | 创建项目 `{ "name": "", "description": "" }` |
| GET | /api/v1/projects/:id | 获取项目详情 |
| PUT | /api/v1/projects/:id | 更新项目 |
| DELETE | /api/v1/projects/:id | 删除项目 |

#### 监控管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/monitors | 获取当前用户的监控列表 |
| POST | /api/v1/monitors | 创建监控 |
| GET | /api/v1/monitors/:id | 获取监控详情 |
| PUT | /api/v1/monitors/:id | 更新监控（类型、目标、间隔等） |
| DELETE | /api/v1/monitors/:id | 删除监控 |
| PATCH | /api/v1/monitors/:id/toggle | 启停监控 |
| GET | /api/v1/monitors/:id/results | 获取监控检测结果（支持分页） |
| GET | /api/v1/monitors/:id/stats | 获取监控统计（可用率、平均延迟） |
| GET | /api/v1/projects/:id/monitors | 获取项目下的监控列表 |

#### 通知管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/notifications | 获取通知列表 |
| POST | /api/v1/notifications | 创建通知 |
| GET | /api/v1/notifications/:id | 获取通知详情 |
| PUT | /api/v1/notifications/:id | 更新通知 |
| DELETE | /api/v1/notifications/:id | 删除通知 |
| POST | /api/v1/notifications/:id/test | 发送测试通知 |

#### 实时通信

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ws | WebSocket 连接（监控状态变更实时推送） |

### 管理员接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users | 获取用户列表 |
| POST | /api/v1/users | 创建用户 |
| PUT | /api/v1/users/:id/password | 修改用户密码 |
| DELETE | /api/v1/users/:id | 删除用户 |
| GET | /api/v1/settings | 获取系统设置 |
| POST | /api/v1/settings | 更新系统设置（UPSERT） |

---

## 🎨 主题定制

UPP 采用深蓝色科技感暗色主题，所有颜色通过 CSS 变量控制，可轻松定制：

### 核心色彩变量

```css
:root {
  /* 背景色 */
  --upp-bg-primary: #0a1628;        /* 页面主背景 */
  --upp-bg-secondary: #0e1729;      /* 容器/卡片背景 */
  --upp-bg-card: rgba(14, 23, 41, 0.92);  /* 毛玻璃卡片 */

  /* 边框 */
  --upp-border: rgba(255, 255, 255, 0.06);  /* 默认边框 */
  --upp-border-hover: rgba(0, 212, 255, 0.35);  /* 悬停边框 */

  /* 文字 */
  --upp-text-primary: #fff;                        /* 主文字 */
  --upp-text-secondary: rgba(255, 255, 255, 0.65); /* 次要文字 */
  --upp-text-tertiary: rgba(255, 255, 255, 0.35);  /* 辅助文字 */

  /* 强调色 */
  --upp-accent: #00d4ff;                              /* 科技蓝 */
  --upp-accent-gradient: linear-gradient(135deg, #00a8ff, #0066ff);  /* 渐变 */

  /* 状态色 */
  --upp-success: #52c41a;
  --upp-danger: #ff4d4f;
  --upp-warning: #faad14;
}
```

### TDesign 暗色变量覆盖

通过覆盖 `--td-*` 变量，将 TDesign 默认的灰色暗色主题替换为深蓝色系：

```css
:root[theme-mode="dark"] {
  --td-bg-color-page: #0a1628;
  --td-bg-color-container: #0e1729;
  --td-bg-color-container-hover: #111d33;
  --td-brand-color: #00d4ff;
  /* ... 更多变量见 dark-theme.css */
}
```

---

## 💻 本地开发

### 环境要求

- Go >= 1.25
- Node.js >= 20
- PostgreSQL >= 16
- MongoDB >= 7

### 后端开发

```bash
cd backend

# 安装依赖
go mod download

# 启动 PostgreSQL 和 MongoDB（可用 Docker）
docker run -d --name upp-postgres -e POSTGRES_USER=upp -e POSTGRES_PASSWORD=upp123 -e POSTGRES_DB=upp -p 5432:5432 postgres:16-alpine
docker run -d --name upp-mongo -p 27017:27017 mongo:7

# 修改 config.yaml 中的数据库地址为 localhost

# 启动后端
go run ./cmd/main.go
```

后端默认运行在 `http://localhost:8080`。

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端开发服务器默认运行在 `http://localhost:5173`，API 请求会代理到后端。

> 开发时需修改 `vite.config.ts` 中的代理地址指向本地后端。

---

## ⚙️ 配置说明

### 后端配置 (`backend/config.yaml`)

```yaml
server:
  port: "8080"           # 服务端口

database:
  host: postgres          # PostgreSQL 主机
  port: "5432"            # PostgreSQL 端口
  user: upp               # 数据库用户名
  password: upp123        # 数据库密码
  dbname: upp             # 数据库名
  sslmode: disable        # SSL 模式

mongodb:
  uri: mongodb://mongo:27017  # MongoDB 连接地址
  database: upp               # MongoDB 数据库名

jwt:
  secret: upp-jwt-secret-key-2024  # JWT 签名密钥（生产环境务必修改）
  expire: 72                        # Token 有效期（小时）
```

### 环境变量覆盖

Docker 部署时可通过环境变量覆盖配置：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `DB_HOST` | PostgreSQL 主机 | postgres |
| `MONGO_URI` | MongoDB 连接地址 | mongodb://mongo:27017 |

### Nginx 配置 (`frontend/nginx.conf`)

- 前端静态文件托管
- `/api/` 路径反向代理到后端
- `/api/v1/ws` WebSocket 代理

---

## 📊 数据模型

### 核心实体关系

```
User ──1:N──▶ Project ──1:N──▶ Monitor ──1:N──▶ MonitorResult
                                    │
                                    └──1:N──▶ Notification
```

| 模型 | 存储位置 | 说明 |
|------|----------|------|
| User | PostgreSQL | 用户信息（用户名、密码哈希、角色） |
| Project | PostgreSQL | 项目（名称、描述、所属用户） |
| Monitor | PostgreSQL | 监控配置（类型、目标、间隔、启用状态） |
| MonitorResult | MongoDB | 监控结果（状态、延迟、消息、时间戳） |
| Notification | PostgreSQL | 通知配置（类型、Webhook 地址、启用状态） |
| Setting | PostgreSQL | 系统设置（键值对，UPSERT 更新） |

---

## 🔔 通知类型配置

### 飞书 (Feishu)

```json
{
  "type": "feishu",
  "name": "飞书告警",
  "config": {
    "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
  }
}
```

### 钉钉 (DingTalk)

```json
{
  "type": "dingtalk",
  "name": "钉钉告警",
  "config": {
    "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=xxx"
  }
}
```

### WebHook

```json
{
  "type": "webhook",
  "name": "自定义 WebHook",
  "config": {
    "url": "https://your-server.com/webhook",
    "method": "POST"
  }
}
```

### 邮箱 (Email)

```json
{
  "type": "email",
  "name": "邮箱告警",
  "config": {
    "smtp_host": "smtp.example.com",
    "smtp_port": 465,
    "username": "alert@example.com",
    "password": "smtp-password",
    "to": "admin@example.com"
  }
}
```

---

## 🛡️ 安全说明

- 密码使用 bcrypt 哈希存储，不可逆
- JWT Token 有效期可配置，过期后需重新登录
- API 接口通过中间件校验 Token 和用户角色
- 管理员接口额外校验用户角色
- CORS 中间件限制跨域访问
- 生产环境建议：
  - 修改 `jwt.secret` 为高强度随机字符串
  - 修改数据库默认密码
  - 配置 HTTPS（Nginx 层添加 SSL 证书）
  - 限制 PostgreSQL 和 MongoDB 端口不对外暴露

---

## 📄 License

[MIT License](LICENSE)
