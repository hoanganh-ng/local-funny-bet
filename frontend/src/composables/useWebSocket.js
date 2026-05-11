import { ref, onMounted, onUnmounted } from 'vue'

export function useWebSocket() {
  const ws = ref(null)
  const isConnected = ref(false)
  const reconnectTimeout = ref(null)
  const eventHandlers = new Map()

  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'

  const connect = () => {
    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        isConnected.value = true
        console.log('[WS] Connected')
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          const handlers = eventHandlers.get(data.type) || []
          handlers.forEach(handler => handler(data))
        } catch (err) {
          console.error('[WS] Failed to parse message:', err)
        }
      }

      ws.value.onerror = (error) => {
        console.error('[WS] Error:', error)
      }

      ws.value.onclose = () => {
        isConnected.value = false
        console.log('[WS] Disconnected')
        scheduleReconnect()
      }
    } catch (err) {
      console.error('[WS] Connection failed:', err)
      scheduleReconnect()
    }
  }

  const scheduleReconnect = () => {
    if (reconnectTimeout.value) return

    reconnectTimeout.value = setTimeout(() => {
      reconnectTimeout.value = null
      console.log('[WS] Reconnecting...')
      connect()
    }, 3000)
  }

  const disconnect = () => {
    if (reconnectTimeout.value) {
      clearTimeout(reconnectTimeout.value)
      reconnectTimeout.value = null
    }

    if (ws.value) {
      ws.value.close()
      ws.value = null
    }

    eventHandlers.clear()
    isConnected.value = false
  }

  const on = (eventType, handler) => {
    if (!eventHandlers.has(eventType)) {
      eventHandlers.set(eventType, [])
    }
    eventHandlers.get(eventType).push(handler)

    return () => {
      const handlers = eventHandlers.get(eventType)
      if (handlers) {
        const index = handlers.indexOf(handler)
        if (index > -1) {
          handlers.splice(index, 1)
        }
      }
    }
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected,
    on
  }
}
