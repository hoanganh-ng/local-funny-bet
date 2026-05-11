import { ref, onMounted, onUnmounted } from 'vue'
import { leaderboardService } from '../services/leaderboard.service.js'
import { useWebSocket } from './useWebSocket.js'

export function useLeaderboard(leaderboardId) {
  const scores = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const { on } = useWebSocket()

  const fetchScores = async () => {
    isLoading.value = true
    error.value = null

    try {
      const data = await leaderboardService.getOne(leaderboardId.value || leaderboardId)
      scores.value = data.scores || []
    } catch (err) {
      error.value = err.message || 'Failed to fetch scores'
      console.error('Failed to fetch leaderboard scores:', err)
    } finally {
      isLoading.value = false
    }
  }

  let unsubscribe = null

  onMounted(() => {
    fetchScores()

    unsubscribe = on('leaderboard_updated', (data) => {
      const targetId = leaderboardId.value || leaderboardId
      if (data.leaderboard_id === targetId) {
        fetchScores()
      }
    })
  })

  onUnmounted(() => {
    if (unsubscribe) {
      unsubscribe()
    }
  })

  return {
    scores,
    isLoading,
    error,
    refetch: fetchScores
  }
}
