<template>
    <div class="usage-guide">
      <h1>文绘星河 - 使用指南</h1>
  
      <!-- 界面介绍 -->
      <section class="guide-section">
        <h2>一、界面介绍</h2>
        <ul>
          <li><strong>表单容器（Form Container）</strong>: 包含用户输入和选择的区域。</li>
          <li><strong>结果容器（Result Container）</strong>: 展示生成的图片及操作按钮。</li>
        </ul>
      </section>
  
      <!-- 详细步骤 -->
      <section class="guide-section">
        <h2>二、详细步骤</h2>
  
        <!-- 表单部分 -->
        <article>
          <h3>1. 表单部分</h3>
          <dl>
            <dt>画面描述（Prompt）</dt>
            <dd>输入想要生成的画面描述，最大长度为100个字符。</dd>
            <dt>API调用方式（API Mode）</dt>
            <dd>选择“远程API”或“本地API”。</dd>
            <dt>尺寸（Size）</dt>
            <dd>选择生成图片的尺寸。</dd>
            <dt>步数（Steps）</dt>
            <dd>选择生成图片时使用的步数。</dd>
            <dt>种子（Seed）</dt>
            <dd>输入整数作为随机种子，点击刷新图标可生成随机种子。</dd>
            <dt>提交按钮（Submit Button）</dt>
            <dd>点击按钮提交表单，开始生成图片。</dd>
          </dl>
        </article>
  
        <!-- 结果部分 -->
        <article>
          <h3>2. 结果部分</h3>
          <dl>
            <dt>应用名称（App Name）</dt>
            <dd>显示当前应用名称。</dd>
            <dt>再次生成（Regenerate）</dt>
            <dd>点击按钮重新生成新的图片。</dd>
            <dt>生成的图片展示（Generated Images）</dt>
            <dd>每张图片会显示在结果区域内，包含信息及操作选项。</dd>
          </dl>
        </article>
      </section>
  
      <!-- 注意事项 -->
      <section class="guide-section">
        <h2>三、注意事项</h2>
        <ul>
          <li>确保输入有效的整数作为种子值。</li>
          <li>每次生成图片消耗20积分，请合理使用。</li>
          <li>遇到问题或有疑问，通过“使用指南”链接获取更多信息。</li>
        </ul>
      </section>
    </div>
  </template>
  
  <script>
  export default {
    name: 'UsageGuide'
  }
  </script>
  
  <style scoped>
  .usage-guide {
    font-family: Arial, sans-serif;
    color: #333;
    padding: 2rem;
  }
  
  h1, h2, h3 {
    color: #069;
  }
  
  .guide-section {
    margin-bottom: 2rem;
  }
  
  .guide-section h2 {
    border-bottom: 1px solid #ddd;
    padding-bottom: .5rem;
  }
  
  dl {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: .5rem;
  }
  
  dt {
    font-weight: bold;
  }
  
  dd {
    margin: 0;
  }
  
  ul {
    list-style-type: disc;
    padding-left: 1.5rem;
  }
  
  @media (max-width: 768px) {
    dl {
      grid-template-columns: 1fr;
    }
  }
  </style>