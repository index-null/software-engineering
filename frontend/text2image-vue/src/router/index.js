import Vue from 'vue'
import VueRouter from 'vue-router'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'

Vue.use(VueRouter)

const routes = [{
    path: '/',
    name: 'LoginView',
    component: LoginView
}, {
    path: '/register',
    name: 'register',
    component: RegisterView
}, {
    path: '/login',
    name: 'login',
    component: LoginView
}, {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
        import ( /* webpackChunkName: "about" */ '../views/AboutView.vue')
}, {
    path: '/home',
    name: 'home',
    component: () =>
        import ('../views/MainView.vue')
}, {
    path: '/info',
    name: 'info',
    component: () =>
        import ('../views/UserProfileViews.vue')

}]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router