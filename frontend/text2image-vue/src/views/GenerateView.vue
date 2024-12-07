<template>
   <section>
  <div class="main-container">  
    <el-card class="image-container">
      <img :src="imageUrl" alt="Generated Image" v-if="imageUrl" class="generated-image" />
      <p v-else>生成的图片将显示在这里</p>
      <div class="action-buttons" v-if="imageUrl">
        <el-button type="danger" @click="discardImage">舍弃</el-button>
        <el-button type="success" @click="downloadImage">导出</el-button>
        <el-button type="primary" @click="saveImage">保存</el-button>
        <el-button type="info" @click="favoriteImage">收藏</el-button>
      </div>
    </el-card>

    <el-card class="form-container">
      <h2>参数设置</h2>
      <el-form
        :model="form"
        label-position="left"
        label-width="120px"
        class="form"
      >
        <el-form-item label="提示词">
          <el-input v-model="form.prompt" placeholder="请输入提示词"></el-input>
        </el-form-item>
        <el-form-item label="宽度">
          <el-input-number v-model="form.width" :min="128" :max="1024" placeholder="宽度"></el-input-number>
        </el-form-item>
        <el-form-item label="高度">
          <el-input-number v-model="form.height" :min="128" :max="1024" placeholder="高度"></el-input-number>
        </el-form-item>
        <el-form-item label="步数">
          <el-input-number v-model="form.steps" :min="1" :max="100" placeholder="步数"></el-input-number>
        </el-form-item>
        <el-form-item label="采样方法">
          <el-select v-model="form.samplingMethod" placeholder="选择采样方法">
            <el-option label="DDIM" value="ddim"></el-option>
            <el-option label="PLMS" value="plms"></el-option>
            <el-option label="K-LMS" value="k-lms"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="种子">
          <el-input v-model="form.seed" placeholder="种子"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generateImage">生成图片</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
  </section>
</template>

<script>
import axios from 'axios';

export default {

  data() {
    return {
      imageUrl: 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202407272335308.gif',
      form: {
        prompt: '',
        width: 512,
        height: 512,
        steps: 50,
        samplingMethod: 'ddim',
        seed: ''
      },
    };
  },
  methods: {
    goToHistory() {
      this.$router.push('/HistoryPage');
    },
    generateImage() {
      const params = {
        prompt: this.form.prompt,
        width: this.form.width,
        height: this.form.height,
        steps: this.form.steps,
        sampling_method: this.form.samplingMethod,
        seed: this.form.seed
      };

      axios.post('/api/generate-image', params)
        .then(response => {
          if (response.data.success) {
            this.imageUrl = response.data.image_url;
          } else {
            this.$message.error(response.data.message);
          }
        })
        .catch(error => {
          console.error('生成图片失败:', error);
          if (error.response) {
            this.$message.error(`服务器错误: ${error.response.status} - ${error.response.data.message}`);
          } else if (error.request) {
            this.$message.error('请求未响应，请检查网络连接');
          } else {
            this.$message.error('请求发送失败，请稍后再试');
          }
        });
    },
    discardImage() {
      this.imageUrl = '';
      this.$message.success('图片已舍弃');
    },
    downloadImage() {
      const link = document.createElement('a');
      link.href = this.imageUrl;
      link.download = 'generated_image.png';
      link.click();
      this.$message.success('图片已下载');
    },
    saveImage() {
      // 假设这里有一个 API 用于保存图片到服务器
      axios.post('/api/save-image', { image_url: this.imageUrl })
        .then(response => {
          if (response.data.success) {
            this.$message.success('图片已保存');
          } else {
            this.$message.error(response.data.message);
          }
        })
        .catch(error => {
          console.error('保存图片失败:', error);
          if (error.response) {
            this.$message.error(`服务器错误: ${error.response.status} - ${error.response.data.message}`);
          } else if (error.request) {
            this.$message.error('请求未响应，请检查网络连接');
          } else {
            this.$message.error('请求发送失败，请稍后再试');
          }
        });
    },
    favoriteImage() {
      // 假设这里有一个 API 用于收藏图片
      axios.post('/api/favorite-image', { image_url: this.imageUrl })
        .then(response => {
          if (response.data.success) {
            this.$message.success('图片已收藏');
          } else {
            this.$message.error(response.data.message);
          }
        })
        .catch(error => {
          console.error('收藏图片失败:', error);
          if (error.response) {
            this.$message.error(`服务器错误: ${error.response.status} - ${error.response.data.message}`);
          } else if (error.request) {
            this.$message.error('请求未响应，请检查网络连接');
          } else {
            this.$message.error('请求发送失败，请稍后再试');
          }
        });
    }
  }
};
</script>

<style scoped>
.main-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  height: 100vh;
  padding: 20px;
  background-color: #f5f5f5;
  /* background-image: url('https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202411251503682.png');
  background-size: cover;
  background-position: center; */
}

.image-container {
  width: 50%;
  max-width: 600px;
  text-align: center;
}

.generated-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: space-around;
}

.form-container {
  width: 40%;
  max-width: 400px;
  padding: 20px;
}

.form {
  margin: 0 auto;
  width: 100%;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333333;
}
</style>