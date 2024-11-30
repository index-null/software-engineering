<template>
  <section>
    <div class="box">
      <div class="square" v-for="n in 6" :key="n"></div>
    </div>
    <div class="container">
      <div class="form">
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
      </div>
    </div>
  </section>
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
          let formData={username: this.ruleForm.username, password: hashedPassword}
          axios.post('http://localhost:8080/login', formData)
            .then(response => {
              if (response.data.message === '登录成功') {
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
@import url('https://fonts.googleapis.com/css2?family=El+Messiri:wght@700&display=swap');

* {
  margin: 0;
  padding: 0;
  font-family: 'El Messiri', sans-serif;
}

body {
  background: #031323;
  overflow: hidden;
}

section {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5; /* 纯色背景 */
  background-image: url(https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202411251503682.png);
  background-size: cover;
  background-position: center;
}

@keyframes gradient {
  0%,100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.box {
  position: relative;
}
  .square {
    position: absolute;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(5px);
    box-shadow: 0 25px 45px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 15px;
    animation: square 10s linear infinite;
    animation-delay: calc(-1s * var(--i));
  }
  @keyframes square {
    0%,100% {
      transform: translateY(-20px);
    }
    50% {
      transform: translateY(20px);
    }
  }
  .square:nth-child(1) {
    width: 100px;
    height: 100px;
    top: -15px;
    right: -45px;
  }
  .square:nth-child(2) {
    width: 150px;
    height: 150px;
    top: 105px;
    left: -125px;
    z-index: 2;
  }
  .square:nth-child(3) {
    width: 60px;
    height: 60px;
    bottom: 85px;
    right: -45px;
    z-index: 2;
  }
  .square:nth-child(4) {
    width: 50px;
    height: 50px;
    bottom: 35px;
    left: -95px;
  }
  .square:nth-child(5) {
    width: 50px;
    height: 50px;
    top: -15px;
    left: -25px;
  }
  .square:nth-child(6) {
    width: 85px;
    height: 85px;
    top: 165px;
    right: -155px;
    z-index: 2;
  }
}

.container {
  position: relative;
  padding: 50px;
  width: 260px;
  min-height: 380px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(5px);
  border-radius: 10px;
  box-shadow: 0 25px 45px rgba(0, 0, 0, 0.2);
}

.container::after {
  content: '';
  position: absolute;
  top: 5px;
  right: 5px;
  bottom: 5px;
  left: 5px;
  border-radius: 5px;
  pointer-events: none;
  background: linear-gradient( to bottom, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.1) 2%
  );
}

.form {
  position: relative;
  width: 100%;
  height: 100%;

  h2 {
    color: #fff;
    letter-spacing: 2px;
    margin-bottom: 30px;
  }

  .inputBx {
    position: relative;
    width: 100%;
    margin-bottom: 20px;

    input {
      width: 80%;
      outline: none;
      border: none;
      border: 1px solid rgba(255, 255, 255, 0.2);
      background: rgba(255, 255, 255, 0.2);
      padding: 8px 10px;
      padding-left: 40px;
      border-radius: 15px;
      color: #fff;
      font-size: 16px;
      box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
    }

    .password-control {
      position: absolute;
      top: 11px;
      right: 10px;
      display: inline-block;
      width: 20px;
      height: 20px;
      background: url(https://snipp.ru/demo/495/view.svg) 0 0 no-repeat;
      transition: 0.5s;
    }

    .view {
      background: url(https://snipp.ru/demo/495/no-view.svg) 0 0 no-repeat;
      transition: 0.5s;
    }

    .fas {
      position: absolute;
      top: 13px;
      left: 13px;
    }

    input[type="submit"] {
      background: #fff;
      color: #111;
      max-width: 100px;
      padding: 8px 10px;
      box-shadow: none;
      letter-spacing: 1px;
      cursor: pointer;
      transition: 1.5s;
    }

    input[type="submit"]:hover {
      background: linear-gradient(115deg, 
        rgba(0,0,0,0.10), 
        rgba(255,255,255,0.25));
      color: #fff;
      transition: .5s;
    }

    input::placeholder {
      color: #fff;
    }

    span {
      position: absolute;
      left: 30px;
      padding: 10px;
      display: inline-block;
      color: #fff;
      transition: .5s;
      pointer-events: none;
    }

    input:focus ~ span,
    input:valid ~ span {
      transform: translateX(-30px) translateY(-25px);
      font-size: 12px;
    }
  }

  p {
    color: #fff;
    font-size: 15px;
    margin-top: 5px;

    a {
      color: #fff;
    }

    a:hover {
      background-color: #000;
      background-image: linear-gradient(to right, #434343 0%, black 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }
}

.remember {
  position: relative;
  display: inline-block;
  color: #fff;
  margin-bottom: 10px;
  cursor: pointer;
}
</style>
