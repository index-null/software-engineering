<template>
  <div class="login-container">
    <el-card class="box-card">
      <h2>登录</h2>
      <el-form
        :model="ruleForm"
        status-icon
        :rules="rules"
        ref="ruleForm"
        label-position="left"
        label-width="70px"
        class="login-form"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="ruleForm.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            type="password"
            v-model="ruleForm.password"
            autocomplete="off"
          ></el-input>
        </el-form-item>
      </el-form>
      <div class="btn-group">
        <el-button type="primary" @click="submitForm('ruleForm')">登录</el-button>
        <el-button @click="resetForm('ruleForm')">重置</el-button>
        <router-link to="/register">
          <el-button style="margin-left: 10px">注册</el-button>
        </router-link>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      ruleForm: {
        username: '',
        password: '',
      },
      rules: {
        username: [
          { required: true, message: "用户名不能为空！", trigger: "blur" },
        ],
        password: [
        { required: true, message: "密码不能为空！", trigger: "blur" },
        { min: 6, max: 16, message: '密码长度必须在6到16个字符之间', trigger: 'change' }
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
    async submitForm(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          let hashedPassword = await this.hashPassword(this.ruleForm.password);
          let formData={
            username: this.ruleForm.username,
            password: hashedPassword
          }
          axios.post('http://localhost:8080/login', formData)
            .then(response => {
              if (response.data.message === '登录成功') {
                localStorage.setItem('token', response.data.token);
                this.$message.success('登录成功');
                this.$router.push('/home');
                
              } else {
                this.$message.error(response.data.message);
              }
            })
            .catch(error => {
              console.error('登录失败:', error);
              this.$message.error('登录过程中发生错误');
            });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
  },
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5; /* 纯色背景 */
  background-image: url(https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202411251503682.png);
  background-size: cover;
  background-position: center;
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

.login-form {
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
</style>