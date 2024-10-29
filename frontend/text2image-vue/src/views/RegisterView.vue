<template>
  <div class="register-container">
    <h1>注册</h1>
    <el-form ref="registerForm" :model="registerForm" :rules="rules" label-width="100px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="registerForm.username" ></el-input>
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="registerForm.email" type="email"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" v-model="registerForm.password"  autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input type="password" v-model="registerForm.confirmPassword"  autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm">注册</el-button>
        <el-button @click="resetForm('registerForm')">重置</el-button>
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
  rules: {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 3, max: 10, message: '用户名长度在 3 到 10 个字符', trigger: 'blur' }
    ],
    email: [
      { required: true, message: '请输入邮箱地址', trigger: 'blur' },
      { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
    ],
    password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { validator: this.validatePass, trigger: 'blur' }
    ],
    confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: this.validatePass2, trigger: 'blur' }
      ]
  },
  methods: {
    validatePass(rule, value, callback) {
      if (value === "") {
        callback(new Error("请输入密码"));
      } else {
        if (this.registerForm.password !== "") {
          this.$refs.registerForm.validateField("confirmPassword");
        }
        callback();
      }
    },
    validatePass2(rule, value, callback) {
      if (value === "") {
        callback(new Error("请再次输入密码"));
      } else if (value !== this.registerForm.password) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    },
    submitForm() {
      this.$refs.registerForm.validate((valid) => {
        if (valid) {
          alert('submit!');
        } else {
          console.log('error submit!!');
          return false;
        }
      });
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