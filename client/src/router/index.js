import Vue from 'vue'
import Router from 'vue-router'
import Editplus from '../components/Editplus.vue'
import Charles from '../components/Charles.vue'

Vue.use(Router)

export default new Router({
    mode: 'history',
    routes: [
        {
            path: '/',
            redirect: '/keygen/editplus',
        },
        {
            path: '/keygen/editplus',
            component: Editplus,
        },
        {
            path: '/keygen/charles',
            component: Charles,
        },
    ],
})
