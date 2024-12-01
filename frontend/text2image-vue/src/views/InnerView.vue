<template>
  <div class="inner-view">
    <div class="left-side-container">
      <div class="logo-container">
        <img src='../assets/文生图.svg' alt="Logo" class="logo"/>
        <span class="title">{{ appName }}</span>
      </div>
      <nav class="nav-container">
        <router-link to="/main/explore" class="nav-link">探索发现</router-link>
        <router-link to="/main/generate" class="nav-link">文字作画</router-link>
        <router-link to="/main/favorites" class="nav-link">我的收藏</router-link>
        <router-link to="/main/history" class="nav-link">作画历史</router-link>
      </nav>
      <div class="avatar-container"></div>
      <div class="nickname-container"></div>
      <div class="setting-container"></div>
    </div>
    <div class="view-container">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    appName() {
      return this.$store.state.appName;
    },
  },
}
</script>

<style scoped>
.inner-view {
  display: flex;
  height: 100vh; /* 使整个容器高度占满视口 */
  background-color: #f0f2f5; /* 背景颜色 */
}

.left-side-container {
  flex: 1; /* 占1份 */
  background-color: #333333; /* 左侧背景颜色 */
  color: #ffffff; /* 文字颜色 */
  padding: 20px; /* 内边距 */
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.logo-container, .avatar-container, .nickname-container, .setting-container {
  height: 50px; /* 高度可以根据需要调整 */
  background-color: #444444; /* 背景颜色 */
  margin-bottom: 10px; /* 下边距 */
}

.nav-container {
  display: flex;
  flex-direction: column;
  gap: 10px; /* 项目之间的间距 */
}

.nav-link {
  color: #ffffff; /* 链接颜色 */
  text-decoration: none; /* 去掉下划线 */
  font-size: 16px; /* 字体大小 */
  padding: 10px; /* 内边距 */
  border-radius: 4px; /* 圆角 */
  transition: background-color 0.3s ease; /* 平滑过渡效果 */
}

.nav-link:hover {
  background-color: #555555; /* 鼠标悬停时的背景颜色 */
}

.view-container {
  flex: 10; /* 占10份 */
  background-color: #ffffff; /* 右侧背景颜色 */
  padding: 20px; /* 内边距 */
}
</style>