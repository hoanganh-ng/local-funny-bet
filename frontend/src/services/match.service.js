import { request } from './api.js'

function transformMatch(match) {
  return {
    id: match.id,
    tournamentId: match.tournament_id,
    homeTeam: match.home_team,
    awayTeam: match.away_team,
    homeTeamCode: match.home_team_code,
    awayTeamCode: match.away_team_code,
    group: match.group,
    homeScore: match.home_score,
    awayScore: match.away_score,
    kickoffAt: match.kickoff_at,
    status: match.status,
    result: match.result,
    points: match.points,
    userPrediction: match.user_prediction,
    externalId: match.external_id
  }
}

export const matchService = {
  async list(status) {
    const query = status ? `?status=${status}` : ''
    const matches = await request(`/matches${query}`)
    return matches.map(transformMatch)
  },

  async getOne(id) {
    const match = await request(`/matches/${id}`)
    return transformMatch(match)
  }
}
