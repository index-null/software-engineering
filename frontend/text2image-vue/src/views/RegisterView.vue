<template>
    <div class="register-container">
      <el-form ref="registerForm" :model="registerForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input v-model="registerForm.username"></el-input>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="registerForm.email" type="email"></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input type="password" v-model="registerForm.password"></el-input>
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input type="password" v-model="registerForm.confirmPassword"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm">注册</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        registerForm: {
          username: '',
          email: '',
          password: '',
          confirmPassword: ''
        }
      };
    },
    methods: {
      submitForm() {
        if (this.registerForm.password !== this.registerForm.confirmPassword) {
          this.$message.error('两次输入的密码不一致');
          return;
        }
        // 发送到服务器的注册请求
        console.log('注册信息:', this.registerForm);
        // 这里可以添加 axios 请求发送数据到服务器
      },
      resetForm() {
        this.$refs.registerForm.resetFields();
      }
    }
  };
  </script>
  
  <style scoped>
  .register-container {
    max-width: 400px;
    margin: 50px auto;
    padding: 20px;
    border: 1px solid #ebeef5;
    border-radius: 4px;
  }
  
  .el-form-item {
    margin-bottom: 20px;
  }
  </style>