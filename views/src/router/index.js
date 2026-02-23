import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '@/views/home'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/records',
    component: Home,
    children: [
      {
        path: 'records',
        name: 'Record',
        component: () => import(/* webpackChunkName: "records" */ '../views/records.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import(/* webpackChunkName: "records" */ '../views/login.vue')

  }

]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})
router.beforeEach((to, from, next) => {
  window.from = from

  if (to.meta.title) {
    document.title = to.meta.title
  }
  if (to.matched.length === 0) { // 如果未匹配到路由
    from.path ? next({ path: from.path }) : next('/') // 如果上级也未匹配到路由则跳转主页面，如果上级能匹配到则转上级路由
    return
  }
  if (to.path !== '/login' && !(localStorage.jwtToken)) {
    next({ name: 'Login', query: { next: to.path } })
  } else {
    next()
  }
})
export default router
