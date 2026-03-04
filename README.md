# 站点监控服务

一个基于 Gin 和 SQLite 的轻量级监控服务，用于接收站点（JS/内容）变化数据、存储文件快照，并提供动态 Hook 脚本配置。

## ✨ 功能特性

- **数据上报**：提供 RESTful API，接收监控工具上传的 JS 或内容变化（差异 + 原始内容）。
- **内容快照**：存储每次变化的原始文件内容，并计算 MD5 哈希，便于调试和回溯。
- **动态 Hook**：为每个站点配置自定义 Hook 脚本（覆盖原生方法、内容脚本、代码片段），客户端可实时拉取执行。
- **轻量存储**：使用 SQLite 嵌入式数据库，无需额外部署，支持并发读写。
- **API 文档**：集成 Swaggo，自动生成 OpenAPI 规范文档，并提供 Swagger UI 交互式调试界面。

## 🛠 技术栈

- **Go 1.21+**
- **Gin**：Web 框架
- **SQLite**：嵌入式数据库（[mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)）
- **Swaggo**：API 文档生成
- **UUID**：唯一标识

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- Git
- C 编译器（SQLite 驱动依赖 CGO）

### 安装步骤

1. 安装依赖
   ```bash
   go mod download
   ```

2. （可选）安装 swag 命令行工具，用于生成 API 文档
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### 运行服务

```bash
go run main.go
```

服务默认启动在 `http://localhost:8080`，数据文件将保存为项目根目录下的 `monitor.db`。

## 📚 API 文档（Swagger）

### 生成文档

如果你修改了 API 注释，需要重新生成文档：

```bash
swag init
```

该命令会扫描代码中的注释，在 `docs` 目录下生成 `docs.go`、`swagger.json` 和 `swagger.yaml`。

### 访问 Swagger UI

服务启动后，访问：

```
http://localhost:8080/swagger/index.html
```

在 UI 中你可以查看所有 API 端点、请求/响应格式，并直接进行在线调试。

## 🗄 数据库说明

SQLite 数据库文件为 `monitor.db`，启动时自动创建以下表：

- `monitored_sites`：被监控的站点信息
- `change_records`：变化记录（差异）
- `file_contents`：文件内容快照（含 MD5）
- `hook_scripts`：Hook 脚本库
- `site_hooks`：站点与脚本的关联配置（支持 JSON 配置）

## 📁 项目结构

```
.
├── docs/                # Swagger 生成的文档
├── db/                  # 数据库初始化
│   └── sqlite.go
├── handlers/            # HTTP 请求处理
│   ├── upload.go        # 上传变化
│   ├── query.go         # 查询站点/变化
│   ├── content.go       # 查询内容快照
│   └── hooks.go         # 获取站点 Hook 脚本
├── models/              # 数据模型
│   └── models.go
├── static/              # 静态文件（前端界面）
├── templates/           # HTML 模板
├── main.go              # 程序入口
├── go.mod
└── go.sum
```

## 🧪 API 调用示例（curl）

### 1. 上传变化记录

```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: application/json" \
  -d '{
    "site_url": "https://example.com",
    "site_name": "Example Site",
    "change_type": "js",
    "file_path": "/js/app.js",
    "change_diff": "@@ -1,3 +1,4 @@\n old line\n+new line",
    "snapshot_hash": "d41d8cd98f00b204e9800998ecf8427e",
    "content": "console.log(\"hello\");\nconsole.log(\"world\");"
  }'
```

**成功响应**：
```json
{
  "message": "change recorded successfully",
  "site_id": "550e8400-e29b-41d4-a716-446655440000",
  "record_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### 2. 获取站点 Hook 脚本

```bash
curl "http://localhost:8080/api/site-hooks?url=https://example.com"
```

**响应示例**：
```json
[
  {
    "type": "override",
    "content": "window.fetch = function() { console.log('fetch intercepted'); }",
    "config": {
      "logLevel": "debug"
    }
  },
  {
    "type": "content_script",
    "content": "document.body.style.backgroundColor = 'red';"
  }
]
```

### 3. 查询变化记录

```bash
curl "http://localhost:8080/api/changes?site_id=550e8400-e29b-41d4-a716-446655440000&type=js"
```

### 4. 获取所有站点

```bash
curl "http://localhost:8080/api/sites"
```

### 5. 获取文件内容快照

```bash
curl "http://localhost:8080/api/content?site_id=550e8400-e29b-41d4-a716-446655440000&file_path=/js/app.js"
```
