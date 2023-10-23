# SQLRest 通用数据网关

SQLRest通用数据网关是一个轻量级的 RESTful 服务，允许用户使用单个 API 从多个数据库中获取数据。通过提供统一的 API 并自动处理不同数据库的连接和查询，它简化了从不同来源获取数据的过程。该服务基于 `Gorilla Mux` 路由器构建，并使用 `database/sql` 包进行数据库交互。

## 功能

- 支持多个数据库源。
- 基于输入参数执行动态查询。
- 方便地添加新的数据源和 API。
- 使用 API 密钥简化身份验证。
- 日志记录和错误处理。

## 入门指南

### 配置

配置设置可以在 `config.yaml` 文件中找到。

```yaml
# Base configuration 
base:
  url: "http://127.0.0.1:8080"

# Database configuration  
database:  
  # Database name  
  name: ""  
  # Database user  
  user: ""  
  # Database password  
  password: ""  
  # Database host (IP address)  
  host: ""  
  # Database port  
  port:   

# Logger parameters  
logger:  
  # Log file path  
  log_filename: "logs/app.log"  
  # Log level (info, debug, etc.)  
  log_level: "debug"  
  # Maximum size of each log file in MB  
  max_size: 10  
  # Maximum number of log files to keep  
  max_backups: 5  
  # Maximum age of log files in days  
  max_age: 28  
  # Whether to compress old log files (true or false)  
  compress: true  

```

### 运行SQLRest
配置config.yaml，导入数据库sql/sqlrest.sql（如果需要使用sqlrest ui，并注意修改db_source的连接记录），启动 go run cmd/main.go

打开：http://localhost:8080/ （sqlrest 管理 ui）

右上角登陆，默认管理Authkey为：101e1c66e7b2e821

sqlrest ui项目 Project Id 为：5f51665a


## 发送请求

要向自动数据 API 发送请求，请发送带有以下参数的 HTTP POST 请求：

- `api_id`：API 的 ID（get参数）。
- `authkey`：API 密钥（header中）


请求示例：

```
POST http://[sqlrest host]/api?api_id=1
Content-Type: application/json
authkey: your_auth_key

{
    "param1": "value1",
    "param2": "value2"
}
```

请求成功时，服务器将返回以下格式的 JSON 响应：

```
{
    "data": [ {},...,{} ],
    "code": 200,
    "msg": "success"
}
```

如果请求失败，服务器将返回错误消息和适当的 HTTP 状态代码。

#### 多API联合查询

- `api_id`：支持以逗号分隔传入多个api_id，例如 2,3 （get参数）

请求示例：

```
POST http://[sqlrest host]/api?api_id=2,3
Content-Type: application/json
authkey: your_auth_key

{
    "param1": "value1",
    "param2": "value2"
}
```

请求成功时，服务器将返回以下格式的 JSON 响应：

```
{
    "data": [ [{},...,{}],...,[{},...,{}]] ],
    "code": 200,
    "msg": "success"
}
```

当传入多API时，data内容会将两次查询结果进行拼接，并返回多个[]

#### API子查询

当传入api子查询id时，子查询数据参数将以主查询返回内容为基础，对返回内容进行逐条子查询，例如返回内容为{"username":"xxx","userid":"yyy"} 时，子查询的传入参数即位username、userid，并且子查询不能继承用户post传入的参数，如需要可以主查询中构造，查询后子查询将以sub_query_x为key值，拼接入数据中，使用时注意效率问题

- `sub_api_id`：支持1个或多个子查询，多个以逗号分隔例如 2,3 （get参数）

请求示例：

```
POST http://[sqlrest host]/api?api_id=1&sub_api_id=2
Content-Type: application/json
authkey: your_auth_key

{
    "param1": "value1",
    "param2": "value2"
}
```

请求成功时，服务器将返回以下格式的 JSON 响应：

```
{
    "data": [{...,"sub_query_0":"xxx"},...,{....,"sub_query_0":"xxx"}] ],
    "code": 200,
    "msg": "success"
}
```

多API联合查询可以与API子查询同时使用

### 通用变量替换

在SQLRest中，变量替换功能允许您在 SQL 查询中使用动态值。为了实现这一功能，您需要在 data_api 表的 sourcesql 字段中使用占位符。占位符的格式为 {param}，其中 param 是参数名称。

当用户向 API 发送请求时，他们可以通过 GET 或 POST 请求提供这些参数。Auto Data API 会将这些参数的值替换为对应的占位符，然后执行生成的 SQL 查询。

假设您有一个用户表，包含以下字段：id、name 和 email，您希望创建一个 API 来根据给定的 name 和 email 过滤用户。您可以在 data_api 表中添加以下记录：
- apiname: "get_users"
- db_id: 1（对应于 db_source 表中的数据库源）
- sourcesql: "SELECT * FROM users WHERE name = '{name}' AND email = '{email}'"
- is_del: 0

您可以使用以下 POST 请求调用此 API：

```
POST http://[sqlrest host]/api?api_id=1  
Content-Type: application/json  
authkey: your_auth_key  
  
{  
	"name": "John",  
	"email": "john@example.com"  
}  
```

在这些请求中，Auto Data API 会将 {name} 和 {email} 占位符替换为实际值，然后执行以下查询：

```sql
SELECT * FROM users WHERE name = 'John' AND email = 'john@example.com'  
```

返回结果将是符合给定条件的用户列表。通过这种方式，您可以轻松地为不同的参数值创建动态查询。

请注意，如果您在 SQL 查询中使用了不存在于请求中的占位符，那么占位符将不会被替换。为了避免潜在的错误，请确保为所有占位符提供相应的参数。

### 内置变量替换

当启用 Joint 身份验证时，SQLRest 可以根据所选的身份验证模式（如 cloud 或 huhang）自动获取用户信息，并将其作为内置变量提供。这些内置变量可以在 SQL 查询中使用，就像使用普通请求参数一样。

在启用 Joint 身份验证的情况下，SQLRest 将根据用户的 token 自动提取以下内置变量：
- userid
- employeeid
- alldeptname
- deptname
- account
- username

这些内置变量可以在 data_api 表中的 sourcesql 字段里使用，与其他请求参数一起作为占位符。只需在占位符名称前加上对应的内置变量名称，如 {userid} 或 {employeeid}。

### 变量限制

不允许自定义变量与内置变量冲突。自定义变量的名称不应与任何内置变量（如 userid、employeeid 等）相同，以避免潜在的冲突和错误。
不允许在请求参数中包含 SQL 关键词，如 SELECT、INSERT、UPDATE、DROP 和 TRUNCATE 等。SQLRest 将检查请求参数，以防止潜在的 SQL 注入攻击。如果发现包含这些关键词的参数，请求将被拒绝。

遵循这些限制可确保 SQLRest 正常运行并提供安全的数据访问。


## 开发相关

### manage-ui

manage-ui为基于sqlrest，vue所开发的管理sqlrest本身的ui界面

### 数据库结构

您需要在主数据库中创建以下表：

- `db_source`：存储有关不同数据库源的信息。
- `data_api`：存储有关不同 API 及其对应 SQL 查询的信息。
- `auth_key`：存储 API 密钥及其允许的 API ID。
- `project`：存储有关不同项目的信息。

参考项目描述中提供的 SQL 命令创建这些表。

## 用法

### 添加新的项目

要添加新的项目，请在 `project` 表中插入一行新数据，包含以下字段：

- `project_name`：项目的唯一名称
- `project_type`：项目类型（可选：'cloud' 或 'huhang'，根据需要启用不同的身份验证模式）
- `memo`：（可选）项目的简短描述
- `is_del`：对于活动项目，请将其设置为 0
- `created_stime`：记录创建时间
- `modified_stime`：记录修改时间
- `created_user`：记录创建者
- `modified_user`：记录修改者

项目表主要用于管理和组织 API。根据项目类型，可以为每个项目启用不同的身份验证模式（如 cloud 或 huhang）。这有助于将 API 分组，并根据需要为其应用不同的身份验证策略。

在添加新的 API 时，您需要在 data_api 表中引用项目的 ID。这将使 API 与项目关联，并允许根据项目设置应用正确的身份验证模式。

### 添加新的数据库源

要添加新的数据库源，请在 `db_source` 表中插入一行新数据，包含以下字段：

- `dbname`：数据库名称。
- `ip`：数据库服务器的 IP 地址。
- `port`：数据库服务器的端口号。
- `user`：数据库连接的用户名。
- `pwd`：数据库连接的密码。
- `memo`：（可选）数据库源的简短描述。
- `is_del`：对于活动源，请将其设置为 `0`。
- `created_stime`：记录创建时间。
- `modified_stime`：记录修改时间。
- `created_user`：记录创建者。
- `modified_user`：记录修改者。

### 添加新的 API

要添加新的 API，请在 `data_api` 表中插入一行新数据，包含以下字段：

- `apiname`：API 的唯一名称。
- `project_id`：项目的 ID（来自 project 表）。
- `db_id`：数据库源的 ID（来自 db_source 表）。
- `sourcesql`：用于从数据库中提取数据的 SQL 查询。对于动态值，请使用占位符（例如 {param}）。
- `memo`：（可选）API 的简短描述。
- `is_del`：对于活动 API，请将其设置为 0。
- `created_stime`：记录创建时间。
- `modified_stime`：记录修改时间。
- `created_user`：记录创建者
- `modified_user`：记录修改者。


### 添加新的 API 密钥

要添加新的 API 密钥，请在 `auth_key` 表中插入一行新数据，包含以下字段：

- `authkey`：API 密钥的唯一字符串。
- `joint`：联合验证类型（可选：'huhang' 或 'cloud'）。
- `project_id`：项目的 ID（来自 project 表）。
- `api_ids`：允许的 API ID 的逗号分隔列表。
- `memo`：（可选）API 密钥的简短描述。
- `is_del`：对于活动 API 密钥，请将其设置为 0。
- `created_stime`：记录创建时间。
- `modified_stime`：记录修改时间。
- `created_user`：记录创建者。
- `modified_user`：记录修改者。