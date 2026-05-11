import { request } from './api.js'

export const authService = {
  login() {
    window.location.href = '/api/auth/google'
  },

  async logout() {
    await request('/auth/logout', { method: 'POST' })
  },

  async refresh() {
    return request('/auth/refresh', { method: 'POST' })
  }
}
