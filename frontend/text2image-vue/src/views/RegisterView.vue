<template>
  <div>
    <el-card class="box-card">
      <h2>注册</h2>
      <el-form
        :model="registerForm"
        status-icon
        :rules="rules"
        ref="registerForm"
        label-position="left"
        label-width="80px"
        class="demo-ruleForm"
      >
        <el-form-item label="用户名" prop="uname">
          <el-input v-model="registerForm.uname"></el-input>
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
      <div class="btnGroup">
        <el-button type="primary" @click="submitForm('registerForm')"
          >提交</el-button
        >
        <el-button @click="resetForm('registerForm')">重置</el-button>
        <el-button @click="goBack">返回</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请输入密码"));
      } else {
        if (this.registerForm.checkPass !== "") {
          this.$refs.registerForm.validateField("checkPass");
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
      registerForm: {
        uname: "",
        pass: "",
        password: "",
      },
      rules: {
        uname: [
          { required: true, message: "用户名不能为空！", trigger: "blur" },
        ],
        pass: [{ required: true, validator: validatePass, trigger: "blur" }],
        password: [
          { required: true, validator: validatePass2, trigger: "blur" },
        ],
      },
    };
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
        // 后端设置路由处理注册信息
        axios.post('/api/register', this.registerForm)
          .then(response => {
            if (response.data.success) {
              // 注册成功，保存用户名到本地浏览器中
              localStorage.setItem('registeredUsername', this.registerForm.uname);
              // 跳转到登录界面
              this.$router.push('/login');
            } else {
              // 注册失败，显示错误消息
              this.$message.error(response.data.message);
            }
          })
          .catch(error => {
             console.error('注册失败:', error);
            if (error.response) {
            // 服务器端错误，可以获取错误状态和消息
            this.$message.error(`服务器错误: ${error.response.status} - ${error.response.data.message}`);
            } else if (error.request) {
            // 请求已发出，但未收到响应
            this.$message.error('请求未响应，请检查网络连接');
            } else {
            // 发送请求时出错
            this.$message.error('请求发送失败，请稍后再试');
            }
          });
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
    goBack() {
      //返回登陆界面
      if(this.$router.path!=='/login'){
        this.$router.push('/login');
      }else{
        this.$router.go(-1);
      }
    }
  }
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