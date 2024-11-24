<template>
  <div class="history-container">
    <h2>历史记录</h2>
    <!-- 没有历史记录时显示默认图片 -->
    <div v-if="historyRecords.length === 0" class="no-records">
      <img src="https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202407272335308.gif" alt="No History" class="default-image" />
      <p>目前没有历史记录噢，快去生成一个吧！</p>
    </div>
    <div class="history-records" v-else>
      <div class="record" v-for="record in historyRecords" :key="record.id">
        <img :src="record.image_url" alt="History Image" class="history-image" />
        <div class="record-details">
          <p>提示词: {{ record.prompt }}</p>
          <p>生成时间: {{ record.created_at }}</p>
          <el-button type="text" @click="regenerateImage(record)">重新生成</el-button>
          <el-button type="text" @click="deleteRecord(record.id)">删除</el-button>
          <el-button type="text" @click="viewImage(record)">查看</el-button>
          <el-button type="text" @click="downloadImage(record)">下载</el-button>
        </div>
      </div>
    </div>
    <div class="pagination">
      <el-button @click="fetchHistoryRecords(currentPage - 1, pageSize)" :disabled="currentPage <= 1">上一页</el-button>
      <span>第 {{ currentPage }} 页，共 {{ Math.ceil(totalRecords / pageSize) }} 页</span>
      <el-button @click="fetchHistoryRecords(currentPage + 1, pageSize)" :disabled="currentPage >= Math.ceil(totalRecords / pageSize)">下一页</el-button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      historyRecords: [],
      currentPage: 1, 
      pageSize: 10, 
      totalRecords: 0, 
    };
  },
  methods: {
    fetchHistoryRecords(page, size) {
      axios.get('/api/history-records', { params: { page, size } })
      .then(response => {
          this.historyRecords = response.data.records;
          this.totalRecords = response.data.total; // 假设后端返回总记录数
      })
      .catch(error => {
          console.error('获取历史记录失败:', error);
          this.$message.error('获取历史记录失败');
      });
    },
    
    regenerateImage(record) {
      axios.post(`/api/history-records/regenerate/${record.id}`)
        .then(() => {
          this.$message.success('图片重新生成成功');
          // 可能需要重新获取历史记录
          this.fetchHistoryRecords(this.currentPage, this.pageSize);
        })
        .catch(error => {
          console.error('重新生成图片失败:', error);
          this.$message.error('重新生成图片失败');
        });
    },

    deleteRecord(id) {
      axios.delete(`/api/history-records/${id}`)
        .then(() => {
          this.fetchHistoryRecords(this.currentPage, this.pageSize); 
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
      link.download = 'downloaded_image.png'; 
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
  },
  mounted() {
    this.fetchHistoryRecords(this.currentPage, this.pageSize);
  }
};
</script>

<style scoped>
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

.no-records {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.default-image {
  max-width: 80%;
  height: 300px;
  width: 300px;
}
</style>