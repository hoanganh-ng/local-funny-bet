import { ref, onMounted } from 'vue'
import { tournamentService } from '../services/tournament.service.js'

export function useTournament() {
  const tournament = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  const fetchTournament = async () => {
    isLoading.value = true
    error.value = null

    try {
      tournament.value = await tournamentService.getActive()
    } catch (err) {
      error.value = err.message || 'Failed to fetch tournament'
      console.error('Failed to fetch tournament:', err)
    } finally {
      isLoading.value = false
    }
  }

  onMounted(() => {
    fetchTournament()
  })

  return {
    tournament,
    isLoading,
    error,
    refetch: fetchTournament
  }
}
