<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { matchService } from '../services/match.service.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import BaseCard from '../components/base/BaseCard.vue'
import MatchStatusBadge from '../components/match/MatchStatusBadge.vue'
import PredictionForm from '../components/prediction/PredictionForm.vue'
import PredictionHistory from '../components/prediction/PredictionHistory.vue'

const route = useRoute()
const matchId = computed(() => route.params.id)
const leaderboardId = computed(() => route.query.leaderboard)

const match = ref(null)
const predictions = ref([])
const isLoadingMatch = ref(false)
const isLoadingPredictions = ref(false)

const showScore = computed(() => {
  return match.value && (match.value.status === 'finished' || match.value.status === 'live')
})

const kickoffPassed = computed(() => {
  if (!match.value) return false
  return new Date(match.value.kickoffAt) < new Date()
})

const currentPrediction = computed(() => {
  return predictions.value.find(p => p.isCurrentUser) || null
})

onMounted(async () => {
  await loadMatch()
  if (leaderboardId.value) {
    await loadPredictions()
  }
})

async function loadMatch() {
  isLoadingMatch.value = true
  try {
    match.value = await matchService.getOne(matchId.value)
  } catch (err) {
    console.error('Failed to load match:', err)
  } finally {
    isLoadingMatch.value = false
  }
}

async function loadPredictions() {
  if (!leaderboardId.value) return

  isLoadingPredictions.value = true
  try {
    predictions.value = await leaderboardService.getMatchPredictions(
      leaderboardId.value,
      matchId.value
    )
  } catch (err) {
    console.error('Failed to load predictions:', err)
    predictions.value = []
  } finally {
    isLoadingPredictions.value = false
  }
}

async function onPredictionSuccess() {
  if (leaderboardId.value) {
    await loadPredictions()
  }
}

function formatKickoffTime(kickoffAt) {
  const date = new Date(kickoffAt)
  return date.toLocaleString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="match-view">
    <div class="container">
      <div v-if="isLoadingMatch" class="loading-state">
        Loading match details...
      </div>
      <template v-else-if="match">
        <BaseCard class="match-detail-card">
          <div class="match-header">
            <MatchStatusBadge
              :status="match.status"
              :kickoff-at="match.kickoffAt"
            />
          </div>

          <div class="match-teams">
            <div class="team-row">
              <span class="team-name">{{ match.homeTeam }}</span>
              <span v-if="showScore" class="score">{{ match.homeScore ?? '-' }}</span>
            </div>

            <div class="vs-divider">VS</div>

            <div class="team-row">
              <span class="team-name">{{ match.awayTeam }}</span>
              <span v-if="showScore" class="score">{{ match.awayScore ?? '-' }}</span>
            </div>
          </div>

          <div class="match-meta">
            <p class="kickoff-time">
              <strong>Kickoff:</strong> {{ formatKickoffTime(match.kickoffAt) }}
            </p>
          </div>
        </BaseCard>

        <section class="prediction-section">
          <h2 class="section-title">Your Prediction</h2>
          <BaseCard>
            <PredictionForm
              :match="match"
              :current-prediction="currentPrediction"
              @success="onPredictionSuccess"
            />
          </BaseCard>
        </section>

        <section v-if="kickoffPassed && leaderboardId" class="history-section">
          <h2 class="section-title">All Predictions</h2>
          <BaseCard>
            <div v-if="isLoadingPredictions" class="loading-state">
              Loading predictions...
            </div>
            <PredictionHistory v-else :predictions="predictions" />
          </BaseCard>
        </section>
      </template>
      <div v-else class="error-state">
        Failed to load match
      </div>
    </div>
  </div>
</template>

<style scoped>
.match-view {
  min-height: 100vh;
  background: var(--color-bg);
  padding: var(--space-6) var(--space-4);
}

.container {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.loading-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.error-state {
  padding: var(--space-12);
  text-align: center;
  color: var(--color-danger);
  font-size: var(--text-base);
}

.match-detail-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  padding: var(--space-8);
}

.match-header {
  display: flex;
  justify-content: flex-end;
}

.match-teams {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.team-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.team-name {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.score {
  font-size: var(--text-4xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  min-width: var(--space-12);
  text-align: center;
}

.vs-divider {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-align: center;
  text-transform: uppercase;
  letter-spacing: var(--tracking-widest);
  padding: var(--space-2) 0;
}

.match-meta {
  border-top: var(--border-hairline);
  padding-top: var(--space-4);
}

.kickoff-time {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

.kickoff-time strong {
  color: var(--color-text-primary);
  font-weight: var(--font-semibold);
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-4);
}

.prediction-section,
.history-section {
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .team-name {
    font-size: var(--text-xl);
  }

  .score {
    font-size: var(--text-3xl);
  }
}
</style>
