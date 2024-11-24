<template>
  <div id="app">
    <nav>
      <router-link to="/about" class="nav-link" :class="{ active: isActive('/about') }">About</router-link>
      <router-link to="/login" class="nav-link" :class="{ active: isActive('/login') }">立即登录</router-link>
      <router-link to="/home" class="nav-link" :class="{ active: isActive('/home') }">文字作画</router-link>
      <router-link to="/history" class="nav-link" :class="{ active: isActive('/history') }">作画历史</router-link>
      <router-link to="/info" class="nav-link" :class="{ active: isActive('/info') }">个人中心</router-link>
      
    </nav>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    isActive(path) {
      return this.$route.path === path;
    }
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  display: flex;
  height: 100vh;
  background-color: #f0f8ff;
}

nav {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  width: 120px; /* 确保这个宽度足够容纳最宽的导航项 */
  height: 100%;
  background-color: #ffffff;
  border-radius: 30px;
  padding: 20px 0;
  box-shadow: none;
}

.nav-link {
  color: #2c3e50;
  text-decoration: none;
  margin: 15px 0;
  transition: color 0.3s ease, background-color 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 40px;
  border-radius: 20px;
  background-color: transparent;
  cursor: pointer;
  padding: 0 10px; /* 减少内边距以避免超出边界 */
  text-align: center;
}

.nav-link:hover, .nav-link.router-link-exact-active {
  color: #42b983;
  background-color: #e0f7fa;
}

.nav-link.active {
  background-color: #ffffff; /* 激活时的背景色 */
  box-shadow: inset 2px 2px 2px rgba(66, 185, 131, 0.5); /* 调整圆框的大小 */
}
</style>