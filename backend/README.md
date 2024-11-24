包说明
1. main.go程序入口用于初始化各个模块并启动服务。
2. api/
API 层

处理所有的 HTTP 请求和响应。
定义与前端交互的 RESTful API 路由，处理用户请求并将请求转发到服务层。
定义了三个接口分别是：auth——认证，history——查询历史记录，generate——生成图
3. models/
数据模型层

定义数据库中使用的所有实体模型和结构体。
该层负责将数据库中的数据映射到 Go 语言的结构体中（ORM）。
示例文件：
init:初始化
image：建image表
user：用户表
repository/
数据持久化层
负责与数据库的交互，进行数据的存储、更新、删除、查询等操作。
该层会调用模型，并根据业务需求编写具体的数据库查询逻辑。

5. services/
业务逻辑层
该层包含系统的核心业务逻辑。
接收来自 api/ 层的请求，调用 repository/ 层来操作数据库，并返回结果。
示例文件：
对应api层里的三个接口，在这里实现其逻辑
6. config/
配置层

存储应用程序的配置文件（如数据库连接信息、API 密钥等）。
负责读取和管理项目的全局配置参数。
（方便修改程序信息，昨晚建漏了这个包，没事全部写完再加也可以哈哈）
7. utils/
工具层

提供常用的工具类函数和公共方法，这些工具通常会被多个模块调用。
该层可以包含日志工具、错误处理工具、加密工具等。
示例文件：
我想到的工具就jwt生成密码的那个hash，还有response
8. middleware/
中间件层
处理请求生命周期中的中间件逻辑，如认证、日志、跨域请求处理等。
中间件通常用于拦截和处理 HTTP 请求的预处理。
示例文件：
jwt：登录要用到的登录验证中间件他会返回一个token用于身份验证登录
9. docs/
文档

项目的文档存放处，包含API文档、架构图、使用说明等。
该目录帮助开发者和使用者理解项目的功能和使用方法。（我们好像没有，先建着吧）
### 开发内容

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
          Code: 400（StatusBadRequest）,
          Msg:  "请求数据格式错误"
        - ```json
          Code: 500（StatusInternalServerError,
          message: "用户创建失败",
          error:   err.Error()
        - ```json  
          Code: 200（StatusOK）,
		  Msg:  "注册成功",

   - 登录访问url：http://localhost:8080/login
   - 登录数据格式：
   - ```json
     {
        "username": "root1",
        "password": "sssssss"
     }
   - 响应格式：
   - ```json
     Code：400（StatusBadRequest）,
     Msg："请求数据格式错误"
   - ```json
     Code：401（Unauthorized）,
     Msg: "用户不存在"
   - ```json
     Code：500（StatusInternalServerError）,
     Msg: "数据库查询错误"
   - ```json
     Code：401（Unauthorized）,
     Msg: "密码错误"
   - ```json
     Code：200（StatusOK）,
     Msg："登录成功",
2. **文生图接口**
   - 部署本地的文生图模型，编写接口进行传参和调用
   - 接收前端的参数，调用本地部署的大模型，生成对应的图片，返回给前端，并将记录存入数据库
   - 假文生图url：http://localhost:8080/auth/generate
   - 参数格式：
   - ```json
     {
       "height": 200,
       "width": 220,
       "prompt": "string",
       "sampling_method": "DDIM",
       "seed": "string",
       "steps": 100
     }
    - 响应格式：
   - ```json
     Code：200（StatusOK）,
     Msg："成功响应"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："缺乏提示词"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："宽度不在范围内"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："高度不在范围内"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："步数不在范围内"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："采样方法不在范围内"
   - ```json
     Code：400（StatusBadRequest）,
     Msg："缺乏种子"
   - ```json
     Code：500（StatusInternalServerError）,
     Msg："图片生成失败"
   3. **个人信息界面**
      - 结合数据库的用户信息，使用查询函数查询出需要的信息返回给前端
      - 头像上传功能，获取功能set，get
        - 修改头像url：http://localhost:8080/auth/setavatar
        - 参数格式：
        - ```json
          携带一个"Authorization"的token
          "token": "string"(jwt生成的token)
          "url": "string"(更换头像的url)
        - 响应格式：
        - ```json
          Code: 401（Unauthorized）,
          Msg:  "请求头中缺少Token"
        - ```json
          Code: 401（Unauthorized）,
          Msg:  "无效的Token"
        - ```json
          Code: 401（Unauthorized）,
          Msg:  "请求头中缺少Token"
        - ```json
          Code: 500(Error),
          Msg:  "更新头像失败"
        - ```json
          Code: 200(Success),
          Msg:  "头像更新成功",
          Data: newURL
        - 获取头像url：http://localhost:8080/auth/getavatar
        - 参数格式：
        - ```json
          携带一个"Authorization"的token
          "token": "string"(jwt生成的token)
        - 响应格式：
        - ```json
          Code: 401（Unauthorized）,
          Msg:  "Token已过期"
        - ```json
          Code: 401（Unauthorized）,
          Msg:  "无效的Token"
        - ```json
          Code: Error,
          Msg:  "查询头像失败"
        - ```json  
          Code: Success,
		  Msg:  "获取头像成功",
		  Data: usera.Avatar_url
4. **文生图历史记录**
   - 总的记录
   - 按照时间排序，可以列出一段时间内的查询信息
   - 按照参数排序，可以列出所需查询参数的查询信息
   - 按照结果排序，可以列出所需结果图片的查询信息
   - 按照用户排序，可以列出所需用户查询的查询信息
   - 待定   

5. **图片收藏界面**
   - 查询展示出用户的收藏图片

6. **数据库设计**
   - 用户登录表：id，email，user，password，token
   - 用户查询表：id，user（外键），params，result，time
   - 收藏表：id，user（外键），result
