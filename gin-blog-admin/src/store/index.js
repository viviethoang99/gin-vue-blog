import { createPinia } from 'pinia'

// https://github.com/prazdevs/pinia-plugin-persistedstate
// Pinia data persistence, solving the problem of data loss on refresh
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export function setupStore(app) {
  const pinia = createPinia()
  pinia.use(piniaPluginPersistedstate)
  app.use(pinia)
}

export * from './modules/permission'
export * from './modules/tag'
export * from './modules/theme'
export * from './modules/user'
export * from './modules/auth'
