<template>
  <div class="history-container">
    <h2>历史记录</h2>
    <div v-if="historyRecords.length === 0">没有历史记录。</div>
    <div class="history-records" v-else>
      <div class="record" v-for="record in historyRecords" :key="record.id" >
        <img :src="record.image_url" alt="History Image" class="history-image" />
        <div class="record-details">
          <p>提示词: {{ record.prompt }}</p>
          <p>生成时间: {{ record.created_at }}</p>
          <el-button type="text" @click="regenerateImage(record)">重新生成</el-button>
          <el-button type="danger" icon="el-icon-delete" circle @click="deleteRecord(record.id)"></el-button>
          <el-button type="text" @click="viewImage(record)">查看</el-button>
          <el-button type="text" @click="downloadImage(record)">下载</el-button>
      </div>
      </div>
    <div class="pagination">
      <el-button @click="fetchHistoryRecords(currentPage - 1, pageSize)" :disabled="currentPage <= 1">上一页</el-button>
      <span>第 {{ currentPage }} 页，共 {{ Math.ceil(totalRecords / pageSize) }} 页</span>
      <el-button @click="fetchHistoryRecords(currentPage + 1, pageSize)" :disabled="currentPage >= Math.ceil(totalRecords / pageSize)">下一页</el-button>
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
      currentPage: 1, 
      pageSize: 10, 
      totalRecords: 0, 
      isLoading: true,
    };
  },
  methods: {
      fetchHistoryRecords(page, size) {
          // axios.get('/api/history-records', { params: { page, size } })
          // .then(response => {
          //     this.historyRecords = response.data.records;
          //     this.totalRecords = response.data.total; // 假设后端返回总记录数
          // })
          // .catch(error => {
          //     console.error('获取历史记录失败:', error);
          //     this.$message.error('获取历史记录失败');
          // });
      // 模拟 API 调用
      this.currentPage = page; // 更新当前页码
      const mockData = {
        records: [],
          total: 100 
      };

      // 模拟分页逻辑
      const start = (page - 1) * size;
      const end = Math.min(start + size, mockData.total);
      for (let i = start; i < end && i < mockData.total; i++) {
      mockData.records.push({
          id: i + 1,
          image_url: `https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202408311347060.jpg`,
          prompt: 'Test prompt',
          created_at: '2024-01-' + (i + 1) 
      });
      }

      // 延迟模拟网络请求
      setTimeout(() => {
          this.historyRecords = mockData.records;
          this.totalRecords = mockData.total;
          this.isLoading = false;
      }, 1000); // 模拟1秒的网络延迟
          
      },
    regenerateImage() {
      // 重新生成图片的逻辑，可能需要重新调用 generateImage 方法，并传入相应的参数
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
  onPreviousPage() {
      this.fetchHistoryRecords(this.currentPage - 1, this.pageSize);
  },
  onNextPage() {
      this.fetchHistoryRecords(this.currentPage + 1, this.pageSize);
  }
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
</style>