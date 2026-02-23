# CoreDNS Admin

CoreDNS 的 Web 管理界面，用于通过 etcd 管理 CoreDNS 的 DNS 记录。后端使用 Go + Gin，前端使用 Vue 2 + Element UI，数据与用户信息均存储在 etcd 中。

## 项目结构

```
coredns-admin-0.4/
├── main.go              # 程序入口，支持 -adduser 添加用户、-C 指定配置文件
├── config.yaml          # 配置文件（etcd 地址、端口、路径前缀等）
├── go.mod / go.sum      # Go 依赖
├── config/              # 配置加载与默认值
├── controller/          # HTTP 接口（登录、记录 CRUD、域名列表）
├── middleware/          # JWT 鉴权、CORS
├── model/               # 数据模型（Record、User、Etcd 结构）
├── router/              # 路由与静态资源（生产时提供 dist）
├── service/             # 业务逻辑（etcd 读写、密码哈希校验）
└── views/               # 前端 Vue 项目
    ├── package.json
    ├── src/
    │   ├── views/       # 登录页、首页、记录管理页
    │   ├── router/      # 前端路由（/login, /records）
    │   └── main.js
    └── vue.config.js
```

## 环境要求

- **Go** 1.22+（依赖见 `go.mod`，含 etcd client v3、gin、jwt v5 等）
- **Node.js**（用于构建前端，建议 14+）
- **Yarn** 或 npm（前端依赖）
- **etcd**（CoreDNS 使用的 etcd 集群需可访问）

## 配置说明

编辑项目根目录下的 `config.yaml`：

```yaml
# 服务监听地址，留空表示所有网卡
# host: ""

# 用户信息在 etcd 中的路径前缀
# user_etcd_path: /user/coredns

# 服务端口
port: 8088

etcd:
  # etcd 集群地址（与 CoreDNS 使用同一 etcd）
  endpoint: [http://10.1.1.224:2379]
  # DNS 记录在 etcd 中的路径前缀（需与 CoreDNS 配置一致）
  path_prefix: /coredns
  timeout: 5
  # 若 etcd 开启认证可取消注释：
  # username: ""
  # password: ""
  # TLS 证书（cert, key, ca 三个文件路径）
  # tls: [cert, key, ca]
```

请根据实际环境修改 `etcd.endpoint` 和 `etcd.path_prefix`，确保与 CoreDNS 的 etcd 配置一致。

**本地无 etcd 时**：可用下面的 **Docker Compose 一键部署**（同时包含 etcd + coredns-admin + 前端 views）。

## Docker 一键部署（etcd + coredns-admin + CoreDNS + views）

项目根目录提供 `Dockerfile` 与 `docker-compose.yml`，可同时启动 etcd、CoreDNS、后端和已打包的前端。

```bash
# 构建并启动（默认管理员 admin / admin123456，可通过 .env 或环境变量覆盖）
docker compose up -d --build
```

- 访问：**http://localhost:8088**，使用 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 登录（默认 `admin` / `admin123456`）。
- 首次启动时自动将配置的管理员写入 etcd，无需手动 `-adduser`。
- 自定义账号：复制 `.env.example` 为 `.env`，修改 `ADMIN_USERNAME`、`ADMIN_PASSWORD` 后重启。
- 镜像内已使用 `docker/config.yaml`，etcd 地址为 compose 服务名 `http://etcd:2379`。
- 仅重建后端镜像：`docker compose build coredns-admin && docker compose up -d coredns-admin`。

## 使用步骤

### 1. 添加管理员用户（首次使用）

在项目根目录执行：

```bash
go run . -adduser
```

按提示输入用户名和密码（密码至少 6 位）。用户信息会写入 etcd（路径由 `user_etcd_path` 决定）。

也可先编译再执行：

```bash
go build -o coredns-admin
./coredns-admin -adduser
```

### 2. 构建前端

进入前端目录并安装依赖、构建：

```bash
cd views
yarn install
yarn build
```

构建产物在 `views/dist/`，后端会从 `./dist/` 提供静态文件和前端页面。

### 3. 启动后端服务

在项目根目录（即包含 `config.yaml` 和 `views/dist` 的目录）执行：

```bash
go run .
```

或使用自定义配置文件：

```bash
go run . -C /path/to/your/config.yaml
```

生产环境建议先编译再运行：

```bash
go build -o coredns-admin
./coredns-admin -C config.yaml
```

默认会在 `0.0.0.0:8088` 启动（端口以 `config.yaml` 为准）。

### 4. 访问管理界面

浏览器打开：

- 地址：`http://localhost:8088`（或你配置的 host:port）
- 未登录会跳转到登录页，使用步骤 1 中创建的用户名和密码登录
- 登录后可管理 DNS 记录：查看、添加、编辑、删除

## 开发说明

- **仅改前端**：在 `views` 下执行 `yarn serve`，会启动前端开发服务器（如 8080）。此时需后端单独运行，并配置前端请求的 API 地址（如代理到 `http://localhost:8088`），否则会有跨域或接口不可用问题。
- **后端单独运行**：在根目录 `go run .` 会读取 `config.yaml` 并连接 etcd；若前端未执行 `yarn build`，则 `./dist` 不存在，访问首页可能 404，需先完成上述“构建前端”步骤。
- **API 概览**：
  - `POST /login`：登录，获取 JWT
  - `GET /api/v1/records`：获取记录列表（需 JWT）
  - `POST /api/v1/record`：新增记录（需 JWT）
  - `PUT /api/v1/record/:key`：更新记录（需 JWT）
  - `DELETE /api/v1/record/:key`：删除记录（需 JWT）
  - `GET /api/v1/domains`：获取域名列表（需 JWT）

## 支持的 DNS 记录类型

- A、AAAA、CNAME、PTR、MX、TXT、SRV（NS 在模型中存在但前端为禁用状态）

SRV 的 Content 格式为：`权重 端口 目标`，例如：`10 80 target.example.com`。

## 注意事项

1. 确保本程序连接的 **etcd 与 CoreDNS 使用的 etcd 一致**，且 `path_prefix` 与 CoreDNS 中 etcd 插件配置一致，否则修改不会生效到 CoreDNS。
2. 首次使用必须先执行 `-adduser` 至少一次，否则无法登录。
3. 生产部署时请使用 HTTPS 并妥善保管 `config.yaml`（若配置了 etcd 密码或 TLS 路径）。

## 许可证

请以项目原有许可证为准。
