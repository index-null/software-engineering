import Vue from 'vue';
import VueRouter from 'vue-router';
import { Message } from 'element-ui'; // 导入 Message 组件
import AboutView from '@/views/AboutView.vue';

Vue.use(VueRouter);

// 定义路由配置
const routes = [
    {
        path: '/',
        redirect: '/about', // 默认重定向到 /about 页面
    },
    {
        path: '/about',
        name: 'about',
        component: AboutView, // 关于页面
    },
    {
        path: '/main',
        name: 'main',
        meta: {
            requiresAuth: true // 需要认证
        },
        redirect: '/main/explore', // 默认重定向到 /main/explore 页面
        component: () => import('@/views/InnerView.vue'), // 内部视图组件
        children: [
            {
                path: 'generate',
                name: 'generate',
                component: () => import('@/views/GenerateView.vue') // 生成页面
            },
            {
                path: 'history',
                name: 'history',
                component: () => import('@/views/HistoryView.vue') // 历史页面
            },
            {
                path: 'explore',
                name: 'explore',
                component: () => import('@/views/ExploreView.vue') // 探索页面
            },
            {
                path: 'favorites',
                name: 'favorites',
                component: () => import('@/views/FavoritesView.vue') // 收藏页面
            },
            {
                path: 'setting',
                name: 'setting',
                component: () => import('@/views/SettingView.vue') // 设置页面
            }
        ]
    },
    {
        path: '/log-reg',
        name: 'log-reg',
        component: () => import('@/views/LogRegView.vue') // 登录注册页面
    },
    {
        path: '/usage',
        name: 'usage',
        component: () => import('@/views/UsageGuide.vue') // 使用指南页面
    }
];

// 创建 VueRouter 实例
const router = new VueRouter({
    mode: 'history', // 使用 HTML5 History 模式
    base: process.env.BASE_URL, // 基础路径
    routes // 路由配置
});

// 全局前置守卫
router.beforeEach((to, from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth); // 检查是否需要认证

    if (requiresAuth) {
        // 检查 localStorage 中是否有 token
        const token = localStorage.getItem('token');

        if (token) {
            next(); // 如果有 token，则允许访问
        } else {
            // 如果没有 token，则重定向到登录页面
            Message.error('您还未登录，请先登录!'); // 使用 Message 组件显示错误信息
            next({
                path: '/log-reg',
            });
        }
    } else {
        next(); // 不需要认证的路由直接放行
    }
});

export default router;