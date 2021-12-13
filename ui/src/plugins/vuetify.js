/*
 * @Author       : jayj
 * @Date         : 2021-12-13 15:12:52
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-13 15:23:21
 */
import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import '@mdi/font/css/materialdesignicons.css'

Vue.use(Vuetify)
const opts = {
  icons: {
    iconfont: 'mdi'
  }
}
export default new Vuetify(opts)
