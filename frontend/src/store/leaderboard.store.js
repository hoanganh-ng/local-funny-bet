import { defineStore } from 'pinia'
import { ref } from 'vue'
import { leaderboardService } from '../services/leaderboard.service.js'

export const useLeaderboardStore = defineStore('leaderboard', () => {
  const leaderboards = ref([])
  const current = ref(null)

  async function fetchAll() {
    const data = await leaderboardService.list()
    leaderboards.value = data
  }

  async function fetchOne(id) {
    const data = await leaderboardService.getOne(id)
    current.value = data
  }

  function setCurrent(leaderboard) {
    current.value = leaderboard
  }

  return {
    leaderboards,
    current,
    fetchAll,
    fetchOne,
    setCurrent
  }
})
