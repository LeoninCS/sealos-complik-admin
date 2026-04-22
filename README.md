# sealos-complik-admin

[English](#english) | [中文](#zh-cn)

<a id="english"></a>

## English

Administrative backend service for managing project configuration and user compliance records. The service is built with Gin and GORM, stores data in MySQL, and runs automatic schema migration on startup.

### Features

- Health check endpoint for service monitoring
- CRUD APIs for project configs and user commitments
- Record and query APIs for violations, bans, and unban operations
- User status endpoints for ban and violation checks
- YAML-based configuration and file logging
- Docker image for containerized deployment

### Tech Stack

- Go `1.26.1`
- Gin
- GORM
- MySQL
- Docker

### Project Structure

```text
.
|-- cmd/                     # Application entrypoint
|-- configs/                 # YAML configuration files
|-- internal/
|   |-- infra/               # Config, logger, database, migration
|   |-- modules/             # Domain modules
|   |   |-- ban/
|   |   |-- commitment/
|   |   |-- projectconfig/
|   |   |-- unban/
|   |   `-- violation/
|   `-- router/              # HTTP route registration
|-- test/postman.json        # Postman collection
|-- Dockerfile
`-- start.sh                 # Local startup helper
```

### Modules

- `projectconfig`: store and manage project-level configuration records
- `commitment`: manage uploaded commitment file metadata per user
- `violation`: track user violation records
- `ban`: track ban history and active ban status
- `unban`: track unban actions

### Requirements

- Go `1.26.1` or later
- MySQL instance reachable by the application
- An existing database named `sealos-complik-admin` or a custom database configured in the YAML file

Example database creation:

```sql
CREATE DATABASE `sealos-complik-admin`
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;
```

### Configuration

The application reads `configs/config.yaml` by default.

```yaml
port: 8080

database:
  host: localhost
  port: 3306
  username: root
  password: sealos123
  name: sealos-complik-admin

log_dir: logs
```

Notes:

- Tables are auto-migrated at startup.
- The application does not create the MySQL database itself.
- Logs are written to the directory configured by `log_dir`.
- `oss.endpoint` now uses an S3-compatible endpoint (such as MinIO), for example `http://minio.objectstorage-system.svc.cluster.local`.
- `oss.public_base_url` should be a URL reachable by the browser if you want frontend users to open uploaded files directly.

### Run Locally

1. Update `configs/config.yaml` for your environment.
2. Make sure MySQL is running and the target database already exists.
3. Start the service with one of the following commands:

```bash
go run ./cmd
```

```bash
./start.sh
```

`start.sh` checks port `8080`, stops any process already listening on that port, and then starts the app.

### Run with Docker

Build the image:

```bash
docker build -t sealos-complik-admin .
```

Run the container:

```bash
docker run --rm -p 8080:8080 sealos-complik-admin
```

If MySQL is not running inside the same network as the container, update `configs/config.yaml` with a reachable database host before building or running the image.

### API Overview

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/configs` | Create project config |
| `GET` | `/api/configs` | List project configs |
| `GET` | `/api/configs/:config_name` | Get project config by name |
| `PUT` | `/api/configs/:config_name` | Update project config |
| `DELETE` | `/api/configs/:config_name` | Delete project config |
| `POST` | `/api/commitments` | Create commitment |
| `GET` | `/api/commitments` | List commitments |
| `GET` | `/api/commitments/:namespace` | Get commitment by namespace |
| `PUT` | `/api/commitments/:namespace` | Update commitment |
| `DELETE` | `/api/commitments/:namespace` | Delete commitment |
| `POST` | `/api/violations` | Create violation record |
| `GET` | `/api/violations` | List violation records |
| `GET` | `/api/violations/:namespace` | Get violations by namespace |
| `DELETE` | `/api/violations/:namespace` | Delete violations by namespace |
| `GET` | `/api/namespaces/:namespace/violations-status` | Check whether a namespace has violations |
| `POST` | `/api/bans` | Create ban record |
| `GET` | `/api/bans` | List ban records |
| `GET` | `/api/bans/:namespace` | Get bans by namespace |
| `DELETE` | `/api/bans/id/:id` | Delete a ban record by id |
| `GET` | `/api/namespaces/:namespace/ban-status` | Check whether a namespace is banned |
| `POST` | `/api/unbans` | Create unban record |
| `GET` | `/api/unbans` | List unban records |
| `GET` | `/api/unbans/:namespace` | Get unban records by namespace |
| `DELETE` | `/api/unbans/id/:id` | Delete an unban record by id |

### Example Requests

Health check:

```bash
curl http://localhost:8080/health
```

Create a project config:

```bash
curl -X POST http://localhost:8080/api/configs \
  -H "Content-Type: application/json" \
  -d '{
    "config_name": "project-config-demo",
    "config_type": "json",
    "config_value": {"enabled": true, "threshold": 3},
    "description": "Demo config"
  }'
```

Create a commitment:

```bash
curl -X POST http://localhost:8080/api/commitments \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "ns-demo",
    "file_name": "commitment.pdf",
    "file_url": "https://oss.example.com/commitments/commitment.pdf"
  }'
```

### API Collection

Import `test/postman.json` into Postman to quickly try the available APIs.

---

<a id="zh-cn"></a>

## 中文

这是一个用于管理项目配置和用户合规记录的管理后台服务，基于 Gin 和 GORM 构建，数据存储在 MySQL 中，并在启动时自动执行表结构迁移。

### 功能特性

- 提供健康检查接口，便于服务监控
- 提供项目配置和用户承诺书的 CRUD 接口
- 提供违规、封禁、解封相关记录的录入与查询接口
- 提供用户封禁状态与违规状态查询接口
- 使用 YAML 配置文件和本地日志目录
- 提供 Docker 镜像构建能力，便于容器化部署

### 技术栈

- Go `1.26.1`
- Gin
- GORM
- MySQL
- Docker

### 项目结构

```text
.
|-- cmd/                     # 应用入口
|-- configs/                 # YAML 配置文件
|-- internal/
|   |-- infra/               # 配置、日志、数据库、迁移
|   |-- modules/             # 业务模块
|   |   |-- ban/
|   |   |-- commitment/
|   |   |-- projectconfig/
|   |   |-- unban/
|   |   `-- violation/
|   `-- router/              # HTTP 路由注册
|-- test/postman.json        # Postman 集合
|-- Dockerfile
`-- start.sh                 # 本地启动辅助脚本
```

### 模块说明

- `projectconfig`：管理项目级配置项
- `commitment`：管理用户承诺文件元数据
- `violation`：记录用户违规信息
- `ban`：记录封禁历史和当前封禁状态
- `unban`：记录解封操作

### 运行要求

- Go `1.26.1` 或更高版本
- 可被服务访问的 MySQL 实例
- 已存在的 `sealos-complik-admin` 数据库，或在配置文件中指定其他数据库名

示例建库语句：

```sql
CREATE DATABASE `sealos-complik-admin`
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;
```

### 配置说明

应用默认读取 `configs/config.yaml`。

```yaml
port: 8080

database:
  host: localhost
  port: 3306
  username: root
  password: sealos123
  name: sealos-complik-admin

log_dir: logs
```

说明：

- 服务启动时会自动执行数据表迁移。
- 应用不会自动创建 MySQL 数据库本身。
- 日志会写入 `log_dir` 指定的目录。

### 本地运行

1. 根据你的环境修改 `configs/config.yaml`。
2. 确认 MySQL 已启动，且目标数据库已经存在。
3. 使用以下任一命令启动服务：

```bash
go run ./cmd
```

```bash
./start.sh
```

`start.sh` 会先检查 `8080` 端口，占用时会停止对应进程，然后再启动应用。

### Docker 运行

构建镜像：

```bash
docker build -t sealos-complik-admin .
```

运行容器：

```bash
docker run --rm -p 8080:8080 sealos-complik-admin
```

如果 MySQL 不在容器同一网络内，请在构建或运行前先把 `configs/config.yaml` 中的数据库地址改成容器可访问的地址。

### 接口概览

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/health` | 健康检查 |
| `POST` | `/api/configs` | 创建项目配置 |
| `GET` | `/api/configs` | 查询项目配置列表 |
| `GET` | `/api/configs/:config_name` | 按名称查询项目配置 |
| `PUT` | `/api/configs/:config_name` | 更新项目配置 |
| `DELETE` | `/api/configs/:config_name` | 删除项目配置 |
| `POST` | `/api/commitments` | 创建承诺记录 |
| `GET` | `/api/commitments` | 查询承诺记录列表 |
| `GET` | `/api/commitments/:namespace` | 按 namespace 查询承诺记录 |
| `PUT` | `/api/commitments/:namespace` | 更新承诺记录 |
| `DELETE` | `/api/commitments/:namespace` | 删除承诺记录 |
| `POST` | `/api/violations` | 创建违规记录 |
| `GET` | `/api/violations` | 查询违规记录列表 |
| `GET` | `/api/violations/:namespace` | 按 namespace 查询违规记录 |
| `DELETE` | `/api/violations/:namespace` | 删除 namespace 违规记录 |
| `GET` | `/api/namespaces/:namespace/violations-status` | 查询 namespace 是否有违规记录 |
| `POST` | `/api/bans` | 创建封禁记录 |
| `GET` | `/api/bans` | 查询封禁记录列表 |
| `GET` | `/api/bans/:namespace` | 按 namespace 查询封禁记录 |
| `DELETE` | `/api/bans/id/:id` | 按 id 删除封禁记录 |
| `GET` | `/api/namespaces/:namespace/ban-status` | 查询 namespace 是否处于封禁状态 |
| `POST` | `/api/unbans` | 创建解封记录 |
| `GET` | `/api/unbans` | 查询解封记录列表 |
| `GET` | `/api/unbans/:namespace` | 按 namespace 查询解封记录 |
| `DELETE` | `/api/unbans/id/:id` | 按 id 删除解封记录 |

### 请求示例

健康检查：

```bash
curl http://localhost:8080/health
```

创建项目配置：

```bash
curl -X POST http://localhost:8080/api/configs \
  -H "Content-Type: application/json" \
  -d '{
    "config_name": "project-config-demo",
    "config_type": "json",
    "config_value": {"enabled": true, "threshold": 3},
    "description": "Demo config"
  }'
```

创建承诺记录：

```bash
curl -X POST http://localhost:8080/api/commitments \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "ns-demo",
    "file_name": "commitment.pdf",
    "file_url": "https://oss.example.com/commitments/commitment.pdf"
  }'
```

### 接口调试

可以将 `test/postman.json` 导入 Postman，快速体验当前仓库提供的接口。
