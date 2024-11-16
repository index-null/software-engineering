<template>
  <el-card class="user-profile">
    <div class="profile-header">
      <div class="avatar-container" @click="handleAvatarClick">
        <img :src="user.avatar" alt="User Avatar" class="avatar">
        <div class="avatar-overlay" v-show="isEditing"></div>
      </div>
      <div class="profile-info" v-if="!isEditing">
        <p class="username">{{ user.username }}</p>
        <p class="email">{{ user.email }}</p>
        <p class="account">{{ user.account }}</p>
      </div>
    </div>
    <el-button type="primary" @click="toggleEdit" class="edit-button" v-show="!isEditing">
      编辑
    </el-button>
    <el-form v-if="isEditing" label-width="80px" class="edit-form">
      <el-form-item label="昵称">
        <el-input v-model="user.username"></el-input>
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="user.email"></el-input>
      </el-form-item>
      <el-form-item label="账号">
        <el-input v-model="user.account" disabled></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="success" @click="saveChanges" class="save-button">保存</el-button>
      </el-form-item>
    </el-form>
    <!-- 隐藏头像上传控件 -->
    <input type="file" ref="fileInput" @change="handleFileChange" style="display: none;">
  </el-card>
</template>

<script>
import OSS from 'ali-oss';

export default {
  data() {
    return {
      user: {
        avatar: 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202408311347060.jpg',
        username: 'Chuhsing',
        email: 'zhuxing.halcyon@gmail.com',
        account: 'index-null'
      },
      isEditing: false,
      client: null
    };
  },
  methods: {
    handleAvatarClick() {
      if (this.isEditing) {
        this.deleteOldAvatar().then(() => {
          this.$refs.fileInput.click();
        });
      }
    },
    async deleteOldAvatar() {
      if (this.user.avatar) {
        const urlParts = this.user.avatar.split('/');
        const fileName = urlParts[urlParts.length - 1];
        try {
          await this.initOSSClient();
          await this.client.delete(`avator/${fileName}`);
          console.log('旧头像删除成功:', fileName);
        } catch (error) {
          console.error('删除旧头像失败:', error);
        }
      }
    },
    async initOSSClient() {
      this.client = new OSS({
        region: process.env.VUE_APP_OSS_REGION,
        accessKeyId: process.env.VUE_APP_OSS_ACCESS_KEY_ID,
        accessKeySecret: process.env.VUE_APP_OSS_ACCESS_KEY_SECRET,
        bucket: process.env.VUE_APP_OSS_BUCKET,
      });
    },
    async handleFileChange(event) {
      const file = event.target.files[0];
      if (!file) return;

      try {
        await this.initOSSClient();
        const result = await this.client.put(`avator/${file.name}`, file);
        this.user.avatar = result.url;
      } catch (error) {
        console.error('上传失败:', error);
      }
    },
    toggleEdit() {
      this.isEditing = !this.isEditing;
    },
    saveChanges() {
      // 这里可以添加保存到后端的逻辑
      console.log('保存用户信息:', this.user);
      this.isEditing = false;
    }
  }
};
</script>

<style scoped>
.user-profile {
  width: 100%; /* 使卡片宽度适应父容器 */
  max-width: 400px; /* 最大宽度为 400px */
  margin: 0 auto;
  text-align: center;
}

.profile-header {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-bottom: 20px;
}

.avatar-container {
  position: relative;
  cursor: pointer;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.avatar-container:hover {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
}

.profile-info {
  margin-top: 10px;
}

.username, .email, .account {
  font-size: 16px;
  margin: 5px 0;
}

.edit-button {
  margin-top: 20px;
}

.save-button {
  margin-left: -10%;
}

.edit-form {
  margin-top: 20px;
  margin-left: -10%; 
}
</style>