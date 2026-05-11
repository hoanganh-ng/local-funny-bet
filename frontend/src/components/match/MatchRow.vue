<script setup>
import MatchStatusBadge from './MatchStatusBadge.vue'

const props = defineProps({
  match: {
    type: Object,
    required: true
  },
  userPrediction: {
    type: Object,
    default: null
  }
})

const predictionText = (prediction) => {
  if (!prediction) return '-'
  if (prediction.outcome === 'home_win') return 'Home'
  if (prediction.outcome === 'draw') return 'Draw'
  if (prediction.outcome === 'away_win') return 'Away'
  return '-'
}
</script>

<template>
  <div class="match-row">
    <div class="team-cell">{{ match.homeTeam }}</div>
    <div class="score-cell">
      <span v-if="match.status !== 'scheduled'">
        {{ match.homeScore ?? '-' }} - {{ match.awayScore ?? '-' }}
      </span>
      <span v-else class="scheduled-dash">-</span>
    </div>
    <div class="team-cell">{{ match.awayTeam }}</div>
    <div class="status-cell">
      <MatchStatusBadge
        :status="match.status"
        :kickoff-at="match.kickoffAt"
      />
    </div>
    <div v-if="userPrediction !== undefined" class="prediction-cell">
      {{ predictionText(userPrediction) }}
    </div>
  </div>
</template>

<style scoped>
.match-row {
  display: grid;
  grid-template-columns: 1fr auto 1fr auto auto;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-bottom: var(--border-hairline);
}

.match-row:last-child {
  border-bottom: none;
}

.team-cell {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.team-cell:first-child {
  text-align: right;
}

.team-cell:nth-child(3) {
  text-align: left;
}

.score-cell {
  font-size: var(--text-base);
  font-weight: var(--font-bold);
  font-family: var(--font-mono);
  color: var(--color-text-primary);
  min-width: var(--space-12);
  text-align: center;
}

.scheduled-dash {
  color: var(--color-text-secondary);
}

.status-cell {
  display: flex;
  justify-content: flex-end;
  min-width: 120px;
}

.prediction-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  min-width: var(--space-12);
  text-align: center;
}
</style>
