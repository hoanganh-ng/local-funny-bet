<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { matchService } from '../services/match.service.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import BaseCard from '../components/base/BaseCard.vue'
import MatchHeroCard from '../components/match/MatchHeroCard.vue'
import CrowdPicks from '../components/match/CrowdPicks.vue'
import PredictionForm from '../components/prediction/PredictionForm.vue'
import PredictionHistory from '../components/prediction/PredictionHistory.vue'

const route = useRoute()
const matchId = computed(() => route.params.id)
const leaderboardId = computed(() => route.query.leaderboard)

const match = ref(null)
const predictions = ref([])
const isLoadingMatch = ref(false)
const isLoadingPredictions = ref(false)

const kickoffPassed = computed(() => {
  if (!match.value) return false
  return new Date(match.value.kickoffAt) < new Date()
})

const currentPrediction = computed(() => {
  return predictions.value.find(p => p.is_current_user) || null
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
</script>

<template>
  <div class="match-view">
    <div class="container">
      <div v-if="isLoadingMatch" class="loading-state">
        Loading match details...
      </div>
      <template v-else-if="match">
        <MatchHeroCard :match="match" />

        <section v-if="leaderboardId" class="crowd-section">
          <div v-if="isLoadingPredictions" class="loading-state">
            Loading predictions...
          </div>
          <CrowdPicks
            v-else
            :predictions="predictions"
            :home-team="match.homeTeam"
            :away-team="match.awayTeam"
          />
        </section>

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

.crowd-section {
  display: flex;
  flex-direction: column;
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
</style>
