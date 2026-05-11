const API_BASE = '/api'

let isRefreshing = false
let refreshPromise = null

async function request(path, options = {}) {
  const token = localStorage.getItem('token')

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
      refreshPromise = fetch(`${API_BASE}/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' }
      })
        .then(async (res) => {
          if (!res.ok) throw new Error('Refresh failed')
          const data = await res.json()
          localStorage.setItem('token', data.token)
          return data.token
        })
        .catch(() => {
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          throw new Error('Session expired')
        })
        .finally(() => {
          isRefreshing = false
          refreshPromise = null
        })
    }

    try {
      const newToken = await refreshPromise
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
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch {}
    throw new Error(errorMessage)
  }

  return response.json()
}

export { request }
