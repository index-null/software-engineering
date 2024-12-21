<template>
  <div class="main-container">
    <div class="form-container">
      <div class="form-header">
        <div class="form-title">文字作画</div>
        <div class="tutorial">使用指南</div>
      </div>
      <div class="form-body">
        <div class="form-item">
          <div class="form-appname">文绘星河</div>
          <el-input
            type="textarea"
            autosize
            :rows="8"
            placeholder="试试输入你心中的画面,尽量描述具体,可以尝试一些风格修饰词辅助你的表达"
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
          <el-input v-model="form.seed" placeholder="请输入种子值">
            <template #suffix>
              <i class="el-input__icon el-icon-refresh" @click="generateRandomSeed"></i>
            </template>
          </el-input>
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
      <div class="image-card" v-for="(img, index) in temp_generatedImg_results" :key="index">
  <el-image
    style="width: 100%; height: 100%; border-radius: 8px;"
    :src="img.img_url"
    fit="contain"
    lazy
  />
  <div class="overlay">
    <div class="image-info">
      <p>Prompt: {{ img.prompt }}</p>
      <p>Width: {{ img.width }}</p>
      <p>Height: {{ img.height }}</p>
      <p>Seed: {{ img.seed }}</p>
      <p>Steps: {{ img.steps }}</p>
    </div>
    <div class="image-buttons">
      <el-button type="primary" icon="el-icon-edit" circle @click="reuseParameters(img)"></el-button>
      <el-button type="warning" icon="el-icon-star-off" circle @click="favoriteImage(img)"></el-button>
      <el-button type="danger" icon="el-icon-delete" circle @click="deleteImage(index)"></el-button>
    </div>
  </div>
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
        prompt: '',
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
      temp_generatedImg_results: [ {
            "prompt": "陈绮贞的专辑图片<<Pussy>>",
            "width": 1024,
            "height": 1024,
            "seed": 1024,
            "steps": 10,
            "img_url": `https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com/chuhsing/202408311347062.jpg`,
          }],
      loading: false,
      sizeOptions: [
      { label: '1024 x 1024px', value: '1024x1024' },
      { label: '720x1280px', value: '720x1280' },
      { label: '768x1152px', value: '768x1152' },
      { label: '1280x720px', value: '1280x720' }
    ],
    selectedSize: '1024x1024'
    };
  },
  methods: {
    reuseParameters(img) {
    this.form.prompt = img.prompt;
    this.form.width = img.width;
    this.form.height = img.height;
    this.form.seed = img.seed;
    this.form.steps = img.steps;
    this.selectedSize = `${img.width}x${img.height}`;
  },
  favoriteImage(img) {
    console.log('Favorite image:', img);
    this.$axios.post('http://localhost:8080/auth/addFavoritedImage', {url: img.img_url}, {
                        headers: {
                            'Content-Type': 'application/json', // 设置请求头
                        },                      
                    })
      .then(response => {
        if (response.status === 200) {
          this.$message.success('收藏成功');
        } else {
          this.$message.error(response.data.message);
        }
      })
      .catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      });
  },
  deleteImage(index) {
    this.temp_generatedImg_results.splice(index, 1);
  },
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
        img_url: 'https://via.placeholder.com/150'}
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
    },
    generateRandomSeed() {
      this.form.seed = Math.floor(Math.random() * 4369000);
    },
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
  font-family: Arial, Helvetica, sans-serif;
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
  font-size: 12px;
  font-weight: bolder;
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
.form-appname{
  text-align: left;
  font-size: 12px;
  font-weight: bolder;
  color: #000000;
  background-color: #F7F8FC;
  border-radius: 30px;
  padding: 6px 20px;
  margin-bottom: 5vh;
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
  flex: 3.5;
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

.result-content {
  display: flex;
  flex-wrap: wrap; /* 允许换行 */
  justify-content: flex-start; /* 从左侧开始 */
  align-items: flex-start;
  padding: 20px;
  gap: 10px; /* 设置间距 */
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

.image-card {
  position: relative; /* 添加相对定位 */
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease, opacity 0.3s ease; /* 添加 opacity 过渡 */
  height: 350px;
  width: 350px;
}

.image-card:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
  opacity: 0.9; /* 可选：调整悬停时的透明度 */
}

.overlay {
  position: absolute; /* 设置为绝对定位 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5); /* 半透明黑色背景 */
  display: flex;
  justify-content: space-around;
  align-items: center;
  opacity: 0; /* 初始透明度为0 */
  transition: opacity 0.3s ease; /* 添加过渡效果 */
  flex-direction: column; /* 确保内容垂直居中 */
}

.image-card:hover .overlay {
  opacity: 1; /* 鼠标悬停时透明度为1 */
}

.image-info {
  color: white;
  text-align: center;
  margin-bottom: 10px;
  overflow-y: auto; /* 添加垂直滚动 */
  max-height: 200px; /* 设置最大高度 */
}

.image-buttons {
  display: flex;
  gap: 10px;
}
/* 美化滚动条 */
.image-info::-webkit-scrollbar {
  width: 8px;
}

.image-info::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

.image-info::-webkit-scrollbar-thumb {
  background: rgba(100, 94, 94, 0.6);
  border-radius: 4px;
}

.image-info::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 1);
}

/* 结果内容样式 */
.result-content {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
}

/* 刷新图标样式 */
.el-input__icon.el-icon-refresh {
  cursor: pointer;
}

::v-deep .el-textarea__inner {
  background-color: #F7F8FC; /* 浅灰色背景 */
  color: #5951f2; /* 深灰色文本 */
  font-weight: bold; /* 加粗 */
  border-radius: 6px;
  font-family: Arial, Helvetica, sans-serif;
}
</style>