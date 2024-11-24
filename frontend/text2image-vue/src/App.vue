<template>
  <div id="app">
    <nav>
      <router-link to="/about" class="nav-link">梗概</router-link>
      <router-link to="/login" class="nav-link">登录</router-link>
      <router-link to="/home" class="nav-link">文生图</router-link>
      <router-link to="/history" class="nav-link">历史记录</router-link>
      <router-link to="/info" class="nav-link">个人资料</router-link>
      <router-link to="/background" class="nav-link">背景</router-link>
    </nav>
    <router-view/>
  </div>
</template>

<script>

export default {
  name: 'App',
};
</script>

<style>

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  
}

nav {
  padding: 30px;
  display: flex;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 为导航栏增加阴影 */
  background-color: #e3f0f8; /* 设置导航栏背景颜色 */
}

.nav-link {
  font-weight: bold;
  color: #2c3e50;
  text-decoration: none;
  margin: 0 100px;
  transition: color 0.3s ease;
}

.nav-link:hover,
.nav-link.router-link-exact-active {
  color: #42b983;
  text-decoration: underline;
}

</style>
