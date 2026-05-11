<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useLeaderboard } from '../composables/useLeaderboard.js'
import { matchService } from '../services/match.service.js'
import { predictionService } from '../services/prediction.service.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import LeaderboardHeader from '../components/leaderboard/LeaderboardHeader.vue'
import LeaderboardTable from '../components/leaderboard/LeaderboardTable.vue'
import MatchCard from '../components/common/MatchCard.vue'
import PredictionForm from '../components/prediction/PredictionForm.vue'
import AdSlot from '../components/base/AdSlot.vue'
import BaseCard from '../components/base/BaseCard.vue'

const route = useRoute()
const leaderboardId = computed(() => route.params.id)

const { scores, isLoading: isLoadingScores } = useLeaderboard(leaderboardId)

const leaderboard = ref(null)
const matches = ref([])
const predictions = ref({})
const isLoadingLeaderboard = ref(false)
const isLoadingMatches = ref(false)

onMounted(async () => {
  await Promise.all([
    loadLeaderboard(),
    loadMatches()
  ])
})

async function loadLeaderboard() {
  isLoadingLeaderboard.value = true
  try {
    leaderboard.value = await leaderboardService.getOne(leaderboardId.value)
  } catch (err) {
    console.error('Failed to load leaderboard:', err)
  } finally {
    isLoadingLeaderboard.value = false
  }
}

async function loadMatches() {
  isLoadingMatches.value = true
  try {
    matches.value = await matchService.list()

    const predictionPromises = matches.value.map(async (match) => {
      try {
        const preds = await leaderboardService.getMatchPredictions(
          leaderboardId.value,
          match.id
        )
        predictions.value[match.id] = preds
      } catch (err) {
        predictions.value[match.id] = []
      }
    })

    await Promise.all(predictionPromises)
  } catch (err) {
    console.error('Failed to load matches:', err)
  } finally {
    isLoadingMatches.value = false
  }
}

function getCurrentPrediction(matchId) {
  const matchPredictions = predictions.value[matchId] || []
  return matchPredictions.find(p => p.isCurrentUser) || null
}

async function onPredictionSuccess(matchId) {
  try {
    const preds = await leaderboardService.getMatchPredictions(
      leaderboardId.value,
      matchId
    )
    predictions.value[matchId] = preds
  } catch (err) {
    console.error('Failed to reload predictions:', err)
  }
}
</script>

<template>
  <div class="leaderboard-view">
    <div class="container">
      <div v-if="isLoadingLeaderboard" class="loading-state">
        Loading leaderboard...
      </div>
      <template v-else-if="leaderboard">
        <LeaderboardHeader :leaderboard="leaderboard" />

        <section class="scores-section">
          <h2 class="section-title">Rankings</h2>
          <div v-if="isLoadingScores" class="loading-state">
            Loading scores...
          </div>
          <LeaderboardTable v-else :scores="scores" />
        </section>

        <AdSlot position="between-matches" />

        <section class="matches-section">
          <h2 class="section-title">Matches & Predictions</h2>
          <div v-if="isLoadingMatches" class="loading-state">
            Loading matches...
          </div>
          <div v-else class="matches-list">
            <BaseCard
              v-for="match in matches"
              :key="match.id"
              class="match-prediction-card"
            >
              <div class="match-section">
                <MatchCard :match="match" />
              </div>
              <div class="prediction-section">
                <PredictionForm
                  :match="match"
                  :current-prediction="getCurrentPrediction(match.id)"
                  @success="onPredictionSuccess(match.id)"
                />
              </div>
            </BaseCard>
          </div>
        </section>
      </template>
      <div v-else class="error-state">
        Failed to load leaderboard
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard-view {
  min-height: 100vh;
  background: var(--color-bg);
  padding: var(--space-6) var(--space-4);
}

.container {
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-12);
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-6);
}

.scores-section,
.matches-section {
  display: flex;
  flex-direction: column;
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

.matches-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.match-prediction-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  padding: var(--space-6);
}

.match-section {
  flex: 1;
}

.prediction-section {
  border-top: var(--border-hairline);
  padding-top: var(--space-6);
}

@media (min-width: 768px) {
  .match-prediction-card {
    flex-direction: row;
  }

  .prediction-section {
    border-top: none;
    border-left: var(--border-hairline);
    padding-top: 0;
    padding-left: var(--space-6);
    min-width: 280px;
  }
}
</style>
