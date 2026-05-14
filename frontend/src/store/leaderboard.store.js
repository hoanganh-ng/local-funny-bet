import { defineStore } from 'pinia'
import { ref } from 'vue'
import { leaderboardService } from '../services/leaderboard.service.js'

export const useLeaderboardStore = defineStore('leaderboard', () => {
  const leaderboards = ref([])
  const current = ref(null)
  const scoresByBoard = ref({})

  async function fetchAll() {
    const data = await leaderboardService.list()
    leaderboards.value = data
  }

  async function fetchOne(id) {
    const data = await leaderboardService.getOne(id)
    current.value = data
  }

  async function fetchBoardScores(id) {
    const data = await leaderboardService.getScores(id)
    scoresByBoard.value = { ...scoresByBoard.value, [id]: data.scores || [] }
  }

  function boardSummary(boardId, userId) {
    const scores = scoresByBoard.value[boardId] || []
    const leader = scores[0] || null
    const userIdx = userId ? scores.findIndex(s => s.user_id === userId) : -1
    const userScore = userIdx >= 0 ? scores[userIdx] : null
    return {
      memberCount: scores.length,
      leaderName: leader?.name ?? null,
      leaderPoints: leader?.points ?? null,
      userRank: userIdx >= 0 ? userIdx + 1 : null,
      userPoints: userScore?.points ?? null,
      pointsFromLead: userScore != null && leader != null
        ? userScore.points - leader.points
        : null,
      topMembers: scores.slice(0, 3),
    }
  }

  function setCurrent(leaderboard) {
    current.value = leaderboard
  }

  return {
    leaderboards,
    current,
    scoresByBoard,
    fetchAll,
    fetchOne,
    fetchBoardScores,
    boardSummary,
    setCurrent,
  }
})
