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
  "Authorization": "your_jwt_token",
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
- password:111111(加密后"c4ca4238a0b923820dcc509a6f75849b")
1. **登陆注册接口**
   - 注册界面接收前端传来邮箱，进行数据库查询，判断用户是否存在，不存在则注册，存在则返回错误信息，并对密码进行加密，保存到数据库中
   - 登陆界面接收前端传来用户名和密码，进行数据库查询，判断用户是否存在，存在则比对密码，不存在则返回错误信息
   - jwt返回token用于登录验证
   - 注册访问url：http://localhost:8080/register
   - 注册数据格式：
   - ```json
     {
        “email": "root@qq.com",
        "username": "root1",
        "password": "sssssss"
     }
   - 响应格式：
        - ```json
          code: 400 (StatusBadRequest),
          message:  "请求数据格式错误"
        - ```json
          code: 500 (StatusInternalServerError),
          message: "用户创建失败",
          error:   err.Error()
          (err.Error()可能的情况为：
            1、“名字为空”
            2、“邮箱为空”
            3、“密码少于6位”
            4、“邮箱格式不正确”
            5、“用户名已存在”
            6、“邮箱已存在”
            7、“查询用户名时发生错误”
            8、“查询邮箱时发生错误”
            9、“插入用户表失败”
          )
        - ```json  
          code: 200,
		      message:  "注册成功",

   - 登录访问url：http://localhost:8080/login
   - 登录数据格式：
   - ```json
     {
        "username": "root1",
        "password": "sssssss"
     }
   - 响应格式：
   - ```json
     code：400（StatusBadRequest）,
     message："请求数据格式错误"
   - ```json
     code：401（Unauthorized）,
     message: "用户不存在"
   - ```json
     code：500（StatusInternalServerError）,
     message: "数据库查询错误"
   - ```json
     code：401（Unauthorized）,
     message: "密码错误"
   - ```json
     code：500（StatusInternalServerError）,
     message: "生成 token 错误"
   - ```json
     code：500（StatusInternalServerError）,
     message: "登录时更新用户 token 失败",
     "error": "用户不存在" 或 "查询用户信息失败" 或 "更新用户信息失败"
   - ```json
     code：200（StatusOK）,
     message："登录成功",
     token: "fake-jwt-token"
     


2. **文生图接口**
   - 用户每次生成消耗20积分
   - 部署本地的文生图模型，编写接口进行传参和调用，
   - 接收前端的参数，调用本地部署的大模型，生成对应的图片，返回给前端，并将记录存入数据库
   - 文生图url：http://localhost:8080/auth/generate
   - 参数格式：
   - ```json
     请求头携带一个"Authorization"的token
     参数：
     {
       "prompt": "string",
       "width": 220,
       "height": 200,,
       "steps": 100
       "sampling_method": "DDIM",
       "seed": "string"
     }
    - 响应格式：
   - ```json
     code：200（StatusOK）,
     image_url: "New_Image_Url" 
     message："用户当前积分为",
     success: true
   - ```json
     code：400（StatusBadRequest）,
     message："缺乏提示词",
     success: false
   - ```json
     code：400（StatusBadRequest）,
     message："宽度不在范围内",
		 success: false
   - ```json
     code：400（StatusBadRequest）,
     message："高度不在范围内",
		 success: false
   - ```json
     code：400（StatusBadRequest）,
     success: false,
     message："步数不在范围内"
   - ```json
     code：400（StatusBadRequest）,
     success: false,
     message："采样方法不在范围内"
   - ```json
     code：400（StatusBadRequest）,
     success: false,
     message："缺乏种子"
   - ```json
     code：401（StatusUnauthorized）,
     message："请求头中缺少Token"
   - ```json
     code：401（StatusUnauthorized）,
     message："无效的Token"
   - ```json
     code：401（StatusUnauthorized）,
		 success: false,
     message："未找到用户信息"
   - ```json
     code：401（StatusUnauthorized）,
     success：false,
     message："用户信息查询失败"
   - ```json
     code：401（StatusUnauthorized）,
     success：false,
     message："用户积分不足"
   - ```json
     code：401（StatusUnauthorized）,
     success：false,
     message："用户积分更新失败"
   - ```json
     code：401（StatusUnauthorized）,
     success：false,
     message："积分记录创建失败"
   - ```json
     code：500（StatusInternalServerError）,
     success：false,
     message："图片生成失败"
   - ```json
     code：500（StatusInternalServerError）,
     message："图片生成失败"


3. **个人信息界面**
  - 结合数据库的用户信息，使用查询函数查询出需要的信息返回给前端
  - 头像上传功能，获取功能set，get
    - 修改头像url：http://localhost:8080/auth/setavator
    - 参数格式：
    - ```json
      携带一个"Authorization"的token
      "url": "string"(更换头像的url)
    - 响应格式：
    - ```json
      code: 401（Unauthorized）,
      message:  "请求头中缺少Token"
    - ```json
      code: 401（Unauthorized）,
      message:  "无效的Token"
    - ```json
      code: 500(Error),
      message:  "更新头像失败"
    - ```json
      code: 200(Success),
      msg:  "头像更新成功",
      data: "newURL"
    - 获取头像url：http://localhost:8080/auth/getavator
    - 参数格式：
    - ```json
      携带一个"Authorization"的token
    - 响应格式：
    - ```json
      code: 401（Unauthorized）,
      message:  "Token已过期"
    - ```json
      code: 401（Unauthorized）,
      message:  "无效的Token"
    - ```json
      code: 500 (Error),
      message:  "查询头像失败"
    - ```json  
      code: 200 (Success),
      msg:  "获取头像成功",
      data: "user.Avator_url"  

  - 用户信息查询
  - 查询当前登录用户的信息
    - url: http:localhost:8080/auth/user/info
    - 参数格式 
    - ```json
      请求头携带一个"Authorization"的token
    - 响应格式：
    - ```json
      code: StatusBadRequest (400)
      message: "Invalid request data"
    - ```json
      code: StatusNotFound (404)
      message: "用户未找到"
    - ```json
      code: StatusInternalServerError (500)
      message: "查询失败"
      error: 
    - ```json
      code: StatusOK
      user:{
      "id": 6,
      "email": "czh@qq.com",
      "username": "czh",
      "password": "chenzanhong",
      "avatar_url": "https://www.chen.com",
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImN6aCIsImV4cCI6MTczMjU0MjU3NH0.IoDSg08xnSNK9jlJBr_xdVyYKUIXIbEZm_UsXk0vPmM",
      "create_time": "2024-11-24T21:49:24.78802Z"
      }
    
  - 查询所有用户信息
    - url: http:localhost:8080/user/all
    - 参数格式：无
    - 响应格式
    - ```json
      code: StatusInternalServerError
      message: "获取用户列表失败"
    - ```json
      code: StatusOK
      message: "获取用户列表成功"
      users: [
        {
          //用户信息
        },{
          //……
        }
      ]

  - 更新当前登录用户的信息
    - url: http:localhost:8080/auth/user/update
    - 参数格式： （PUT方法）
    - ```json
      请求头携带一个"Authorization"的token
      {
        //所有参数都是可选的，而且无法更新用户名和id
        "id":,
        "username":,
        "email":,
        "password":,
        "avator_url":,
        "token":,
        "create_time":
      }
    - 响应格式
    - ```json
      code: StatusBadRequest (400)
      message: "请求数据格式错误"
      error:
    - ```json
      code: 401,
      message: "未找到用户信息"
    - ```json
      code: StatusInternalServerError (500)
      message: "更新用户信息失败"
      error: 可能为：
        "用户不存在" "查询用户时发生错误" "用户名不可修改" "邮箱为空" "密码少于6位" "邮箱格式不正确" "更新用户信息失败"
    - ```json
       code: 200,
       message: "用户信息更新成功"
      
    

4. **文生图历史记录**
  - 总的记录，获取所有的图像信息（按id/create_time升序）：
  - Get方法
    - url: localhost:8080/image/all 
    - 参数格式： 无（GET方法）
    - 响应格式：
    - ```json
      code: StatusInternalServerError (500)
      message: "获取图像列表失败"
    - ```json
      code: StatusInternalServerError (500)
      message: "查询失败"
      images: {
        "images": [
            {
                "id": 1,
                "username": "czh0",
                "params": "Prompt: sun, Width: 400, Height: 400, Steps: 30, SamplingMethod: DDIM",
                "picture": "generate/sun-2024-11-21 23:31:24.png",
                "create_time": "2024-11-21T23:31:25.924231Z"
            },
            {
              //……
            },
            //……
        ]
      }


  - 按照时间排序，获取当前登录用户在一段时间内的生成的图像信息
    - url: localhost:8080/auth/user/images/timeRange
    - 参数格式： ?start_time=YYYY-MM-DD&end_time=YYYY-MM-DD 
      （参数值也可以为完整的时间戳2006-01-02T15:04:05.000000Z）
    - ```json
      请求头携带一个"Authorization"的token
    - 响应格式：
    - ```json
      code: StatusBadRequest (400)
      message: "无效的开始时间格式", 
      error:
    - ```json
      code: StatusBadRequest (400)
      message: "无效的结束时间格式", 
      error:

    - ```json
      code: StatusInternalServerError (500)
      message: "查询图像列表失败", 
      error:
    - ```json
      code: StatusOK
      message: "查询图像列表成功", 
      images: {
        "images": [
            {
                "id": 1,
                "username": "czh0",
                "params": "Prompt: sun, Width: 400, Height: 400, Steps: 30, SamplingMethod: DDIM",
                "picture": "generate/sun-2024-11-21 23:31:24.png",
                "create_time": "2024-11-21T23:31:25.924231Z"
            },
            {
              //……
            },
            //……
        ]
      }
  - 获取指定的某张图像
    - url: localhost:8080/auth/image
    - 参数格式：?username= 或?id= 或?url=

    - 响应格式：
    - ```json
      code: StatusNotFound (404)
      message: "未找到相关图片"
    - ```json
      code: StatusInternalServerError (500)
      message: "查询用户的图片失败"
      error: 
    - ```json
      code: StatusBadRequest (400)
      message: "无效的图像id或用户名"
      error:  
    - ```json
      code: StatusOK
      message: "查询图像成功"
      image:  {
                "id": 1,
                "username": "czh0",
                "params": "Prompt: sun, Width: 400, Height: 400, Steps: 30, SamplingMethod: DDIM",
                "picture": "generate/sun-2024-11-21 23:31:24.png",
                "create_time": "2024-11-21T23:31:25.924231Z"
      }


  - 获取当前登录用户生成的所有图像
    - url: localhost:8080/auth/user/images
    - 参数格式：
    - ```json
      请求头携带一个"Authorization"的token
    - 响应格式：
    - ```json
      code: StatusInternalServerError (500)
      message: "查询用户图像失败"
      error:  
    - ```json
      code: StatusOK
      message: "获取用户的图像成功"
      images: {
        "images": [
            {
                "id": 1,
                "username": "czh0",
                "params": "Prompt: sun, Width: 400, Height: 400, Steps: 30, SamplingMethod: DDIM",
                "picture": "generate/sun-2024-11-21 23:31:24.png",
                "create_time": "2024-11-21T23:31:25.924231Z"
            },
            {
              //……
            },
            //……
        ]
      }



   - 点赞图片功能：
   - url："localhost:8080/auth/like"
   - 参数格式：
   - ```json
     请求头携带一个"Authorization"的token
        {
            “url":,//图像url
        }
   - 响应格式：
   - ```json
     "code":  400,
     "error": "Missing image URL"
   - ```json
     "code":  409,
     "error": "用户已经点赞过该图片"
   - ```json
     "code":  500,
     "error": “返回获取赞数错误的error"
   - ```json
     "code":  500,
     "error": "点赞数据库开始出错"
   - ```json
     "current_likes": 当前赞数,
     “message": "Image liked successfully"
   - 按照参数排序，可以列出所需查询参数的查询信息
   - 待定   

5. **图片收藏界面**
  - 查询展示出用户的收藏图片
  - 获取当前用户收藏的图像
    - url: localhost:8080/auth/user/favoritedimages
    - 参数格式：（GET方法）
    - ```json
      请求头携带一个"Authorization"的token
    - 响应格式：同localhost:8080/user/images（只不过message多了一个“收藏”） 

  - 收藏图像
  - 收藏指定图像
    - url：localhost:8080/auth/addFavoritedImage
    - 参数格式：（POST方法）
    - ```json
      请求头携带一个"Authorization"的token
      {
        //两个参数有一个就行
        "url":,//图像url
        "id":,//图像id
      }
    - 响应格式
    - ```json
      code: StatusBadRequest (400)
      message: "无有效的图像id或url"
      error: "id 必须大于 0 或者 url 不得为空"
    - ```json
      code: StatusNotFound (404)
      message: "未找到对应的图像"
      error:  
    - ```json
      code:  401
      message: "未找到用户信息"//没有token时
      error:  
    - ```json
      code:  StatusInternalServerError （500）
      message: "检查收藏状态失败"
      error:  
    - ```json
      code:  StatusConflict
      message: "该图像已经被收藏过"
    - ```json
      code:  StatusInternalServerError （500）
      message: "收藏图像失败"
      error:  
    - ```json
      code:  200
      message: "收藏图像成功"


  - 取消图像收藏
  - 取消指定图像的收藏
    - url：localhost:8080/auth/deleteFavoritedImage
    - 参数格式：?url 或?id（收藏表的图像id，不是图像表的图像id）  （DELETE方法）
    - ```json
      请求头携带一个"Authorization"的token
    - 响应格式
    - ```json
      code: StatusBadRequest (400)
      message: "无有效的图像id或url"
      error: "id 必须大于 0 或者 url 不得为空"
    - ```json
      code: StatusNotFound (404)
      message: "未找到对应的图像"
      error:  
    - ```json
      code:  401
      message: "未找到用户信息"//没有token时
      error:  
    - ```json
      code:  StatusInternalServerError （500）
      message: "检查收藏状态失败"
      error:  
    - ```json
      code:  StatusConflict
      message: "该图像未被收藏过"
    - ```json
      code:  StatusInternalServerError （500）
      message: "取消图像收藏失败"
      error:  
    - ```json
      code:  200
      message: "取消图像收藏成功"

6. **签到增加积分接口**
    - GET方法，用户每次签到增加100积分，并保存记录
    - url：http://localhost:8080/auth/score
    - 参数格式：
    - ```json
      请求头携带一个"Authorization"的token

    - 响应格式：
    - ```json
      code：200（StatusOK）,
      message："用户当前积分"
    - ```json
      code：401（StatusUnauthorized）,
      message："请求头中缺少Token"
    - ```json
      code：401（StatusUnauthorized）,
      message："无效的Token"
    - ```json
      code：401（StatusUnauthorized）,
      success：false,
      message："用户信息查询失败"
    - ```json
      code：401（StatusUnauthorized）,
      success：false,
      message："积分记录创建失败"
    - ```json
      code：401（StatusUnauthorized）,
      success：false,
      message："用户积分更新失败"
    
7. **token校验功能**
    - 校验用户的token
    - url：http://localhost:8080/checkToken
    - 参数格式：
    - ```json
      请求头携带一个"Authorization"的token
    
    - 响应格式：
    - ```json
      code：StatusUnauthorized(401)
      msg："令牌格式不正确"
      data: null
    - ```json
      code：StatusUnauthorized(401)
      msg："令牌过期或未激活"
      data: null
    - ```json
      code：StatusUnauthorized(401)
      msg："令牌无法处理"
      data: null
    - ```json
      code：StatusUnauthorized(401)
      msg："令牌无效"
      data: null
    - ```json
      code：StatusOK(200)
      msg："令牌有效"
      data: tokenStr

8. **数据库设计**
   - 用户登录表：id，email，user，password，token
   - 用户查询表：id，user（外键），params，picture，time
   - 收藏表：id，user（外键），picture
