<div align="center">

# [![PongHub](imgs/band.png)](https://health.ch3nyang.top)

🌏 [Live Demo](https://health.ch3nyang.top) | 📖 [English](README.md)

</div>

## 简介

PongHub 是一个开源的服务状态监控网站，旨在帮助用户监控和验证服务的可用性。它支持

- **🕵️ 零侵入监控** - 无需改动代码即可实现全功能监控
- **🚀 一键部署** - 通过 Actions 自动构建，一键部署至 Github Pages
- **🌐 全平台支持** - 兼容 OpenAI 等公共服务及私有化部署
- **🔍 多端口探测** - 单服务支持同时监控多个端口状态
- **🤖 智能响应验证** - 精准匹配状态码及正则表达式校验响应体
- **🛠️ 自定义请求引擎** - 自由配置请求头/体、超时和重试策略
- **🔒 SSL 证书监控** - 自动检测 SSL 证书过期并发送通知
- **📊 实时状态展示** - 直观的服务响应时间、响应状态记录
- **⚠️ 异常告警通知** - 利用 GitHub Actions 实现异常告警通知

![浏览器截图](imgs/browser_CN.png)

## 快速开始

1. Star 并 Fork [PongHub](https://github.com/WCY-dt/ponghub)

2. 修改根目录下的 [`config.yaml`](config.yaml) 文件，配置你的服务检查项

3. 修改根目录下的 [`CNAME`](CNAME) 文件，配置你的自定义域名

   > 如果你不需要自定义域名，请删除 `CNAME` 文件

4. 提交修改并推送到你的仓库，GitHub Actions 将自动更新，无需干预

> [!TIP]
> 默认情况下，GitHub Actions 每 30 分钟运行一次。如果你需要更改运行频率，请修改 [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml) 文件中的 `cron` 表达式。
>
> 请不要将频率设置过高，以免触发 GitHub 的限制。

> [!IMPORTANT]
> 如果 GitHub Actions 未正常自动触发，手动触发一次即可。
>
> 请注意打开 GitHub Pages，并开启 GitHub Actions 的通知权限。

## 配置说明

### 基本配置

配置文件 `config.yaml` 的格式如下：

| 字段                                  | 类型  | 描述                        | 必填 | 备注                             |
|-------------------------------------|-----|---------------------------|----|--------------------------------|
| `display_num`                       | 整数  | 首页显示的服务数量                 | ✖️ | 默认 72 个                        |
| `timeout`                           | 整数  | 每次请求的超时时间，单位为秒            | ✖️ | 单位为秒，默认 5 秒                    |
| `max_retry_times`                   | 整数  | 请求失败时的重试次数                | ✖️ | 默认 2 次                         |
| `max_log_days`                      | 整数  | 日志保留天数，超过此天数的日志将被删除       | ✖️ | 默认 3 天                         |
| `cert_notify_days`                  | 整数  | SSL 证书过期前通知的天数            | ✖️ | 默认 7 天                         |
| `services`                          | 数组  | 服务列表                      | ✔️ |                                |
| `services.name`                     | 字符串 | 服务名称                      | ✔️ |                                |
| `services.endpoints`                | 数组  | 端口列表                      | ✔️ |                                |
| `services.endpoints.url`            | 字符串 | 请求的 URL                   | ✔️ |                                |
| `services.endpoints.method`         | 字符串 | 请求的 HTTP 方法               | ✖️ | 支持 `GET`/`POST`/`PUT`，默认 `GET` |
| `services.endpoints.headers`        | 对象  | 请求头内容                     | ✖️ | 键值对形式，支持自定义请求头                 |
| `services.endpoints.body`           | 字符串 | 请求体内容                     | ✖️ | 仅在 `POST`/`PUT` 请求时使用          |
| `services.endpoints.status_code`    | 整数  | 响应体期望的 HTTP 状态码（默认 `200`） | ✖️ | 默认 `200`                       |
| `services.endpoints.response_regex` | 字符串 | 响应体内容的正则表达式匹配             | ✖️ |                                |
| `notifications`                     | 对象  | 通知配置                      | ✖️ | 详见 [自定义通知](#自定义通知)             |

下面是一个示例配置文件：

```yaml
display_num: 72
timeout: 5
max_retry_times: 2
max_log_days: 3
cert_notify_days: 7
services:
  - name: "GitHub API"
    endpoints:
      - url: "https://api.github.com"
      - url: "https://api.github.com/repos/wcy-dt/ponghub"
        method: "GET"
        headers:
          Content-Type: application/json
          Authorization: Bearer your_token
        status_code: 200
        response_regex: "full_name"
  - name: "Example Website"
    endpoints:
      - url: "https://example.com/health"
        response_regex: "status"
      - url: "https://example.com/status"
        method: "POST"
        body: '{"key": "value"}'
```

### 特殊参数

ponghub 现已支持强大的参数化配置功能，允许在配置文件中使用多种类型的动态变量，这些变量会在程序运行时实时生成和解析。

<details>

<summary>点击展开查看支持的参数类型</summary>

<div markdown="1">

#### 📅 日期时间参数

使用 `{{%格式}}` 格式定义日期时间参数：

- `{{%Y-%m-%d}}` - 当前日期，格式：2006-01-02（如：2025-09-22）
- `{{%H:%M:%S}}` - 当前时间，格式：15:04:05（如：17:30:45）
- `{{%s}}` - Unix时间戳（如：1727859600）
- `{{%Y}}` - 当前年份（如：2025）
- `{{%m}}` - 当前月份，格式：01-12
- `{{%d}}` - 当前日期，格式：01-31
- `{{%H}}` - 当前小时，格式：00-23
- `{{%M}}` - 当前分钟，格式：00-59
- `{{%S}}` - 当前秒数，格式：00-59
- `{{%B}}` - 完整月份名称（如：September）
- `{{%b}}` - 简短月份名称（如：Sep）
- `{{%A}}` - 完整星期名称（如：Monday）
- `{{%a}}` - 简短星期名称（如：Mon）

#### 🎲 随机数参数

- `{{rand}}` - 生成0-1000000范围的随机数
- `{{rand_int}}` - 生成大范围随机整数
- `{{rand(min,max)}}` - 生成指定范围的随机数
  - 示例：`{{rand(1,100)}}` - 生成1-100之间的随机数
  - 示例：`{{rand(1000,9999)}}` - 生成4位随机数

#### 🔤 随机字符串参数

- `{{rand_str}}` - 生成8位随机字符串（字母+数字）
- `{{rand_str(length)}}` - 生成指定长度的随机字符串
  - 示例：`{{rand_str(16)}}` - 生成16位随机字符串
- `{{rand_str_secure}}` - 生成16位加密安全的随机字符串
- `{{rand_hex(length)}}` - 生成指定长度的十六进制随机字符串
  - 示例：`{{rand_hex(8)}}` - 生成8位十六进制字符串
  - 示例：`{{rand_hex(32)}}` - 生成32位十六进制字符串

#### 🆔 UUID参数

- `{{uuid}}` - 生成标准UUID（带连字符）
  - 示例：`bf3655f7-8a93-4822-a458-2913a6fe4722`
- `{{uuid_short}}` - 生成短UUID（无连字符）
  - 示例：`14d44b7334014484bb81b015fb2401bf`

#### 🌍 环境变量参数

- `{{env(变量名)}}` - 读取环境变量的值
  - 示例：`{{env(API_KEY)}}` - 读取API_KEY环境变量
  - 示例：`{{env(VERSION)}}` - 读取VERSION环境变量
  - 如果环境变量不存在，返回空字符串

环境变量可通过 GitHub Actions 的 Repository Secrets 设置

#### 📊 序列号和哈希参数

- `{{seq}}` - 基于当前时间的序列号（6位数字）
- `{{seq_daily}}` - 每日序列号（自午夜起的秒数）
- `{{hash_short}}` - 短哈希值（6位十六进制）
- `{{hash_md5_like}}` - MD5风格的长哈希值（32位十六进制）

#### 🌐 网络和系统信息参数

- `{{local_ip}}` - 获取系统本地IP地址
- `{{hostname}}` - 获取系统主机名
- `{{user_agent}}` - 生成随机的User-Agent字符串
- `{{http_method}}` - 生成随机的HTTP方法（GET、POST、PUT、DELETE等）

#### 🔐 编码和解码参数

- `{{base64(内容)}}` - 对提供的内容进行Base64编码
  - 示例：`{{base64(hello world)}}` - 将"hello world"编码为Base64
- `{{url_encode(内容)}}` - 对提供的内容进行URL编码
  - 示例：`{{url_encode(hello world)}}` - 对"hello world"进行URL编码
- `{{json_escape(内容)}}` - 对提供的内容进行JSON转义
  - 示例：`{{json_escape("test")}}` - 转义引号和特殊字符以用于JSON

#### 🔢 数学运算参数

- `{{add(a,b)}}` - 两数相加
  - 示例：`{{add(10,5)}}` - 返回15
- `{{sub(a,b)}}` - 两数相减
  - 示例：`{{sub(10,5)}}` - 返回5
- `{{mul(a,b)}}` - 两数相乘
  - 示例：`{{mul(10,5)}}` - 返回50
- `{{div(a,b)}}` - 两数相除
  - 示例：`{{div(10,5)}}` - 返回2

#### 📝 文本处理参数

- `{{upper(文本)}}` - 将文本转换为大写
  - 示例：`{{upper(hello)}}` - 返回"HELLO"
- `{{lower(文本)}}` - 将文本转换为小写
  - 示例：`{{lower(HELLO)}}` - 返回"hello"
- `{{reverse(文本)}}` - 反转文本
  - 示例：`{{reverse(hello)}}` - 返回"olleh"
- `{{substr(文本,起始位置,长度)}}` - 从文本中提取子字符串
  - 示例：`{{substr(hello world,0,5)}}` - 返回"hello"

#### 🎨 颜色生成参数

- `{{color_hex}}` - 生成随机的十六进制颜色代码
  - 示例：`#FF5733`
- `{{color_rgb}}` - 生成随机的RGB颜色值
  - 示例：`rgb(255, 87, 51)`
- `{{color_hsl}}` - 生成随机的HSL颜色值
  - 示例：`hsl(120, 50%, 75%)`

#### 📁 文件和MIME类型参数

- `{{mime_type}}` - 生成随机的MIME类型
  - 示例：`application/json`、`image/png`、`text/html`
- `{{file_ext}}` - 生成随机的文件扩展名
  - 示例：`.jpg`、`.pdf`、`.txt`

#### 👤 虚拟数据生成参数

- `{{fake_email}}` - 生成逼真的虚拟邮箱地址
  - 示例：`john.smith@example.com`
- `{{fake_phone}}` - 生成虚拟电话号码
  - 示例：`+1-555-0123`
- `{{fake_name}}` - 生成虚拟人名
  - 示例：`张三`
- `{{fake_domain}}` - 生成虚拟域名
  - 示例：`example-site.com`

#### ⏰ 时间计算参数

- `{{time_add(时长)}}` - 在当前时间基础上增加指定时长
  - 示例：`{{time_add(1h)}}` - 在当前时间上增加1小时
  - 示例：`{{time_add(30m)}}` - 在当前时间上增加30分钟
  - 支持的单位：s（秒）、m（分钟）、h（小时）、d（天）
- `{{time_sub(时长)}}` - 在当前时间基础上减去指定时长
  - 示例：`{{time_sub(1d)}}` - 在当前时间上减去1天
  - 示例：`{{time_sub(2h30m)}}` - 在当前时间上减去2小时30分钟

</div>
</details>

下面是一个示例配置文件：

```yaml
services:
  - name: "参数化服务"
    endpoints:
        - url: "https://api.example.com/data?date={{%Y-%m-%d}}&rand={{rand(1,100)}}"
        - url: "https://api.example.com/submit"
          method: "POST"
          headers:
            Content-Type: application/json
            X-Request-ID: "{{uuid}}"
          body: '{"session_id": "{{rand_str(16)}}", "timestamp": "{{%s}}"}'
```

### 自定义通知

PongHub 现在支持多种通知方式，当服务出现问题或证书即将过期时，可以通过多个渠道发送警报通知。

<details>

<summary>点击展开查看支持的通知类型</summary>

<div markdown="1">

PongHub 支持以下通知方式：

- **默认通知** - 通过GitHub Actions工作流失败进行通知
- **邮件通知** - 通过SMTP发送邮件，支持高级安全选项
- **自定义Webhook** - 发送到任意HTTP端点，支持高级配置

使用时，在 `config.yaml` 文件中添加 `notifications` 配置块：

```yaml
notifications:
  enabled: true  # 启用通知功能
  methods:       # 要启用的通知方式
    - email
    - webhook
  
  # 各种通知方式的具体配置...
```

#### ⚙️ 默认通知

默认情况下，PongHub 会在 GitHub Actions 工作流失败时发送通知。

默认通知会在以下情况自动启用：

- 没有配置 `notifications` 字段
- `notifications.enabled: true` 但没有指定 `methods` 或仅指定了非email/webhook方法
- 显式配置 `methods: ["default"]`

如果 `notifications` 配置了 `email` 或 `webhook` 方法，默认通知默认关闭，除非在 `notifications.default.enabled` 中明确启用。

#### 📧 邮件通知

```yaml
email:
  smtp_host: "smtp.gmail.com"       # SMTP服务器地址
  smtp_port: 587                    # SMTP端口，默认587
  from: "alerts@yourdomain.com"     # 发件人邮箱
  to:                               # 收件人列表
    - "admin@yourdomain.com"
    - "ops@yourdomain.com"
  subject: "PongHub Service Alert"  # 邮件主题（可选）
  use_tls: true                     # 使用TLS加密（可选）
  use_starttls: true                # 使用STARTTLS（可选）
  skip_verify: false                # 跳过TLS证书验证（可选）
  timeout: 30                       # 连接超时时间，单位秒（可选）
  username: ""                      # SMTP用户名（可选，留空则使用环境变量）
  password: ""                      # SMTP密码（可选，留空则使用环境变量）
  template: ""                      # 自定义邮件模板路径（可选）
  is_html: true                     # 发送HTML邮件（可选）
```

所需环境变量：

- `SMTP_USERNAME` - SMTP用户名
- `SMTP_PASSWORD` - SMTP密码

#### 🔗 自定义Webhook配置

```yaml
webhook:
  url: "https://your-webhook-endpoint.com/notify"  # 留空则从环境变量读取
  method: "POST"                        # HTTP方法（可选，默认POST）
  headers:                              # 自定义请求头（可选，支持特殊参数）
    Content-Type: "application/json"
    Authorization: "Bearer {{env(API_TOKEN)}}"  # 支持特殊参数
    X-Request-ID: "{{uuid}}"           # 使用特殊参数的动态请求ID
  
  # 认证选项（所有字段均支持特殊参数）
  auth_type: "bearer"                   # 认证类型：bearer, basic, apikey（可选）
  auth_token: "{{env(WEBHOOK_TOKEN)}}"  # 使用环境变量的Bearer令牌（可选）
  auth_username: "user-{{rand(1000,9999)}}"  # 带动态后缀的基本认证用户名（可选）
  auth_password: "{{env(AUTH_PASSWORD)}}"     # 来自环境变量的基本认证密码（可选）
  auth_header: "X-API-Key"              # API密钥自定义头部名称（可选）
  
  # 请求配置
  timeout: 30                           # 请求超时时间，单位秒（可选，默认30）
  retries: 3                            # 重试次数（可选，默认0）
  skip_tls_verify: false                # 跳过TLS证书验证（可选）
  
  # 高级载荷自定义，支持特殊参数
  custom_payload:                       # 自定义请求载荷（可选）
    template: |
      {
        "alert": "{{.Title}}",
        "details": "{{.Message}}",
        "timestamp": "{{%Y-%m-%d %H:%M:%S}}",
        "request_id": "{{uuid}}",
        "environment": "{{env(ENVIRONMENT)}}",
        "random_id": "{{rand(10000,99999)}}"
      }
    content_type: "application/json"    # 载荷内容类型（可选）
    fields:                             # 支持特殊参数的额外字段（可选）
      environment: "prod-{{rand(100,999)}}"
      session_id: "{{uuid_short}}"
      build_number: "{{env(BUILD_NUMBER)}}"
      timestamp_unix: "{{%s}}"
    include_title: true                 # 在额外字段中包含标题（可选）
    include_message: true               # 在额外字段中包含消息（可选）
    title_field: "alert_title"          # 标题自定义字段名（可选）
    message_field: "alert_message"      # 消息自定义字段名（可选）
```

**Webhook中的特殊参数支持：**

Webhook配置现在全面支持在以下字段中使用特殊参数：

- **URL**: `url: "https://hooks.example.com/{{env(HOOK_ID)}}"`
- **请求头**: 所有请求头值都可以使用特殊参数
- **认证信息**: 所有认证字段均支持动态值
- **模板**: 同时支持Go模板语法（`{{.Title}}`）和特殊参数（`{{uuid}}`）
- **自定义字段**: 所有自定义载荷字段均支持特殊参数

**模板语法兼容性：**

Webhook模板系统无缝支持两种语法：

- **Go模板语法**: `{{.Title}}`、`{{.Message}}` - 访问通知数据
- **特殊参数**: `{{uuid}}`、`{{%Y-%m-%d}}`、`{{env(VAR)}}` - 动态值

结合两种语法的示例：

```yaml
custom_payload:
  template: |
    {
      "service_alert": "{{.Title}}",
      "description": "{{.Message}}",
      "alert_id": "{{uuid_short}}",
      "timestamp": "{{%Y-%m-%d %H:%M:%S}}",
      "environment": "{{env(DEPLOY_ENV)}}",
      "correlation_id": "{{rand_str(12)}}"
    }
  fields:
    datacenter: "{{env(DATACENTER)}}"
    version: "{{env(APP_VERSION)}}"
```

所需环境变量：

- `WEBHOOK_URL` - 自定义Webhook URL（如果`url`字段为空）
- 特殊参数中引用的任何环境变量（如：`API_TOKEN`、`ENVIRONMENT`）

</div>
</details>

以上所需的环境变量均可通过 GitHub Actions 的 Repository Secrets 设置。

下面是一个完整的通知配置示例文件：

```yaml
display_num: 72
timeout: 5
max_retry_times: 2
max_log_days: 3
cert_notify_days: 7

services:
  - name: "Example Service"
    endpoints:
      - url: "https://example.com/health"

notifications:
  enabled: true
  methods:
    - email
  
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    from: "alerts@yourdomain.com"
    to:
      - "admin@yourdomain.com"
      - "ops@yourdomain.com"
    subject: "🚨 PongHub 服务告警"
    use_starttls: true
    is_html: true
```

## 本地开发

本项目使用 Makefile 进行本地开发和测试。你可以使用以下命令在本地运行项目：

```bash
make run
```

项目有一些测试用例，可以通过以下命令运行测试：

```bash
make test
```

## 免责声明

[PongHub](https://github.com/WCY-dt/ponghub) 仅用于个人学习和研究，不对程序的使用行为或结果负责。请勿将其用于商业用途或非法活动。

<div align="center">
  <img src="imgs/sponsor.png" alt="Sponsor PongHub"/>
  <p>Made with ❤️ by the PongHub Team</p>
</div>
