/*
 * @Author       : jayj
 * @Date         : 2021-12-13 15:08:25
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-13 15:17:22
 */
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
