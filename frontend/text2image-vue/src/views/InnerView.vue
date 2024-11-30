<template>
  <div class="inner-view">
    <nav class="nav-container">
    <router-link to="/main/explore" class="nav-link">探索发现</router-link>
      <router-link to="/main/generate" class="nav-link">文字作画</router-link>
      <router-link to="/main/favorites" class="nav-link">我的收藏</router-link>
      <router-link to="/main/history" class="nav-link">作画历史</router-link>
    </nav>
    <div class="view-container">
        <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {


}
</script>

<style scoped>

</style>