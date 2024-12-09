import Vue from 'vue';
import VueRouter from 'vue-router';
import RegisterView from '@/views/RegisterView.vue';
import LoginView from '@/views/LoginView.vue';
import AboutView from '@/views/AboutView.vue';




Vue.use(VueRouter)

const routes = [{
    path: '/',
    redirect: '/about',
}, {
    path: '/about',
    name: 'about',
    component: AboutView
}, {
    path: '/register',
    name: 'register',
    component: RegisterView
}, {
    path: '/login',
    name: 'login',
    component: LoginView
}, {
    path: '/main',
    name: 'main',
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
    component: () =>
        import ('@/views/LogRegView.vue')
}]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router