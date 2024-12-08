<template>
  <div>
    <el-button type="primary" icon="el-icon-delete" class="delete-button"> 批量管理</el-button>

    <!-- 显示图片 -->
    <div v-if="images.length > 0" class="image-gallery-container">
      <div v-for="image in images" :key="image.id" class="image-card" @mouseover="hoveredImage = image.id"
        @mouseleave="hoveredImage = null">
        <img :src="image.url" :alt="image.name" class="image">
        <div class="overlay" v-if="hoveredImage === image.id">
          <button @click="toggleFavorite(image)">{{ image.isFavorite ? '收藏' : '取消收藏' }}</button>
          <button @click="downloadImage(image)">下载</button>
        </div>
      </div>
    </div>
    <!-- 提示没有收藏 -->
    <div v-else>
      <img :src="require('@/assets/nofavorites.png')" :alt="暂无收藏">
      <h1>暂无收藏</h1>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      images: [
        { id: 1, name: '图片1', url: require('@/assets/favorites/image1.png'), isFavorite: false },
        { id: 2, name: '图片2', url: require('@/assets/favorites/image2.png'), isFavorite: false },
        { id: 3, name: '图片3', url: require('@/assets/favorites/image3.png'), isFavorite: false },
        { id: 4, name: '图片4', url: require('@/assets/favorites/image4.png'), isFavorite: false },
        // 更多图片...
      ],
      hoveredImage: null,  // 用于追踪当前悬停的图片
    };
  },
  methods: {
    toggleFavorite(image) {
      image.isFavorite = !image.isFavorite;
      if (image.isFavorite) {
        this.$message.success('取消收藏成功');
        this.removeImageFromFavorites(image);
      } else {
        this.$message.success('收藏成功');

      }
    },
    downloadImage(image) {
      const link = document.createElement('a');
      link.href = image.url;
      link.download = image.name;
      link.click();
    },
  },
};
</script>

<style>
.delete-button {
  position: fixed;
  /* 固定位置 */
  top: 50px;
  /* 离页面顶部20px */
  right: 50px;
  /* 离页面右边20px */
  z-index: 900;
  /* 确保按钮显示在页面最上面 */
  border-radius: 30px;
}

.image-gallery-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 80px;
  padding: 20px;
  margin-left: 130px;
  margin-top: 100px;
  width: 100%;
  /* 容器宽度设置为页面宽度的 100% */
}

.image-card {
  position: relative;
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.image-card:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

.image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.image-card:hover .overlay {
  opacity: 1;
}

.overlay button {
  background: white;
  border: none;
  padding: 10px 20px;
  margin: 10px;
  cursor: pointer;
  border-radius: 5px;
}
</style>