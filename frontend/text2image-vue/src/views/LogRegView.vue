<template>
  <div class="container">
   <div class="artistic-nav"></div>
  <div class="log-reg-container">
    <div class="left-side-container">
      <img src="https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202412081939173.png" alt="#">
    </div>
    <div class="right-side-container">
      <LogoAndAppName />
      <el-tabs type="border-card">
        <el-tab-pane label="注册">
          <!-- <el-card class="box-card"> -->
            <h2>注册</h2>
            <el-form
              :model="registerForm"
              status-icon
              :rules="registerRules"
              ref="registerForm"
              label-position="left"
              label-width="80px"
              class="register-form"
            >
              <el-form-item label="用户名" prop="uname">
                <el-input v-model="registerForm.uname"></el-input>
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="registerForm.email"></el-input>
              </el-form-item>
              <el-form-item label="密码" prop="pass">
                <el-input
                  type="password"
                  v-model="registerForm.pass"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
              <el-form-item label="确认密码" prop="password">
                <el-input
                  type="password"
                  v-model="registerForm.password"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
            </el-form>
            <div class="btn-group">
              <el-button type="primary" @click="submitRegisterForm('registerForm')">提交</el-button>
              <el-button @click="resetRegisterForm('registerForm')">重置</el-button>
            </div>
          <!-- </el-card> -->
        </el-tab-pane>
        <el-tab-pane label="登录">
          <!-- <el-card class="box-card"> -->
            <h2>登录</h2>
            <el-form
              :model="loginForm"
              status-icon
              :rules="loginRules"
              ref="loginForm"
              label-position="left"
              label-width="70px"
              class="login-form"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="loginForm.username"></el-input>
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input
                  type="password"
                  v-model="loginForm.password"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
            </el-form>
            <div class="btn-group">
              <el-button type="primary" @click="submitLoginForm('loginForm')">登录</el-button>
              <el-button @click="resetLoginForm('loginForm')">重置</el-button>
            </div>
          <!-- </el-card> -->
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</div>
</template>

<script>
import axios from 'axios';

export default {
  components: {
    LogoAndAppName: () => import('../components/LogoAndAppName.vue')
  },
  data() {
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else if (value.length < 6 || value.length > 16) {
        callback(new Error('密码长度必须在6到16个字符之间'));
      } else {
        if (this.registerForm.password !== '') {
          this.$refs.registerForm.validateField('password');
        }
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请再次输入密码"));
      } else if (value !== this.registerForm.pass) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    };
    return {
      loginForm: {
        username: '',
        password: '',
      },
      loginRules: {
        username: [
          { required: true, message: "用户名不能为空！", trigger: "blur" },
        ],
        password: [
          { required: true, message: "密码不能为空！", trigger: "blur" },
          { min: 6, max: 16, message: '密码长度必须在6到16个字符之间', trigger: 'change' }
        ],
      },
      registerForm: {
        uname: "",
        email: "",
        pass: "",
        password: "",
      },
      registerRules: {
        uname: [{ required: true, message: "请输入用户名", trigger: "blur" }],
        email: [
          { required: true, message: "请输入邮箱地址", trigger: "blur" },
          { type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] },
        ],
        pass: [{ required: true, validator: validatePass, trigger: "blur" }],
        password: [
          { required: true, validator: validatePass2, trigger: "blur" },
        ],
      },
    };
  },
  methods: {
  async hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return hashHex;
  },
  async submitLoginForm(formName) {
    this.$refs[formName].validate(async (valid) => {
      if (valid) {
        let hashedPassword = await this.hashPassword(this.loginForm.password);
        let formData = {
          username: this.loginForm.username,
          password: hashedPassword
        };
        try {
          const response = await axios.post('http://localhost:8080/login', formData);
          if (response.data.code === 200) {
            localStorage.setItem('token', response.data.token);
            this.$message.success('登录成功');
            this.$router.push('/main');
          } 
        } catch (error) {
          if (error.response) {
            // 服务器返回了错误响应
            switch (error.response.status) {
              case 401:
                this.$message.error('用户名或密码错误');
                break;
              default:
                this.$message.error('登录过程中发生错误');
            }
          } else {
            // 其他错误，如网络问题
            this.$message.error('登录过程中发生错误');
          }
        }
      } else {
        console.log("error submit!!");
        return false;
      }
    });
  },
  resetLoginForm(formName) {
    this.$refs[formName].resetFields();
  },
  async submitRegisterForm(formName) {
    this.$refs[formName].validate(async (valid) => {
      if (valid) {
        try {
          const hashedPassword = await this.hashPassword(this.registerForm.password);
          const formattedFormData = {
            email: this.registerForm.email,
            username: this.registerForm.uname,
            password: hashedPassword
          };
          const response = await axios.post('http://localhost:8080/register', formattedFormData);
          if (response.data.code === 200) {
            this.$message.success('注册成功');
          } 
        } catch (error) {
          if (error.response) {
            // 服务器返回了错误响应
            switch (error.response.status) {
              case 401:
                this.$message.error('用户名或密码错误');
                break;
              default:
                this.$message.error('注册过程中发生错误');
            }
          } else {
            // 其他错误，如网络问题
            this.$message.error('注册过程中发生错误');
          }
        }
      }
    });
  },
    resetRegisterForm(formName) {
      this.$refs[formName].resetFields();
    },
  },
};
</script>

<style scoped>
.log-reg-container {
  display: flex;
  height: 100vh;
  background-color: #f8fffe;
}

.left-side-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.right-side-container {
  flex: 1;
  display: flex;
  justify-content: center;
  gap: 00px;
  align-items: center;
  flex-direction: column;
}

.box-card {
  width: 400px;
  padding: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background-color: #ffffff;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333333;
}

.login-form, .register-form {
  margin: 0 auto;
  width: 100%;
}

.btn-group {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.el-button {
  flex: 1;
  margin: 0 5px;
}

.el-button + .el-button {
  margin-left: 10px;
}

.el-input {
  width: 100%;
  margin-bottom: 10px;
  border-radius: 4px;
  transition: border-color 0.3s ease;
}
.artistic-nav {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 50px;
  background: linear-gradient(270deg, #ff7e5f, #feb47b, #86a8e7, #7f7fd5);
  background-size: 800% 800%;
  animation: gradientAnimation 15s ease infinite;
  z-index: 1000;
}

@keyframes gradientAnimation {
  0%{background-position:0% 50%}
  50%{background-position:100% 50%}
  100%{background-position:0% 50%}
}

</style>