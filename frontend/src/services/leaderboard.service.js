import { request } from './api.js'

export const leaderboardService = {
  async create(name) {
    return request('/leaderboards', {
      method: 'POST',
      body: JSON.stringify({ name })
    })
  },

  async list() {
    return request('/leaderboards')
  },

  async getOne(id) {
    return request(`/leaderboards/${id}`)
  },

  async generateInviteLink(id) {
    return request(`/leaderboards/${id}/invite`, {
      method: 'POST'
    })
  },

  async join(inviteToken) {
    return request('/leaderboards/join', {
      method: 'POST',
      body: JSON.stringify({ invite_token: inviteToken })
    })
  },

  async getScores(id) {
    return request(`/leaderboards/${id}/scores`)
  },

  async getMatchPredictions(leaderboardId, matchId) {
    return request(`/leaderboards/${leaderboardId}/matches/${matchId}/predictions`)
  }
}
