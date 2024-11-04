<template>
    <div class="user-profile">
      <h1>个人信息</h1>
      <div class="profile-info">
       <el-row class="demo-avatar demo-basic">
      <el-col :span="24">
        <div class="sub-title"></div>
        <div class="demo-basic--circle">
          <div class="block">
            <el-avatar :size="50" :src="user.circleUrl"></el-avatar>
          </div>
        </div>
      </el-col>
    </el-row>
        <div class="info">
          <p><strong>姓名:</strong> {{ user.name }}</p>
          <p><strong>邮箱:</strong> {{ user.email }}</p>
          <p><strong>电话:</strong> {{ user.phone }}</p>
          <p><strong>居住地址:</strong> {{ user.address }}</p>
          <p><strong>学校:</strong> {{ user.school }}</p>
          <p><strong>个人签名:</strong></p>
          <div style="margin: 20px 0;"></div>
  <el-input
    type="textarea"
    :autosize="{ minRows: 2, maxRows: 4}"
    placeholder="请输入内容"
    v-model="user.textarea2">
  </el-input>
        </div>
      </div>
    </div>
  </template>
  
  <style scoped>
  .user-profile {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
  }
  
  .avatar {
    position: relative;
    display: inline-block;
    width: 100px;
    height: 100px;
    border-radius: 50%;
    overflow: hidden;
    margin-right: 20px;
  }
  
  .change-avatar-btn {
    position: absolute;
    bottom: 0;
    right: 0;
    background-color: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    border-radius: 0 0 5px 0;
    cursor: pointer;
    padding: 2px 8px;
    font-size: 12px;
  }
  
  .info p {
    margin: 10px 0;
  }
  .demo-avatar .sub-title {
    font-size: 16px;
    color: #5e6d82;
    font-weight: bold;
    margin-bottom: 20px;
  }
  
  .demo-basic--circle .block {
    display: inline-flex;
    align-items: center;
    margin-right: 20px;
  }
  </style>
  
  <script>
  export default {
    data() {
      return {
        user: {
          avatar: require('@/assets/logo.png'), // 确保路径正确
          name: '张三',
          email: 'zhangsan@example.com',
          phone: '123-456-7890',
          address: '北京市朝阳区',
          school: '清华大学',
          signature: '',
          circleUrl: "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
          textarea1: '',
        textarea2: ''
        } 
        
      };
    },
     methods: {
      changeAvatar() {
        this.$refs.fileInput.click();
      },
      onFileChange(e) {
        const file = e.target.files[0];
        if (!file) return;
  
        const reader = new FileReader();
        reader.onload = (e) => {
          this.user.avatar = e.target.result;
        };
        reader.readAsDataURL(file);
      }
    }
  };
  </script>
  
  <!-- 隐藏的文件输入元素 -->
  <input type="file" ref="fileInput" @change="onFileChange" style="display: none;" />