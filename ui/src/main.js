/*
 * @Author       : jayj
 * @Date         : 2021-11-15 14:49:15
 * @Description  :
 */
import { createApp } from 'vue'
import ArcoVue from '@arco-design/web-vue'
import App from './App.vue'
import router from './router'
import store from './store'
import '@arco-design/web-vue/dist/arco.css'

createApp(App).use(store).use(router).use(ArcoVue, {
  componentPrefix: 'arco'
}).mount('#app')
