import { request } from './api.js'

export const matchService = {
  async list(status) {
    const query = status ? `?status=${status}` : ''
    return request(`/matches${query}`)
  },

  async getOne(id) {
    return request(`/matches/${id}`)
  }
}
