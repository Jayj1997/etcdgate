/*
 * @Author       : jayj
 * @Date         : 2021-12-13 15:12:52
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 10:37:00
 */
import Vue from 'vue'
import VueRouter from 'vue-router'
import Conf from '../views/Conf.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Conf',
    component: Conf,
    meta: {
      title: 'config'
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = 'etcd-gate ' + to.meta.title
  }

  next()
})

export default router
