<template>
  <div class="main-container">
    <div class="form-container">
      <div class="form-header">
        <div class="form-title">文字作画</div>
        <div class="tutorial">使用指南</div>
      </div>
      <div class="form-body">
        <div class="form-item">
          <div class="form-label">文绘星河</div>
          <el-input
            type="textarea"
            :rows="4"
            placeholder="请输入您的描述"
            v-model="form.prompt"
          />
        </div>
        <div class="form-item">
          <div class="form-label">尺寸</div>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="form-sub-label">宽度</div>
              <el-input-number v-model="form.width" :min="1" :max="1000" placeholder="宽度" />
            </el-col>
            <el-col :span="12">
              <div class="form-sub-label">高度</div>
              <el-input-number v-model="form.height" :min="1" :max="1000" placeholder="高度" />
            </el-col>
          </el-row>
        </div>
        <div class="form-item">
          <div class="form-label">步数</div>
          <el-select v-model="form.steps" placeholder="请选择步数">
            <el-option
              v-for="item in stepsOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </div>
        <div class="form-item">
          <div class="form-label">种子</div>
          <el-input v-model="form.seed" placeholder="请输入种子值" />
        </div>
        <div class="form-submit">
          <el-button type="primary" native-type="submit">生成</el-button>
        </div>
      </div>
    </div>
    <div class="result-container">
      <div class="result-header">
        <div class="appName">文绘星河</div>
        <div class="regenerate">
          <button>再次生成</button>
        </div>
      </div>
      <div class="result-content" v-if="imageUrl">
        <div class="prompt-show">{{ this.form.prompt }}</div>
        <img :src="imageUrl" alt="Generated Image" class="generated-image" />
      </div>
      <div v-else class="placeholder">生成的图片将在这里显示</div>
    </div>
  </div>
</template>
<script>
import axios from 'axios';

export default {
  data() {
    return {
      form: {
        prompt: '',
        width: 512,
        height: 512,
        steps: 10,
        seed: ''
      },
      stepsOptions: [
        { value: 10, label: '10' },
        { value: 15, label: '15' },
        { value: 20, label: '20' },
        { value: 25, label: '25' },
        { value: 30, label: '30' },
        { value: 35, label: '35' },
        { value: 40, label: '40' }
      ],
      imageUrl: 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202411282349099.png' // 添加图片URL字段
    };
  },
  methods: {
    async handleSubmit() {
      try {
        // 模拟发送请求
        const response = await axios.post('http://localhost:8080/auth/generate', this.form);
        this.imageUrl = response.data.imageUrl;
      } catch (error) {
        console.error('生成图片失败', error);
      }
    }
  }
};
</script>
<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.generated-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.placeholder {
  font-size: 16px;
  color: #909399;
  text-align: center;
  padding: 20px;
}

.main-container {
  height: 98vh;
  width: 100%;
  display: flex;
  justify-content: space-between;
  gap: 1vw;
  background-color: #f0f2f5;
}

.form-container {
  flex: 1;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow-y: auto;
}

.result-container {
  flex: 2;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-title {
  font-size: 22px;
  font-weight: bold;
}

.tutorial {
  font-size: 14px;
  font-weight: bold;
  color: #5f6163;
  border: 1px gray solid;
  border-radius: 30px;
  padding: 6px;
  padding-left: 20px;
  padding-right: 20px
}

.form-body {
  display: flex;
  flex-direction: column;
  height: calc(100% - 60px); /* Adjust height to accommodate header and footer */
}

.form-item {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}

.form-label {
  font-size: 14px;
  margin-bottom: 5px;
}

.form-sub-label {
  font-size: 14px;
  margin-bottom: 5px;
}

.form-submit {
  display: flex;
  justify-content: center;
  width: 90%;
  margin: 0 auto;
}
</style>