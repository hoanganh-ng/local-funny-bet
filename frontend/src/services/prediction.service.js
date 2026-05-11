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
  }
}
