// src/store/index.js
import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        appName: '文绘星河' // 全局变量
    },
    mutations: {},
    actions: {
        updateAppName({
            commit
        }, newName) {
            commit('setAppName', newName);
        }
    },
    getters: {
        appName: state => state.appName
    }
});