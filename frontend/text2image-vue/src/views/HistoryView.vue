<template>
    <div class="history-container">
        <h2>历史记录</h2>
         <div   v-if="historyRecords&&historyRecords.length" class="image-gallery-container">
            <div v-for="image in historyRecords" :key="image.id"  class="image-card">
                <img :src="image.url" :alt="image.name" class="image">                 
            </div>
        </div>
        
        <!-- 提示没有收藏 -->
        <div v-else>
            <img :src="require('@/assets/nofavorites.png')" alt="暂无收藏">
            <h1>暂无收藏</h1>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            historyRecords: [],
            //currentPage: 1,
            //pageSize: 10,
            //totalRecords: 0,
            //isLoading: true,
            token: localStorage.getItem('token') || '', 
        };
    },
    methods: {
       async getHistoryImages() {
            try {
                const response = await fetch("http://localhost:8080/auth/user/images", {
                    method: 'GET',
                    headers: {
                        Authorization: this.token,  // 携带 token
                    },
                })

               //const data = await response.json();
                //const images = data.images;

                // 解析JSON字符串
                const data = await response.json();
                console.log(data);

                const images=data.images;
                //console.log(Array.isArray(images)); // 输出: true，表示images是一个数组

                images.forEach(item => {
                   
                //     this.historyRecords.push({
                //         id: item.id,
                //         name: item.username,
                //         grams: item.params,
                //         url: item.picture,              
                //         likecount: item.likecount,
                //         createtime: item.create_time
                // })
                // console.log(this.historyRecords);
                if (!this.historyRecords.some(record => record.id === item.id)) {
                    this.historyRecords.push({
                        id: item.id,
                        name: item.username,
                        grams: item.params,
                        url: item.picture,              
                        likecount: item.likecount,
                        createtime: item.create_time
                    });
                } else {
                    console.log(`记录 id: ${item.id} 已经存在，跳过添加`);
                }
            });
            } catch (error) {
                console.error('获取收藏的图片失败:', error.response?.data || error.message);
            }
        },

        deleteRecord(id) {
            axios.delete(`/api/history-records/${id}`)
                .then(() => {
                    this.fetchHistoryRecords();
                    this.$message.success('记录已删除');
                })
                .catch(error => {
                    console.error('删除记录失败:', error);
                    this.$message.error('删除记录失败');
                });
        },
        viewImage(record) {
            window.open(record.image_url, '_blank');
        },
        downloadImage(record) {
            const link = document.createElement('a');
            link.href = record.image_url;
            link.download = 'downloaded_image.png'; // 你可以根据需要更改文件名
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        },
        
    },
    mounted() {
        this.getHistoryImages();
        //this.fetchHistoryRecords(this.currentPage, this.pageSize);
    }
};
</script>

<style scoped>
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
.history-container {
    padding: 20px;
    margin: 0 auto;
    text-align: center;
}

.history-records {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 20px;
}

.record {
    border: 1px solid #ccc;
    padding: 10px;
    border-radius: 4px;
    width: calc(33.333% - 20px);
    text-align: center;
}

.history-image {
    display: block;
    max-width: 100%;
    height: auto;
    margin: 0 auto;
}

.record-details {
    margin-top: 10px;
}
</style>