import Vue from 'vue';
import VueRouter from 'vue-router';
import {
    Message
} from 'element-ui'; // 导入 Message 组件
import AboutView from '@/views/AboutView.vue';




Vue.use(VueRouter)

const routes = [{
        path: '/',
        redirect: '/about',
    }, {
        path: '/about',
        name: 'about',
        component: AboutView
    },
    //{
    //     path: '/register',
    //     name: 'register',
    //     component: RegisterView
    // }, {
    //     path: '/login',
    //     name: 'login',
    //     component: LoginView
    // },
    {
        path: '/main',
        name: 'main',
        meta: {
            requiresAuth: true
        },
        redirect: '/main/explore',
        component: () =>
            import ('@/views/InnerView.vue'),
        children: [{
            path: 'generate',
            name: 'generate',
            component: () =>
                import ('@/views/GenerateView.vue')
        }, {
            path: 'history',
            name: 'history',
            component: () =>
                import ('@/views/HistoryView.vue')
        }, {
            path: 'explore',
            name: 'explore',
            component: () =>
                import ('@/views/ExploreView.vue')
        }, {
            path: 'favorites',
            name: 'favorites',
            component: () =>
                import ('@/views/FavoritesView.vue')
        }, {
            path: 'setting',
            name: 'setting',
            component: () =>
                import ('@/views/SettingView.vue')
        }]
    }, {
        path: '/log-reg',
        name: 'log-reg', // 添加了名称
        component: () =>
            import ('@/views/LogRegView.vue')
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})
router.beforeEach((to, from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

    if (requiresAuth) {
        // 检查 localStorage 中是否有 token
        const token = localStorage.getItem('token');

        // 如果有 token，则允许访问
        if (token) {
            next();
        } else {
            // 如果没有 token，则重定向到登录页面
            Message.error('您还未登录，请先登录!'); // 使用 Message 组件
            next({
                path: '/log-reg',
            });
        }
    } else {
        next(); // 不需要认证的路由直接放行
    }
});
export default router