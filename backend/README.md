# 包说明

## 1. main.go
- **程序入口**：用于初始化各个模块并启动服务。

## 2. api/
- **API 层**：
  - 处理所有的 HTTP 请求和响应。
  - 定义与前端交互的 RESTful API 路由，处理用户请求并将请求转发到服务层。
  - 定义了三个接口：
    - `auth`：认证
    - `history`：查询历史记录
    - `generate`：生成图

## 3. models/
- **数据模型层**：
  - 定义数据库中使用的所有实体模型和结构体。
  - 该层负责将数据库中的数据映射到 Go 语言的结构体中（ORM）。
  - **示例文件**：
    - `init`：初始化
    - `image`：建 image 表
    - `user`：用户表

## 4. repository/
- **数据持久化层**：
  - 负责与数据库的交互，进行数据的存储、更新、删除、查询等操作。
  - 该层会调用模型，并根据业务需求编写具体的数据库查询逻辑。

## 5. services/
- **业务逻辑层**：
  - 该层包含系统的核心业务逻辑。
  - 接收来自 `api/` 层的请求，调用 `repository/` 层来操作数据库，并返回结果。
  - **示例文件**：
    - 对应 `api` 层里的三个接口，在这里实现其逻辑

## 6. config/
- **配置层**：
  - 存储应用程序的配置文件（如数据库连接信息、API 密钥等）。
  - 负责读取和管理项目的全局配置参数。
  - **备注**：方便修改程序信息，昨晚建漏了这个包，没事全部写完再加也可以哈哈

## 7. utils/
- **工具层**：
  - 提供常用的工具类函数和公共方法，这些工具通常会被多个模块调用。
  - 该层可以包含日志工具、错误处理工具、加密工具等。
  - **示例文件**：
    - `jwt`：生成密码的 hash
    - `response`：响应处理

## 8. middleware/
- **中间件层**：
  - 处理请求生命周期中的中间件逻辑，如认证、日志、跨域请求处理等。
  - 中间件通常用于拦截和处理 HTTP 请求的预处理。
  - **示例文件**：
    - `jwt`：登录要用到的登录验证中间件，他会返回一个 token 用于身份验证登录

## 9. docs/
- **文档**：
  - 项目的文档存放处，包含 API 文档、架构图、使用说明等。
  - 该目录帮助开发者和使用者理解项目的功能和使用方法。
  - **备注**：我们好像没有，先建着吧

---
# 示例接口规范
### 创建用户

### 功能

创建一个新用户

#### URL地址

`POST localhost:8080/api/v1/users`



#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体

```json
{
  "username": "string",
  "password": "string",
  "email": "user@example.com"
}
```

#### 响应

| 响应码 | 描述         | 示例响应体                                            |
| ------ | ------------ | ----------------------------------------------------- |
| 201    | 用户创建成功 | `{"message": "User created successfully","code":201}` |
| 400    | 请求参数错误 | `{"error": "Invalid request parameters"}`             |
| 409    | 用户名已存在 | `{"error": "Username already exists"}`                |
---
### 开发内容
#### 管理员账号：
- username:root
- password:111111(加密后"bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a")
1. **登陆注册接口**
#### 注册数据格式：
- 注册界面接收前端传来邮箱，进行数据库查询，判断用户是否存在，不存在则注册，存在则返回错误信息，并对密码进行加密，保存到数据库中
- 登陆界面接收前端传来用户名和密码，进行数据库查询，判断用户是否存在，存在则比对密码，不存在则返回错误信息
- jwt返回token用于登录验证
#### url：
POST http://localhost:8080/register


```json
{
  "email": "root@qq.com",
  "username": "root1",
  "password": "sssssss"
}
```
#### 响应

| 响应码 | 描述                   | 示例响应体                                            |
|--------|------------------------|-------------------------------------------------------|
| 400    | 请求数据格式错误         | { "code": 400, "message": "请求数据格式错误" }         |
| 500    | 用户创建失败（具体错误信息） | { "code": 500, "message": "用户创建失败", "error": "邮箱已存在" } |
| 200    | 注册成功                | { "code": 200, "message": "注册成功" }                 |

---

#### 登录
- 使用用户名和密码登录
#### 请求头
无
#### 请求体
``` json
{
  "username": "root1",
  "password": "sssssss"
}
```
#### 响应

| 响应码 | 描述                    | 示例响应体                                                     |
|--------|-------------------------|----------------------------------------------------------------|
| 400    | 请求数据格式错误          | { "code": 400, "message": "请求数据格式错误" }                  |
| 401    | 用户不存在               | { "code": 401, "message": "用户不存在" }                       |
| 500    | 数据库查询错误            | { "code": 500, "message": "数据库查询错误" }                    |
| 401    | 密码错误                 | { "code": 401, "message": "密码错误" }                          |
| 500    | 生成token错误             | { "code": 500, "message": "生成 token 错误" }                   |
| 500    | 登录时更新用户token失败     | { "code": 500, "message": "登录时更新用户 token 失败", "error": "查询用户信息失败" } |
| 200    | 登录成功并返回token       | { "code": 200, "message": "登录成功", "token": "fake-jwt-token" } |
---

2. **文生图接口**
生成图像并记录积分消耗,每次消耗20积分

#### URL地址
- (POST) http://localhost:8080/auth/generate

#### 请求参数
- 请求头携带一个"Authorization"的token
- 请求体（JSON格式）：

```json
{
  "prompt": "string",
  "width": 220,
  "height": 200,
  "steps": 100,
  "sampling_method": "DDIM",
  "seed": "string"
}
```

#### 响应

| 响应码 | 描述                 | 示例响应体                                                                  |
|--------|----------------------|-----------------------------------------------------------------------------|
| 200    | 图片生成成功并返回URL | { "code": 200, "image_url": "New_Image_Url", "message": "用户当前积分为", "success": true } |
| 400    | 缺少提示词           | { "code": 400, "message": "缺乏提示词", "success": false }                   |
| 400    | 宽度不在范围内       | { "code": 400, "message": "宽度不在范围内", "success": false }               |
| 400    | 高度不在范围内       | { "code": 400, "message": "高度不在范围内", "success": false }               |
| 400    | 步数不在范围内       | { "code": 400, "message": "步数不在范围内", "success": false }               |
| 400    | 采样方法不在范围内   | { "code": 400, "message": "采样方法不在范围内", "success": false }           |
| 400    | 缺少种子             | { "code": 400, "message": "缺乏种子", "success": false }                     |
| 401    | 请求头中缺少Token     | { "code": 401, "message": "请求头中缺少Token", "success": false }             |
| 401    | 无效的Token          | { "code": 401, "message": "无效的Token", "success": false }                  |
| 401    | 未找到用户信息       | { "code": 401, "message": "未找到用户信息", "success": false }                |
| 401    | 用户信息查询失败     | { "code": 401, "message": "用户信息查询失败", "success": false }              |
| 401    | 用户积分不足         | { "code": 401, "message": "用户积分不足", "success": false }                 |
| 401    | 用户积分更新失败     | { "code": 401, "message": "用户积分更新失败", "success": false }             |
| 401    | 积分记录创建失败     | { "code": 401, "message": "积分记录创建失败", "success": false }             |
| 500    | 图片生成失败         | { "code": 500, "message": "图片生成失败", "success": false }                 |
---

3. **个人信息界面**
  - 结合数据库的用户信息，使用查询函数查询出需要的信息返回给前端
  - 头像上传功能，获取功能set，get
### 修改头像URL

#### URL地址
- (POST) http://localhost:8080/auth/setavator

#### 请求参数
- 必须携带一个"Authorization"的token，和头像URL：

```json
{
  "url": "string"  // 更换头像的URL
}
```

#### 响应

| 响应码 | 描述               | 示例响应体                                    |
|--------|--------------------|-----------------------------------------------|
| 401    | 请求头中缺少Token   |` { "code": 401, "message": "请求头中缺少Token" }` |
| 401    | 无效的Token        | `{ "code": 401, "message": "无效的Token" }  `     |
| 500    | 更新头像失败        | `{ "code": 500, "message": "更新头像失败" }  `    |
| 200    | 头像更新成功        | `{ "code": 200, "msg": "头像更新成功", "data": "newURL" } `|

---

### 获取头像URL

#### URL地址
- (GET) http://localhost:8080/auth/getavator

#### 请求参数
- 必须携带一个"Authorization"的token。

#### 响应

| 响应码 | 描述            | 示例响应体                               |
|--------|-----------------|------------------------------------------|
| 401    | Token已过期     |` { "code": 401, "message": "Token已过期" }` |
| 401    | 无效的Token     | `{ "code": 401, "message": "无效的Token" } `|
| 500    | 查询头像失败     | `{ "code": 500, "message": "查询头像失败" } `|
| 200    | 获取头像成功     | `{ "code": 200, "msg": "获取头像成功", "data": "user.Avator_url" }` |

---

### 查询当前登录用户信息

#### URL地址
- (GET) http://localhost:8080/auth/user/info

#### 请求参数
- 请求头携带一个"Authorization"的token。

#### 响应

| 响应码 | 描述               | 示例响应体                                                                                                         |
|--------|--------------------|--------------------------------------------------------------------------------------------------------------------|
| 400    | 请求数据格式错误    | `{ "code": 400, "message": "Invalid request data" }          `                                                        |
| 404    | 用户未找到         | `{ "code": 404, "message": "用户未找到" }           `                                                                 |
| 500    | 查询失败           | `{ "code": 500, "message": "查询失败", "error": "错误信息" }   `                                                       |
| 200    | 查询成功           | `{ "code": 200, "user": { "id": 6, "email": "czh@qq.com", "username": "czh", "avatar_url": "https://www.chen.com", "create_time": "2024-11-24T21:49:24.78802Z" } }` |

---

### 查询所有用户信息

#### URL地址
- (GET) http://localhost:8080/user/all

#### 请求参数
- 无（GET方法，不需要携带token）

#### 响应

| 响应码 | 描述               | 示例响应体                                    |
|--------|--------------------|-----------------------------------------------|
| 500    | 获取用户列表失败    | `{ "code": 500, "message": "获取用户列表失败" } `|
| 200    | 获取用户列表成功    | `{ "code": 200, "message": "获取用户列表成功", "users": [{ ... }, { ... }] } `|

---

### 更新当前登录用户信息

#### URL地址
- (PUT) http://localhost:8080/auth/user/update

#### 请求参数
- 请求头携带一个"Authorization"的token，参数包含：

```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "avator_url": "string",
  "token": "string",
  "create_time": "string"
}
```

- 所有参数都是可选的，但用户名和ID不可修改。

#### 响应

| 响应码 | 描述                 | 示例响应体                                                                   |
|-----|----------------------|-------------------------------------------------------------------------|
| 400 | 请求数据格式错误      | `{ "code": StatusBadRequest(400),"message":"请求数据格式错误","error": "错误信息"}` 
| 401 | 未找到用户信息     | `{ "code": 401, "message": "未找到用户信息" }`                                 |
| 200 | 用户信息更新成功     | `{ "code": 200, "message": "用户信息更新成功" }`               |
---


4. **文生图历史记录**
#### 获取所有用户所有的图像信息

- **URL地址**  
  `(GET) http://localhost:8080/image/all`

- **请求参数**  
  无（GET方法，不需要携带token）

- **响应**

| 响应码 | 描述                    | 示例响应体                                                                 |
|-------|-----------------------|-------------------------------------------------------------------------|
| 500   | 获取图像列表失败            | `{ "code": StatusInternalServerError (500), "message": "获取图像列表失败" }` |
| 500   | 查询失败                  | `{ "code": StatusInternalServerError (500), "message": "查询失败", "images": [...] }` |

---

#### 获取当前登录用户在一段时间内生成的图像信息

- **URL地址**  
  `(GET) localhost:8080/auth/user/images/timeRange`

- **请求参数**  
  `?start_time=YYYY-MM-DD&end_time=YYYY-MM-DD`  
  或 `?start_time=2006-01-02T15:04:05.000000Z&end_time=2006-01-02T15:04:05.000000Z`

- **请求头**
  ```json
  {
    "Authorization": "your_jwt_token"
  }
  ```
  
| 响应码 | 描述                   | 示例响应体                                                                   |
|-------|----------------------|---------------------------------------------------------------------------|
| 400   | 无效的开始时间格式       | `{ "code": StatusBadRequest (400), "message": "无效的开始时间格式", "error": "" }` |
| 400   | 无效的结束时间格式       | `{ "code": StatusBadRequest (400), "message": "无效的结束时间格式", "error": "" }` |
| 500   | 查询图像列表失败        | `{ "code": StatusInternalServerError (500), "message": "查询图像列表失败", "error": "" }` |
| 200   | 查询图像列表成功        | `{ "code": StatusOK, "message": "查询图像列表成功", "images": [...] }`      |
---
### 获取指定的某张图像

- **URL地址**
  `(GET) http://localhost:8080/image`

- **请求参数**
  `?username=` 或 `?id=` 或 `?url=`

- **响应**

| 响应码 | 描述                 | 示例响应体                                                                   |
|-------|--------------------|---------------------------------------------------------------------------|
| 404   | 未找到相关图片          | `{ "code": StatusNotFound (404), "message": "未找到相关图片" }`             |
| 500   | 查询用户的图片失败       | `{ "code": StatusInternalServerError (500), "message": "查询用户的图片失败", "error": "" }` |
| 400   | 无效的图像id或用户名     | `{ "code": StatusBadRequest (400), "message": "无效的图像id或用户名", "error": "" }` |
| 200   | 查询图像成功            | `{ "code": StatusOK, "message": "查询图像成功", "image": { ... } }`         |
---

### 获取当前登录用户生成的所有图像

- **URL地址**
  `(GET) http://localhost:8080/auth/user/images`

- **请求头**
  ```json
  {
    "Authorization": "your_jwt_token"
  }
  ```

- **响应**

| 响应码 | 描述                | 示例响应体                                                                  |
|-------|-------------------|--------------------------------------------------------------------------|
| 500   | 查询用户图像失败      | `{ "code": StatusInternalServerError (500), "message": "查询用户图像失败", "error": "" }` |
| 200   | 获取用户的图像成功     | `{ "code": StatusOK, "message": "获取用户的图像成功", "images": [...] }`    |
---

### 点赞图片功能

- **URL地址**
  `(POST) http://localhost:8080/auth/like`

- **请求头**
  ```json
  {
    "Authorization": "your_jwt_token"
  }
  ```

- **请求体**
  ```json
  {
    "url": "图像url"
  }
  ```

- **响应**

| 响应码 | 描述                   | 示例响应体                                                                        |
|-------|----------------------|------------------------------------------------------------------------------|
| 400   | 缺少图像URL             | `{ "code": 400, "error": "Missing image URL" }`                              |
| 409   | 用户已经点赞过该图片      | `{ "code": 409, "error": "用户已经点赞过该图片" }`                                     |
| 500   | 返回获取赞数错误        | `{ "code": 500, "error": "返回获取赞数错误的error" }`                                 |
| 500   | 点赞数据库操作出错       | `{ "code": 500, "error": "点赞数据库操作出错" }`                                      |
| 200   | 点赞成功                | `{ “code":200,current_likes": 当前赞数, "message": "Image liked successfully" }` |
---
  
5. **图片收藏界面**
#### 查询功能
- 查询展示出用户的收藏图片
#### URL地址

`(GET) http://localhost:8080/auth/user/favoritedimages`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
无

#### 响应
- 响应格式：同localhost:8080/user/images（只不过message多了一个“收藏”）

#### 收藏图像功能
- 收藏指定图像
#### URL地址

`(POST) http://localhost:8080/auth/addFavoritedImage`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
//两个参数有一个就行
```json
{
   
  "url":"",
  "id":""
}
```

#### 响应
| 响应码 | 描述                | 示例响应体                                                                                    |
|-----|-------------------|------------------------------------------------------------------------------------------|
| 400 | 无有效的图像id或url      | `{ code: StatusBadRequest (400),message: "无有效的图像id或url",error: "id 必须大于 0 或者 url 不得为空"}` |
| 404 | 未找到对应的图像          | `{ "code"：StatusUnauthorized(404),"message"："未找到对应的图像",error}`                           |
| 401 | 未找到用户信息"或没有token时 | `{ "code"：StatusUnauthorized(401),"message"："未找到用户信息"或没有token时",error}`                  |
| 500 | 检查收藏状态失败        | `{ "code"：StatusInternalServerError(500),"message"："检查收藏状态失败",error}`               |
| 401 | 该图像已经被收藏过          | `{ "code"：StatusConflict(409),"message"："该图像已经被收藏过",error}`                         |
| 500 | 收藏图像失败        | `{ "code"：StatusInternalServerError(500),"message"："收藏图像失败",error}`               |
| 200 | 收藏图像成功          | `{ "code"：StatusOK(200),"message"："收藏图像成功"}`                                 |
---

### 取消图像收藏
- 取消指定图像的收藏

#### URL地址

`(DELETE) localhost:8080/auth/deleteFavoritedImage`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
//至少传递一个参数
```json
   
  "url": "",
  "id": ""
```

#### 响应

| 响应码 | 描述                | 示例响应体                                                                                    |
|-----|-------------------|------------------------------------------------------------------------------------------|
| 400 | 无有效的图像id或url      | `{ code: StatusBadRequest (400),message: "无有效的图像id或url",error: "id 必须大于 0 或者 url 不得为空"}` |
| 404 | 未找到对应的图像          | `{ "code"：StatusNotFound(404),"message"："未找到对应的图像",error}`                           |
| 401 | 未找到用户信息"或没有token时 | `{ "code"：StatusUnauthorized(401),"message"："未找到用户信息"或没有token时",error}`                  |
| 500 | 检查收藏状态失败        | `{ "code"：StatusInternalServerError(500),"message"："检查收藏状态失败",error}`               |
| 409 | 该图像未被收藏过          | `{ "code"：StatusConflict(409),"message"："该图像未被收藏过"}`                         |
| 500 | 取消图像收藏失败        | `{ "code"：StatusInternalServerError(500),"message"："取消图像收藏失败"}`               |
| 200 | 取消图像收藏成功          | `{ "code"：StatusOK(200),"message"："取消图像收藏成功"}`                                 |
---
6. **签到增加积分接口**
#### 功能
 - 用户每次签到增加100积分，并保存记录
#### URL地址

`(GET) http://localhost:8080/auth/score`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
无

#### 响应

| 响应码 | 描述          | 示例响应体                                                                  |
|-----|-------------|------------------------------------------------------------------------|
| 401 | 请求头中缺少Token | `{ "code"：StatusUnauthorized(401),"message"："请求头中缺少Token",data: null}` |
| 401 | 无效的Token    | `{ "code"：StatusUnauthorized(401),"message"："无效的Token",data: null}`    |
| 401 | 用户信息查询失败    | `{ "code"：StatusUnauthorized(401),"message"："用户信息查询失败",data: null}`      |
| 401 |   积分记录创建失败          | `{ "code"：StatusUnauthorized(401),"message"："积分记录创建失败",data: null}`        |
| 401 |  用户积分更新失败          | `{ "code"：StatusUnauthorized(401),"message"："用户积分更新失败",data: null}`        |
| 200 | 返回用户当前积分    | `{ "code"：StatusOK(200),"message"："用户当前积分为",data: null}`               |
---
7. **token校验功能**
#### 功能
    - 校验用户的token
#### URL地址

`(GET) http://localhost:8080/checkToken`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
无

#### 响应

| 响应码 | 描述            | 示例响应体                                                               |
| ------ | ------------- |---------------------------------------------------------------------|
| 401    | 令牌格式不正确 | `{ "code"：StatusUnauthorized(401),"message"："令牌格式不正确",data: null}`  |
| 401    | 令牌过期或未激活     | `{ "code"：StatusUnauthorized(401),"message"："令牌过期或未激活",data: null}` |
| 401    | 令牌无法处理    | `{ "code"：StatusUnauthorized(401),"message"："令牌无法处理",data: null}`   |
| 500    | 令牌无效 | `{ "code"：StatusUnauthorized(401),"message"："令牌无效",data: null}`     |
| 200    | 令牌有效 | `{ "code"：StatusOK(200),"message"："令牌有效",data: null}`               |
---

8. **搜索功能**
### 搜素图像
### 功能
  根据前端传来的关键字（可多个）查询Prompt中包含至少一个关键字的图像。
  如果?isOwn=true，则只查询当前登录用户的。

#### URL地址

`(DELETE) http://localhost:8080/auth/image/feature`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
无。
通过查询参数:   ?feature=……&feature=……（可多个feature）（可选&isOwn=true ）

#### 响应

| 响应码 | 描述               | 示例响应体                                                                              |
| ------ | ----------------- | -------------------------------------------------------------------------------------- |
| 401    | 请求头中缺少Token" | `{"message"："请求头中缺少Token"}`                                                       |
| 401    | 无效的Token       | `{"message"："无效的Token"}`                                                             |
| 401    | 未找到用户信息     | `{"message"："未找到用户信息"}`                                                          |
| 500    | 根据关键字查询失败 | `{"message"："根据关键字查询失败"}`                                                       |
| 200    | 根据关键字查询成功 | `{"images": [{"id":, "username":, "params":, "picture":, "likecount":, "create_time":,},{……},…………]}`  |
---
  

9. **管理员操作**
### 删除用户
### 功能
  根据前端传来的?username删除指定用户
  如果?isOwn=true表示账号注销（即删除自己的账号信息，如果没有这个，则表示root用户删除违规账号）

#### URL地址

`(DELETE) http://localhost:8080/auth/root/deleteOneUser`

#### 请求头

```json
{
  "Authorization": "your_jwt_token"
}
```

#### 请求体
无。
通过查询参数:   ?username= (可选isOwn=true)

#### 响应

| 响应码 | 描述               | 示例响应体                                        |
| ------ | ----------------- | ------------------------------------------------- |
| 401    | 请求头中缺少Token  | `{"message"："请求头中缺少Token"}`                 |
| 401    | 无效的Token        | `{"message"："无效的Token"}`                      |
| 401    | 未找到用户信息      | `{"message"："未找到用户信息"}`                   |
| 400    | 非root用户         | `{"message"："非root用户，不可删除其他某个用户"}`   |
| 401    | 无效的Token        | `{"message"："无效的Token"}`                      |
| 500    | 查询用户是否存在失败| `{"message"："查询用户是否存在失败"}`              |
| 400    | 用户不存在         | `{"message"："用户不存在"}`                       |
| 500    | 删除用户失败       | `{"message"："删除用户失败"}`                      |
| 200    | 成功删除用户       | `{"message"："成功删除用户：（用户名）"}`           |
| 200    | 账号注销成功       | `{"message"："（用户名）的账号注销成功"}`           |
---
 

10. **数据库设计**
   - 用户登录表：id，email，user，password，token
   - 用户查询表：id，user（外键），params，picture，time
   - 收藏表：id，user（外键），picture
