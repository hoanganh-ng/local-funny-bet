import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '../services/auth.service.js'

export const useAuthStore = defineStore('auth', () => {
  // In-memory only - no localStorage
  const token = ref(null)
  const user = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  function setAuth(newToken, newUser) {
    token.value = newToken
    user.value = newUser
  }

  function clearAuth() {
    token.value = null
    user.value = null
  }

  async function refresh() {
    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' }
      })

      if (!response.ok) throw new Error('Refresh failed')

      const data = await response.json()

      // Backend returns { access_token }, fetch user separately
      const userResponse = await fetch('/api/auth/me', {
        headers: {
          'Authorization': `Bearer ${data.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (userResponse.ok) {
        const userData = await userResponse.json()
        setAuth(data.access_token, userData)
      } else {
        setAuth(data.access_token, null)
      }
    } catch (err) {
      clearAuth()
      throw err
    }
  }

  async function logout() {
    try {
      await authService.logout()
    } finally {
      clearAuth()
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    setAuth,
    clearAuth,
    refresh,
    logout
  }
})
