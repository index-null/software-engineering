<template>
    <div>
      <el-card class="box-card">
        <h2>登录</h2>
        <el-form
          :model="ruleForm"
          status-icon
          :rules="rules"
          ref="ruleForm"
          label-position="left"
          label-width="70px"
          class="login-from"
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
        <div class="btnGroup">
          <el-button type="primary" @click="submitForm('ruleForm')"
            >登录</el-button
          >
          <el-button @click="resetForm('ruleForm')">重置</el-button>
          <router-link to="/register">
            <el-button style="margin-left:10px">注册</el-button>
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
    /*created(){
      //检查是否有已注册的用户名存储在 localStorage 中
      const registeredUsername = localStorage.getItem('registeredUsername');
      if(registeredUsername){
        this.ruleForm.uname = registeredUsername; // 填充用户名
        this.$message.success(`欢迎回来，${registeredUsername}`);
      }else{
        this.$message.error('请先注册', {
        duration: 2000 // 显示时间
        });
      }
    },*/
    methods: {
      // 检查用户名是否已注册
    checkUsername() {
      const username = this.ruleForm.uname;
      //后端
      axios.post('/api/check-username', { username })
        .then(response => {
         if (response.data.registered) {
            // 账号已注册，可以进行登录操作
            this.$message.success('该用户名已注册，可以登录');
          } else {
            // 账号未注册，提示用户
            this.$message.error('该用户名未注册，请先注册');
          }
        })
      .catch(error => {
        console.error('检查用户名时发生错误:', error);
        this.$message.error('在检查用户名过程中发生错误',{
        duration: 2000 // 显示时间
        });
      });
  },
  submitForm(formName) {
    this.$refs[formName].validate((valid) => {
      if (valid) {
        // 发送登录请求
        axios.post('/api/login', this.ruleForm)
          .then(response => {
            if (response.data.success) {
              // 登录成功，跳转到首页或其他页面
              this.$router.push('/');
            } else {
              // 登录失败，显示错误消息
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
  .box-card {
    margin: auto auto;
    width: 400px;
  }
  .login-from {
    margin: auto auto;
  }
  </style>