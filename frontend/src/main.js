import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './store/auth.store.js'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Attempt silent token refresh before mounting
;(async () => {
  const authStore = useAuthStore()

  try {
    // Try to refresh using httpOnly cookie
    // If this succeeds, the store is hydrated with a new access token
    await authStore.refresh()
  } catch {
    // If refresh fails (no cookie, expired, network error), clear auth state
    // This ensures a clean state before the router guard runs
    authStore.clearAuth()
  }

  // Now mount — router guard will see the correct auth state
  app.mount('#app')
})()
