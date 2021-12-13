/*
 * @Author       : jayj
 * @Date         : 2021-12-13 15:12:52
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-13 15:36:35
 */
import Vue from 'vue'
import VueRouter from 'vue-router'
import Conf from '../views/Conf.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  // {
  //   path: "/about",
  //   name: "About",
  //   // route level code-splitting
  //   // this generates a separate chunk (about.[hash].js) for this route
  //   // which is lazy-loaded when the route is visited.
  //   component: () =>
  //     import(/* webpackChunkName: "about" */ "../views/About.vue"),
  // },
];

const router = new VueRouter({
  mode: 'history',
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }

  next()
})

export default router
