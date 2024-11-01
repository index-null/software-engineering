<script>
// 引入ref函数，用于创建响应式数据
import { ref } from 'vue';

export default {
  data() {
    // 创建一个响应式数据size，初始值为'medium'
    const size = ref('medium');
    // 返回size
    return {
      size:''
    };
  }
};
</script>

<template>
  <!-- 使用el-radio-group组件，v-model绑定size -->
  <el-radio-group v-model="size">
    <!-- 默认选项 -->
    <el-radio label="">默认</el-radio>
    <!-- 中等选项 -->
    <el-radio label="medium">中等</el-radio>
    <!-- 小型选项 -->
    <el-radio label="small">小型</el-radio>
    <!-- 超小选项 -->
    <el-radio label="mini">超小</el-radio>
  </el-radio-group>

  <!-- 使用el-descriptions组件，class为margin-top，column为3，size绑定size，border为true -->
  <el-descriptions class="margin-top" title="带边框列表" :column="3" :size="size" border>
    <!-- 插槽，用于添加额外内容 -->
    <template slot="extra">
      <!-- 添加一个操作按钮 -->
      <el-button type="primary" size="small">操作</el-button>
    </template>
    <!-- 用户名 -->
    <el-descriptions-item>
      <!-- 插槽，用于添加标签 -->
      <template slot="label">
        <!-- 添加一个用户图标 -->
        <i class="el-icon-user"></i>
        用户名
      </template>
      kooriookami
    </el-descriptions-item>
    <!-- 手机号 -->
    <el-descriptions-item>
      <template slot="label">
        <i class="el-icon-mobile-phone"></i>
        手机号
      </template>
      18100000000
    </el-descriptions-item>
    <!-- 居住地 -->
    <el-descriptions-item>
      <template slot="label">
        <i class="el-icon-location-outline"></i>
        居住地
      </template>
      苏州市
    </el-descriptions-item>
    <!-- 备注 -->
    <el-descriptions-item>
      <template slot="label">
        <i class="el-icon-tickets"></i>
        备注
      </template>
      <!-- 添加一个标签 -->
      <el-tag size="small">学校</el-tag>
    </el-descriptions-item>
    <!-- 联系地址 -->
    <el-descriptions-item>
      <template slot="label">
        <i class="el-icon-office-building"></i>
        联系地址
      </template>
      江苏省苏州市吴中区吴中大道 1188 号
    </el-descriptions-item>
  </el-descriptions>