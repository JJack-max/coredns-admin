import Vue from 'vue'
import axios from 'axios'
import App from './App.vue'
import router from './router'
import '@/assets/cutom.scss'
import './plugins/element.js'
import { decode } from 'js-base64'

Vue.config.productionTip = false

Vue.prototype.$ajax = axios

Vue.prototype.$payloadDecode = function (authToken) {
  if (!authToken || authToken.constructor !== String) {
    return null
  }
  const payload = authToken.split('.')[1]
  if (!payload) {
    return null
  }
  return JSON.parse(decode(payload))
}

const vm = new Vue({
  router,
  render: h => h(App)
}).$mount('#app')

process.env.NODE_ENV === 'development' && (axios.defaults.baseURL = 'http://localhost:8088')
// axios.defaults.headers.common['Authorization'] = store.state.authToken
axios.interceptors.request.use(
  function (config) {
    config.headers = { Authorization: localStorage.jwtToken }
    return config
  })
axios.interceptors.response.use(function (response) {
  // 正确响应
  return response
}, function (error) {
  // 响应错误
  if (error.response.status === 401) {
    vm.$Notice.warn('授权失效，请重新登录')
    vm.$router.replace({
      name: 'Login',
      query: { next: vm.$route.path === '/login' ? vm.$route.query.next : vm.$route.path }
    })
  } else {
    return Promise.reject(error)
  }
})
