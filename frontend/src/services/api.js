import { useAuthStore } from '../store/auth.store.js'

const API_BASE = '/api'

let isRefreshing = false
let refreshPromise = null

async function request(path, options = {}) {
  const authStore = useAuthStore()
  const token = authStore.token

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const config = {
    ...options,
    headers
  }

  let response = await fetch(`${API_BASE}${path}`, config)

  if (response.status === 401 && !options._isRetry) {
    if (!isRefreshing) {
      isRefreshing = true
      refreshPromise = authStore.refresh()
        .catch(() => {
          window.location.href = '/login'
          throw new Error('Session expired')
        })
        .finally(() => {
          isRefreshing = false
          refreshPromise = null
        })
    }

    try {
      await refreshPromise
      const newToken = authStore.token
      headers['Authorization'] = `Bearer ${newToken}`
      response = await fetch(`${API_BASE}${path}`, {
        ...config,
        headers,
        _isRetry: true
      })
    } catch (err) {
      throw err
    }
  }

  if (!response.ok) {
    let errorMessage = 'Request failed'
    let errorData = null
    try {
      errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch {}
    const err = new Error(errorMessage)
    err.status = response.status
    err.data = errorData
    throw err
  }

  return response.json()
}

export { request }
