import { request } from './api.js'

function transformTournament(tournament) {
  return {
    id: tournament.id,
    name: tournament.name,
    season: tournament.season,
    logoUrl: tournament.logo_url,
    status: tournament.status,
    externalId: tournament.external_id
  }
}

export const tournamentService = {
  async getActive() {
    const tournament = await request('/tournaments/active')
    return transformTournament(tournament)
  }
}
