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
    <el-select v-model="selectedSize" placeholder="请选择尺寸" @change="updateSize">
      <el-option
        v-for="item in sizeOptions"
        :key="item.label"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
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
          <el-button type="primary" native-type="submit"  @click="handleSubmit">生成</el-button>
        </div>
      </div>
    </div>
    <div class="result-container">
    <div class="result-header">
      <div class="appName">文绘星河</div>
      <div class="regenerate">
        <button @click="regenerateImage">再次生成</button>
      </div>
    </div>
    <div class="result-content" v-loading="loading">
      <div v-for="(img, index) in temp_generatedImg_results" :key="index" class="image-card">
        <el-image
          style="width: 100%; height: 100%; border-radius: 8px;"
          :src="img.img_url"
          fit="contain"
          lazy
        />
      </div>
    </div>
    <div v-if="temp_generatedImg_results.length === 0" class="placeholder">生成的图片将在这里显示</div>
  </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      form: {
        prompt: 'cute girl with a kite',
        width: 1024,
        height: 1024,
        steps: 10,
        seed: 1024
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
      temp_generatedImg_results: [],
      loading: false,
      sizeOptions: [
      { label: '1024x1024', value: '1024x1024' },
      { label: '720x1280', value: '720x1280' },
      { label: '768x1152', value: '768x1152' },
      { label: '1280x720', value: '1280x720' }
    ],
    selectedSize: '1024x1024'
    };
  },
  methods: {
    updateSize() {
  const selected = this.sizeOptions.find(option => option.value === this.selectedSize);
  if (selected) {
    console.log('Old width:', this.form.width, 'Old height:', this.form.height);
    this.form.width = parseInt(selected.value.split('x')[0], 10);
    this.form.height = parseInt(selected.value.split('x')[1], 10);
    console.log('New width:', this.form.width, 'New height:', this.form.height);
  }
},
    regenerateImage() {
      this.loading = true;
      this.handleSubmit().finally(() => {
        this.loading = false;
      });
    },
    handleSubmit() {
      // 强制将 seed 转换为整数
      this.form.seed = parseInt(this.form.seed, 10) || 0;
  
      // 确保转换后的 seed 是有效的整数
      if (isNaN(this.form.seed)) {
        this.$message.error('种子值必须是有效的整数');
        return;
      }
      this.$message.success('提交成功,正在生成图片...');
      const currentScore = parseInt(localStorage.getItem("score"), 10) || 0;
      localStorage.setItem("score", currentScore - 20);
      // 添加占位图片
      const placeholderImg = {
        prompt: this.form.prompt,
        width: this.form.width,
        height: this.form.height,
        seed: this.form.seed,
        steps: this.form.steps,
        img_url: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==' // 灰色占位图
      };
      this.temp_generatedImg_results.push(placeholderImg);
  
      this.loading = true;

      return this.$axios.post('http://localhost:8080/auth/generate', this.form, {
        timeout: 300000 // 设置超时时间为300秒
      }).then(response => {
        if (response && response.data) {
          console.log(response.data);
          let img_item = {
            "prompt": this.form.prompt,
            "width": this.form.width,
            "height": this.form.height,
            "seed": this.form.seed,
            "steps": this.form.steps,
            "img_url": response.data.image_url,
          };
          console.log(img_item);

          // 替换占位图片
          const lastIdx = this.temp_generatedImg_results.length - 1;
          this.$set(this.temp_generatedImg_results, lastIdx, img_item);

          this.$message.success(response.data.message);
        } else {
          this.$message.error('服务器返回数据异常');
        }
      }).catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      }).finally(() => {
        this.loading = false;
      });
    }
  }
};
</script>
<style scoped>
/* 主容器样式 */
.main-container {
  height: 98vh;
  width: 100%;
  display: flex;
  justify-content: space-between;
  gap: 1vw;
  background-color: #f0f2f5;
}

/* 表单容器样式 */
.form-container {
  flex: 1;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow-y: auto;
}

/* 表单头部样式 */
.form-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* 表单标题样式 */
.form-title {
  font-size: 22px;
  font-weight: bold;
}

/* 教程样式 */
.tutorial {
  font-size: 14px;
  font-weight: bold;
  color: #5f6163;
  border: 1px gray solid;
  border-radius: 30px;
  padding: 6px 20px;
}

/* 表单主体样式 */
.form-body {
  display: flex;
  flex-direction: column;
  height: calc(100% - 60px); /* Adjust height to accommodate header and footer */
}

/* 表单项样式 */
.form-item {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}

/* 表单标签样式 */
.form-label {
  font-size: 18px;
  margin-bottom: 5px;
  text-align: left;
}

/* 子标签样式 */
.form-sub-label {
  font-size: 14px;
  margin-bottom: 5px;
}

/* 输入框样式 */
.el-input,
.el-input-number,
.el-select {
  width: 100%;
}

/* 数字输入框样式 */
.el-input-number {
  width: 100%;
}

/* 选择框样式 */
.el-select {
  width: 100%;
}

/* 选项样式 */
.el-option {
  width: 100%;
}

/* 提交按钮样式 */
.form-submit {
  display: flex;
  justify-content: center;
  width: 90%;
  margin: 0 auto;
}

/* 结果容器样式 */
.result-container {
  flex: 2;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

/* 结果头部样式 */
.result-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* 应用名称样式 */
.appName {
  font-size: 22px;
  font-weight: bold;
}

/* 重新生成按钮样式 */
.regenerate button {
  font-size: 14px;
  font-weight: bold;
  color: #5f6163;
  border: 1px gray solid;
  border-radius: 30px;
  padding: 6px 20px;
  background-color: transparent;
}

/* 结果内容样式 */
.result-content {
  display: flex;
  flex-wrap: wrap; /* 允许换行 */
  align-items: flex-start;
  padding: 20px;
}

/* 生成的图片样式 */
.generated-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

/* 占位符样式 */
.placeholder {
  font-size: 16px;
  color: #909399;
  text-align: center;
  padding: 20px;
}

/* 其他未指定样式的组件类 */
.el-row,
.el-col {
  width: 100%;
}

.el-button {
  width: 100%;
}
/* 图片卡片样式 */
.image-card {
  width: 200px; /* 根据需要调整宽度 */
  height: 200px; /* 根据需要调整高度 */
  margin: 10px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 结果内容样式 */
.result-content {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
}
</style>