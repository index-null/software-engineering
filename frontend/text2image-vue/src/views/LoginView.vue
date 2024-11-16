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
        <el-form-item label="用户名" prop="uname">
          <el-input v-model="ruleForm.uname" @blur="checkUsername"></el-input>
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
        uname: localStorage.getItem('registeredUsername') || '', // 从localStorage获取注册的账号
        password: '',
      },
      rules: {
        uname: [
          { required: true, message: "用户名不能为空！", trigger: "blur" },
        ],
        password: [
          { required: true, message: "密码不能为空！", trigger: "blur" },
        ],
      },
    };
  },
  methods: {
    checkUsername() {
      if (this.ruleForm.uname.trim().length > 0) {
        const username = this.ruleForm.uname;
      axios.post('/api/check-username', { username })
        .then(response => {
          if (response.data.registered) {
            this.$message.success('该用户名已注册，可以登录');
          } else {
            this.$message.error('该用户名未注册，请先注册');
          }
        })
        .catch(error => {
          console.error('检查用户名时发生错误:', error);
          this.$message.error('在检查用户名过程中发生错误', { duration: 2000 });
        });
      }
    },
    submitForm(formName) {
      /* eslint-enable no-unreachable */
      //测试代码
      this.$router.push('/home');
      /* eslint-enable no-unreachable */
      this.$refs[formName].validate((valid) => {
        if (valid) {
          axios.post('/api/login', this.ruleForm)
            .then(response => {
              if (response.data.success) {
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