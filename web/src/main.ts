import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router.ts'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import VueCookies from 'vue-cookies'

import { createPinia } from 'pinia'


const app = createApp(App)

app.use(router)
app.use(ElementPlus)
app.use(createPinia())
app.config.globalProperties.VueCookies = VueCookies

app.mount('#app')
