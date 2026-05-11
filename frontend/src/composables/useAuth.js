import { computed } from 'vue'
import { useAuthStore } from '../store/auth.store.js'

export function useAuth() {
  const authStore = useAuthStore()

  return {
    user: computed(() => authStore.user),
    isAuthenticated: computed(() => authStore.isAuthenticated),
    logout: () => authStore.logout()
  }
}
