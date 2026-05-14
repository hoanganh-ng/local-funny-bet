<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useLeaderboard } from '../composables/useLeaderboard.js'
import { matchService } from '../services/match.service.js'
import { predictionService } from '../services/prediction.service.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import LeaderboardHeader from '../components/leaderboard/LeaderboardHeader.vue'
import LeaderboardPodium from '../components/leaderboard/LeaderboardPodium.vue'
import LeaderboardTable from '../components/leaderboard/LeaderboardTable.vue'
import PredictionForm from '../components/prediction/PredictionForm.vue'
import AdSlot from '../components/base/AdSlot.vue'

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
  return matchPredictions.find(p => p.is_current_user) || null
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
          <template v-else>
            <LeaderboardPodium :scores="scores" />
            <div class="table-wrapper">
              <h3 class="table-subtitle">Full Standings</h3>
              <LeaderboardTable :scores="scores" />
            </div>
          </template>
        </section>

        <AdSlot position="between-matches" />

        <section class="matches-section">
          <h2 class="section-title">Matches & Predictions</h2>
          <div v-if="isLoadingMatches" class="loading-state">
            Loading matches...
          </div>
          <div v-else class="matches-list">
            <PredictionForm
              v-for="match in matches"
              :key="match.id"
              :match="match"
              :current-prediction="getCurrentPrediction(match.id)"
              @success="onPredictionSuccess(match.id)"
            />
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

.table-wrapper {
  margin-top: var(--space-8);
}

.table-subtitle {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
  margin-bottom: var(--space-4);
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

</style>
