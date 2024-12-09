import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import apiClient from './api/index'; // Axios-拦截器实例
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css';
Vue.config.productionTip = false
Vue.prototype.$AppName = "";
Vue.use(ElementUI);
Vue.prototype.$axios = apiClient;
new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app')