<template>
  <el-row class="background-overlay">
    <el-col :span="24">
      <div class="content">
        <h1>这是背景覆盖效果</h1>
        <el-button type="primary" @click="handleClick">按钮</el-button>
      </div>
    </el-col>
  </el-row>
</template>

<script>
export default {
  name: 'BackgroundOverlay',
  methods: {
    handleClick() {
      // 这里可以添加按钮点击时的逻辑
      console.log('Button clicked');    
    }
  }
}
</script>

<style scoped>
.background-overlay {
  position: relative;
  width: 100%;
  height: 100vh;  /* 设置背景覆盖整个视口 */
  background-image: url('~@/assets/3.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.content {
  position: relative;
  z-index: 1;  /* 确保内容位于背景之上 */
  color: white;
  text-align: center;
  padding-top: 50px;
}

.background-overlay::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);  /* 半透明黑色覆盖 */
  z-index: 0;  /* 保证覆盖层在背景之下 */
}

/* 添加的悬停效果 */
.content .el-button:hover {
  background-color: blue !important;  /* 改变悬停时的背景颜色 */
  border-color: blue !important;  /* 改变悬停时的边框颜色 */
  color: white !important;  /* 改变悬停时的文字颜色 */
}
</style>