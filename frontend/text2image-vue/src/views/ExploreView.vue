<template>
<div class="explore-container">
  <div class="carousel-img-container">
    <div class="carousel-img-bg">
      <div class="selector">
        <div class="selector-item">精致国潮</div>
        <div class="selector-item-unselected">国风水墨</div>
        <div class="selector-item-unselected">厚涂原画</div>
      </div>
      <div class="text-box">
        <h1>琼楼玉宇</h1>
        <h3>中式色彩国风美学</h3>
      </div>
      <div class="generate-same-style">
        <input class="search-input" type="text" placeholder="华丽的中国传统古代建筑,坐落在山脚下,山脚下有弯曲的河流流过,天空中飘着祥云" disabled>
        <div class="generate-button" @click="$router.push('/main/generate')">生成同款</div>
      </div>
    </div>
  </div>
  <div class="generate-entry-container">
    <div class="generate-entry">
      <div class="text1">文字作画</div>
      <div class="text2">
        文字描述生成画作
      </div>
      <div class="generate-button" @click="$router.push('/main/generate')">去生成</div>
    </div>
  </div>
  <div class="explore-ground-container">
    <h2 class="ground-title">图片广场</h2>
    <div class="image-grid">
      <div class="image-item" v-for="n in 16" :key="n">
        <img src="https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202408311347060.jpg" alt="Image" class="grid-image">
      </div>
    </div>
  </div>
</div>
</template>

<script>
export default {

}
</script>

<style scoped>
* {
    box-sizing: border-box; /* 确保所有元素的宽度和高度包括内边距和边框 */
    margin: 0;
    padding: 0;
}

.explore-container {
    display: flex;
    flex-direction: column;
    min-height: 120vh; /* 最小高度 */
    background-color: #F1F6FF;
}

.carousel-img-container {
    flex: 1.3; 
    border-radius: 1.5rem;
    position: relative;
    color: white;
}

.generate-entry-container {
    flex: 0.5; /* 减少 flex 值 */
    border-radius: 1.5rem;
    position: relative;
    color: white;
}
.explore-ground-container{
    flex: 8;
    border-radius: 1.5rem; /* 使用相对单位 */
    position: relative; /* 使子元素可以绝对定位 */
    color: white; /* 设置文本颜色为白色 */
}


.carousel-img-bg {
    background-image: url("@/assets/Picture-Square/1.png");
    background-size: cover;
    width: 100%;
    height: auto;
    top: 0;
    left: 0;
    border-radius: 1.5rem; /* 与容器圆角一致 */
}

.selector {
    display: flex;
    gap: 1rem; /* 使用相对单位 */
    margin-left: 2rem; /* 使用相对单位 */
}

.selector-item,
.selector-item-unselected {
    background-color: #BED7DA;
    border-radius: 3rem; /* 使用相对单位 */
    padding: 1rem 2rem; /* 使用相对单位 */
    color: #7FAEBD; /* 文本颜色设置为白色 */
    cursor: pointer; /* 添加鼠标悬停效果 */
    transition: background-color 0.3s; /* 添加过渡效果 */
}

.selector-item {
    background-color: #297790; /* 添加背景颜色 */
}

.selector-item:hover,
.selector-item-unselected:hover {
    background-color: #1f5a6b; /* 悬停时改变背景色 */
}

.generate-same-style {
    display: flex;
    margin-top: 2rem; /* 使用相对单位 */
    align-items: center; /* 垂直居中对齐 */
    margin-left: 6rem; /* 使用相对单位 */
    margin-right: 6rem; /* 使用相对单位 */
    margin-bottom: rem; /* 使用相对单位 */
    position: relative; /* 确保生成按钮和其他元素不重叠 */
    z-index: 1; /* 确保生成按钮在背景之上 */
    
}

.search-input {
    flex: 10;
    padding: 1rem 1.5rem;
    border: none;
    border-radius: 3rem; /* 圆角 */
    background-color: rgba(255, 255, 255, 0.5); /* 半透明背景 */
    color: #333;
    font-size: 1rem;
    outline: none; /* 去掉聚焦时的默认边框 */
    text-align: left; /* 确保placeholder左侧对齐 */
    transition: background-color 0.3s; /* 添加过渡效果 */
    margin-bottom: 2rem;
}

.search-input:focus {
    background-color: rgba(255, 255, 255, 1); /* 聚焦时背景变白 */
}

.generate-button {
    flex: 1;
    background-color: #ff6347; /* 添加按钮背景色 */
    color: white;
    border: none;
    border-radius: 3rem; /* 圆角 */
    padding: 1rem 2rem;
    cursor: pointer;
    transition: background-color 0.3s; /* 添加过渡效果 */
    margin-bottom: 2rem;
}

.generate-button:hover {
    background-color: #ff4500; /* 悬停时改变背景色 */
}

.generate-entry-container {
    flex: 1;
    display: flex;
    align-items: flex-end;
    margin-bottom: 2rem; /* 使用相对单位 */
    margin-top: 0.5rem; /* 使用相对单位 */
}

.generate-entry {
    width: 50%;
    height: 100%;
    background-image: url("@/assets/Picture-Square/it.gif");
    background-size: cover; /* 保持图片完整显示 */
    background-repeat: no-repeat; /* 防止图片重复 */
    background-position: center; /* 图片居中显示 */
    border-radius: 2.5rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    align-items: flex-start;
    padding: 2rem; /* 添加内边距 */
    box-sizing: border-box;
}

.text1 {
    color: black; /* 设置文本颜色为黑色 */
    margin-bottom: 1rem; /* 使用相对单位 */
}

.text2 {
    color: rgb(97, 90, 90);
    margin-bottom: 1rem; /* 使用相对单位 */
    font-size: 0.875rem; /* 使用相对单位 */
}

.generate-button {
    background-color: #007BFF; /* 蓝色背景 */
    color: white; /* 白色文字 */
    border: none;
    border-radius: 3rem;
    padding: 1rem 2rem;
    cursor: pointer;
    transition: background-color 0.3s; /* 添加过渡效果 */
}

.generate-button:hover {
    background-color: #0056b3; /* 悬停时改变背景色 */
}

.explore-ground-container {
    flex: 8;
    padding: 2rem; /* 使用相对单位 */
    box-sizing: border-box;
}

.ground-title {
    font-size: 2rem; /* 使用相对单位 */
    margin-bottom: 2rem; /* 使用相对单位 */
    text-align: left;
    color: black;
    font-family: Arial, Helvetica, sans-serif;
}

.image-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem; /* 使用相对单位 */
    justify-content: space-between;
}

.image-item {
    flex: 0 0 calc(25% - 1rem); /* 一行展示4张图片 */
    position: relative;
    overflow: hidden;
    border-radius: 1rem; /* 使用相对单位 */
}

.text-box {
    margin-top: 2rem; /* 使用相对单位 */
    display: flex;
    flex-direction: column;
    text-align: left;
    margin-left: 3rem; /* 使用相对单位 */
    margin-top: 9rem; /* 使用相对单位 */
}

.text-box h1 {
    font-size: 2rem; /* 使用相对单位 */
    margin-bottom: 0.5rem; /* 使用相对单位 */
}

.text-box h3 {
    font-size: 1.2rem; /* 使用相对单位 */
}
</style>