import 'virtual:svg-icons-register' // 导入SVG图标注册
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import TDesign from 'tdesign-vue-next'

import App from './App.vue'
import router from './router'

import 'tdesign-vue-next/es/style/index.css'
import './assets/base.css'
import ComSvgIcon from './components/ComSvgIcon.vue'

const app = createApp(App)
app.component('ComSvgIcon', ComSvgIcon)

app.use(createPinia())
app.use(router)
app.use(TDesign)

app.mount('#app')
