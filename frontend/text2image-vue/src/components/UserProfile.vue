<template>
  <div class="user-profile">
    <h1>个人信息</h1>
    <div class="profile-info">
      <img :src="user.avatar" alt="用户头像" class="avatar">
      <div class="info">
        <p><strong>姓名:</strong> {{ user.name }}</p>
        <p><strong>邮箱:</strong> {{ user.email }}</p>
        <p><strong>电话:</strong> {{ user.phone }}</p>
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
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 20px;
}

.info p {
  margin: 10px 0;
}
</style>
<script>
export default {
  name: 'UserProfile',
  data() {
    return {
      user: {
        avatar: 'https://via.placeholder.com/150',
        name: '张三',
        email: 'zhangsan@example.com',
        phone: '123-456-7890'
      }
    };
  }
};
</script>
