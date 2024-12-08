<template>
    <div class="logo-app-name-container">
      <div class="logo-container">
        <img src="../assets/button-icon/文生图-gray.svg" :alt="appName" class="logo"/>
      </div>
      <h1 class="app-name">{{ appName }}</h1>
    </div>
  </template>
  
  <script>
  export default {
    computed: {
      appName() {
        return this.$store.state.appName;
      }
    },
    methods: {
      
    }
  };
  </script>
  
  <style scoped>
  .logo-app-name-container {
    display: flex;
    align-items: center;
  }
  
  .logo-container {
    width: 50px;
    height: 50px;
    margin-right: 20px;
  }
  
  .logo {
    width: 100%;
    height: 100%;
    object-fit: contain;
    cursor: pointer;
  }
  
  .app-name {
    font-size: 3rem;
    margin: 0;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
    text-align: center;
    font-family: 'Arial', sans-serif;
    background: linear-gradient(270deg, rgb(255, 255, 255), rgba(150, 150, 150, 0.8), rgba(255, 255, 255));
    -webkit-background-clip: text;
    background-clip: text; /* 添加标准属性 */
    -webkit-text-fill-color: transparent;
    background-size: 200% 200%;
    animation: gradientAnimation 8s ease-in-out infinite;
    transform: translateY(0);
    animation: floatAnimation 4s ease-in-out infinite;
    cursor: pointer;
  }
  
  @keyframes gradientAnimation {
    0% {
      background-position: 0% 50%;
    }
    50% {
      background-position: 100% 50%;
    }
    100% {
      background-position: 0% 50%;
    }
  }
  
  @keyframes floatAnimation {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-5px);
    }
  }
  </style>