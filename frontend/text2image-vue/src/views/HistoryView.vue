<template>
  <div>
    <h2>历史记录</h2> 
    <div class="block">
      <span class="demonstration"></span>
      <el-date-picker
        v-model="value1"
        type="datetimerange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
      ></el-date-picker>
    <el-button type="success" icon="el-icon-check" circle @click="handleDateChange()"></el-button>
    </div>

    <div class="history-container">
      <div v-if="historyRecords && historyRecords.length" class="image-gallery-container">
        <div v-for="image in historyRecords" :key="image.id" class="image-card" @mouseover="hoveredImage = image.id" 
        @mouseleave="hoveredImage = null">
          <img :src="image.url" :alt="image.name" class="image">
          <!-- <div font="微软雅黑">图片ID:{{ image.id }} 点赞数量:{{ image.likecount }}</div> -->
          <div class="overlay" v-if="hoveredImage === image.id">
            <button  circle  @click="downloadImage(image)">下载图像</button>
            <button  round @click="deleteRecord(image)">删除</button> 
          </div>
        </div>
      </div>

      <!-- 提示没有收藏 -->
      <div v-else>
        <img :src="require('@/assets/nofavorites.png')" alt="暂无历史记录">
        <h1>暂无历史记录</h1>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      historyRecords: [],
      token: localStorage.getItem('token') || '',
      hoveredImage: null,  // 用于追踪当前悬停的图片
      pickerOptions: {
        shortcuts: [
          {
            text: '最近一周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit('pick', [start, end]);
            }
          },
          {
            text: '最近一个月',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit('pick', [start, end]);
            }
          },
          {
            text: '最近三个月',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit('pick', [start, end]);
            }
          }
        ]
      },
      value1: [new Date(2000, 10, 10, 10, 10), new Date(2000, 10, 11, 10, 10)],
      value2: ''
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
        });

        const data = await response.json();
        console.log(data);

        const images = data.images;
        images.forEach(item => {
          if (!this.historyRecords.some(record => record.id === item.id)) {
            this.historyRecords.push({
              id: item.id,
              name: item.username,
              params: item.params,
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
  formatDateToISO(date) {
  // 将传入的日期字符串转换为 Date 对象
  const dateObj = new Date(date);
  
  // 获取各个时间部分
  const year = dateObj.getUTCFullYear();
  const month = String(dateObj.getUTCMonth() + 1).padStart(2, '0'); // 月份从0开始，所以需要加1
  const day = String(dateObj.getUTCDate()).padStart(2, '0');
  const hours = String(dateObj.getUTCHours()).padStart(2, '0');
  const minutes = String(dateObj.getUTCMinutes()).padStart(2, '0');
  const seconds = String(dateObj.getUTCSeconds()).padStart(2, '0');
  const milliseconds = String(dateObj.getUTCMilliseconds()).padStart(3, '0');
  
  // 拼接成目标格式: 2006-01-02T15:04:05.000000Z
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
  },
    handleDateChange() {
      const { value1 } = this;
      if (value1 && value1.length === 2) {
        const [startDate, endDate] =value1;
        const start = this.formatDateToISO(startDate);  // 转换格式
        const end =  this.formatDateToISO(endDate);
        console.log('开始日期:', start);
        console.log('结束日期:', end);
        this.getPeriodImages(start, end);
      }
    },
    async getPeriodImages(start,end) {
      try {
        const response = await fetch("http://localhost:8080/auth/user/images/timeRange?start_time="+start+"&end_time="+end, {
          method: 'GET',
          headers: {
            Authorization: this.token,  // 携带 token
          },
          //params:{startDate,endDate},//传递开始时间和结束时间
        });
        console.log('get开始日期:', start);
        console.log('get结束日期:', end);

        const data = await response.json();
        console.log(data);

        const images = data.images;
        this.historyRecords = [];//置为空
        images.forEach(item => {
          if (!this.historyRecords.some(record => record.id === item.id)) {
            
            this.historyRecords.push({
              id: item.id,
              name: item.username,
              params: item.params,
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
    async deleteRecord(image) {
      try {
                const response = await axios.delete(
                    'http://localhost:8080/auth/root/deleteOneImage', {
                    headers: {
                        Authorization: this.token,
                    },
                    params: { url: image.url },  // 传递图像的收藏表url
                }
                );
                
                if (response.status === 200) {
                    this.historyRecords = this.historyRecords.filter(i => i.url !== image.url);  // 从收藏列表中移除已取消收藏的图像
                    this.$message.success('删除成功');
                }
          } catch (error) {
                console.error('删除失败:', error.response?.data || error.message);
          }
     
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
  }
};
</script>

<style scoped>
.overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 50%;
    top:50%;
    height: 50%;
    background: whitesmoke;
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    transition: opacity 0.3s ease;
}
/* 保持原有的样式 */
.h2 {
  text-align: left;
}
.image-gallery-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(365px, 1fr));
  gap: 80px;
  padding: 20px;
  margin-left: 100px;
  margin-top: 50px;
  width: 80%;
}

.image-card {
  position: relative;
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  height:600px;
  width:400px;
}

.image-card:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

.image {
  width: 90%;
  height: 90%;
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