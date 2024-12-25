# 文绘星河

文绘星河 是一个文本到图片的 Web 应用，实现了账号的注册与管理、文生图的前端抽象集成、个人信息管理以及文生图历史查看等功能。
## 项目结构

### 前端

1. **`frontend/text2image-vue/.env`**
   - 前端环境变量配置文件，用于头像的上传以及本地模型生成结果的存储，便于后续与后端对接。

2. **`frontend/text2image-vue/src/views`**
   - **`AboutView.vue`**: 项目介绍页
   - **`ExploreView.vue`**: 探索界面（含图片广场）
   - **`FavoritesView.vue`**: 收藏夹界面
   - **`GenerateView.vue`**: 文生图界面
   - **`HistoryView.vue`**: 文生图历史记录界面
   - **`InnerView.vue`**: 主导航界面（功能页为其子路由）
   - **`LogRegView.vue`**: 登录注册界面
   - **`SettingView.vue`**: 个人信息设置界面
   - **`UsageGuide.vue`**: 使用说明界面

   > 项目的所有页面视图定义，Vue-Router的路由定义基于此处的页面。

3. **`frontend/text2image-vue/src/components`**
   - 考虑到页面需要复用的组件不多，且定制性较强，所以没有特地进行组件拆分，此处无复用组件。

4. **`frontend/text2image-vue/src/assets`**
   - 存储静态资源，如项目介绍页的图片、按钮图标、首页的常驻图片等。

5. **`frontend/text2image-vue/src/router/index.js`**
   - 路由定义，路由跳转定义，路由守卫，页面懒加载等。

6. **`frontend/text2image-vue/src/api/index.js`**
   - API定义，封装了axios请求，自动在请求带上token，主要便于开发和前后端交互。
   
## 前端启动

### 配置

1. **切换到前端项目目录**：
   
   ```bash
   cd frontend/text2image-vue
   ```
   
2. **创建 `.env` 文件**：
   ```bash
   touch .env
   vim .env
   ```

3. **编辑 `.env` 文件**：
   在 `.env` 文件中，填入阿里云OSS的密钥。以下是一个示例配置：
   ```plaintext
   VUE_APP_OSS_REGION=oss-cn-shenzhen
   VUE_APP_OSS_ACCESS_KEY_ID=your_access_key_id
   VUE_APP_OSS_ACCESS_KEY_SECRET=your_access_key_secret
   VUE_APP_OSS_BUCKET=your_bucket_name
   ```

### 启动前端

1. **安装依赖**：
   ```bash
   npm install
   ```

2. **启动开发服务器**：
   ```bash
   npm run serve
   ```

## 后端启动

### 配置

1. **切换到后端项目目录**：
   ```bash
   cd backend/text-to-picture
   ```

2. **拉取依赖**：
   ```bash
   go mod tidy
   ```
   如果有拉取不完全的报错，可以使用 `go get + 依赖` 进行手动拉取。

3. **修改数据库配置**：
   - 进入 `backend/text-to-picture/config/DBconfig` 目录，修改 `DBconfig.yaml.example` 文件为自己的数据库配置,并重命名为 `DBconfig.yaml`。

4. **配置阿里云OSS**：
    在 `text-to-picture/config/oss` 目录下创建 `oss.env` 文件，写入阿里云OSS存储的配置信息。
   ```bash
   cd config/oss
   touch oss.env
   vim oss.env
   ```
   
   以下是一个示例配置：
     ```plaintext
     OSS_REGION=oss-cn-shenzhen
     OSS_ACCESS_KEY_ID=your_access_key_id
     OSS_ACCESS_KEY_SECRET=your_access_key_secret
     OSS_BUCKET=your_bucket_name
     ```


### 启动后端

1. **运行后端服务**：
   ```bash
   go run main.go
   ```

## 其他注意事项

- **暂无**