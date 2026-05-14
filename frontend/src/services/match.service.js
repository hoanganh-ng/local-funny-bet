import { request } from './api.js'

function transformMatch(match) {
  return {
    id: match.id,
    tournamentId: match.tournament_id,
    homeTeam: match.home_team,
    awayTeam: match.away_team,
    homeTeamCode: match.home_team_code,
    awayTeamCode: match.away_team_code,
    homeScore: match.home_score,
    awayScore: match.away_score,
    kickoffAt: match.kickoff_at,
    status: match.status,
    externalId: match.external_id
  }
}

export const matchService = {
  async list({ status = '', limit = 20, after = null } = {}) {
    const params = new URLSearchParams()
    if (status) params.set('status', status)
    params.set('limit', String(limit))
    if (after) params.set('after', after)
    const query = params.toString() ? `?${params.toString()}` : ''
    const data = await request(`/matches${query}`)
    return {
      matches: data.matches.map(transformMatch),
      nextCursor: data.next_cursor || null,
      hasMore: data.has_more || false,
    }
  },

  async getOne(id) {
    const match = await request(`/matches/${id}`)
    return transformMatch(match)
  }
}
