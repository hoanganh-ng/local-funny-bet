import { request } from './api.js'

export const predictionService = {
  async upsert(matchId, value) {
    return request('/predictions', {
      method: 'PUT',
      body: JSON.stringify({
        match_id: matchId,
        value
      })
    })
  },

  async getHistory() {
    const data = await request('/predictions/history')
    return data.map(h => ({
      ...h,
      isCorrect: h.is_correct,
      predictionLabel: h.prediction_label,
      homeTeam: h.home_team,
      awayTeam: h.away_team,
      homeScore: h.home_score,
      awayScore: h.away_score,
      kickoffAt: h.kickoff_at,
    }))
  }
}
