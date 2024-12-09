<template>
    <div>
        <el-button type="primary" icon="el-icon-delete" class="delete-button">批量管理</el-button>

        <!-- 显示图片 -->
        <div v-if="images.length > 0" class="image-gallery-container">
            <div v-for="image in images" :key="image.id" class="image-card" @mouseover="hoveredImage = image.id"
                @mouseleave="hoveredImage = null">
                <img :src="image.url" :alt="image.name" class="image">
                <div class="overlay" v-if="hoveredImage === image.id">
                    <button @click="toggleFavorite(image)">{{ image.isFavorite ? '取消收藏' : '收藏' }}</button>
                    <button @click="downloadImage(image)">下载</button>
                </div>
            </div>
        </div>

        <!-- 提示没有收藏 -->
        <div v-else>
            <img :src="require('@/assets/nofavorites.png')" >
            <h1>暂无收藏</h1>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            images: [],  // 存储用户的收藏图片
            hoveredImage: null,  // 用于追踪当前悬停的图片
            token: localStorage.getItem('token') || '',  // 获取用户的 token
        };
    },
    mounted() {
        this.getFavoritedImages();  // 组件挂载时获取用户的收藏图片
    },
    methods: {
        // 获取用户收藏的图片
        async getFavoritedImages() {
            try {      
                const token=this.token;
                if (!token) {
                    throw new Error('Token not found in localStorage');
                }
                console.log(token);
                const response =await axios.get('http://localhost:8080/auth/user/favoritedimages', {
                    headers: {
                        Authorization:  this.token,  // 使用从 localStorage 获取的 token
                        
                    },
                });
                this.images = //获取response中每一个image的url
                

            } catch (error) {
                console.error('获取收藏的图片失败:', error.response?.data || error.message);
            }
        },

        // 收藏图像
        async addFavoriteImage(imageId, imageUrl) {
            try {
                const response = await axios.post(
                    'http://localhost:8080/auth/addFavoritedImage',
                    { id: imageId, url: imageUrl },  // 只需传递 id 或 url
                    {
                        headers: {
                            Authorization: `Bearer ${this.token}`,  // 携带 token
                        },
                    }
                );
                if (response.status === 200) {
                    this.getFavoritedImages();  // 收藏成功后重新获取收藏列表
                    this.$message.success('收藏图像成功');
                }
            } catch (error) {
                console.error('收藏图像失败:', error.response?.data || error.message);
            }
        },

        // 取消收藏图像
        async removeFavorite(image) {
            try {
                const response = await axios.delete(
                    'http://localhost:8080/auth/deleteFavoritedImage', {
                    headers: {
                        Authorization: `Bearer ${this.token}`,
                    },
                    params: { id: image.id },  // 传递图像的收藏表id
                }
                );
                if (response.status === 200) {
                    this.images = this.images.filter(i => i.id !== image.id);  // 从收藏列表中移除已取消收藏的图像
                    this.$message.success('取消收藏成功');
                }
            } catch (error) {
                console.error('取消收藏失败:', error.response?.data || error.message);
            }
        },

        // 切换收藏状态
        toggleFavorite(image) {
            if (image.isFavorite) {
                // 如果当前是收藏状态，调用取消收藏
                this.removeFavorite(image);
            } else {
                // 如果当前不是收藏状态，调用添加收藏
                this.addFavoriteImage(image.id, image.url);
            }
            image.isFavorite = !image.isFavorite;  // 切换收藏状态
        },

        // 下载图像
        downloadImage(image) {
            const link = document.createElement('a');
            link.href = image.url;
            link.download = image.name;
            link.click();
        },
    },
};
</script>

<style scoped>
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
