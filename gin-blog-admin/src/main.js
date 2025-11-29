import '@unocss/reset/tailwind.css'
import 'uno.css'

import { createApp } from 'vue'
import App from './App.vue'
import { setupNaiveDiscreteApi, setupNaiveUnocss } from './utils'
import { setupRouter } from './router'
import { setupStore } from './store'

async function bootstrap() {
  const app = createApp(App)
  setupStore(app) // Highest priority
  setupNaiveUnocss()
  setupNaiveDiscreteApi()
  await setupRouter(app)
  app.mount('#app')
}

bootstrap()
