// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueAnalytics from 'vue-analytics'
import './assets/custom.sass'

import * as Rx from 'rxjs'
import VueRx from 'vue-rx'

import App from './App.vue'
import router from './router'

Vue.config.productionTip = true

Vue.use(VueRx, Rx)
Vue.use(VueAnalytics, {
    id: 'UA-86109905-1',
    router,
})

/* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    template: '<App/>',
    components: {App},
})
