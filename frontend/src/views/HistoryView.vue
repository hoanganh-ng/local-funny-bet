<script setup>
import { ref, onMounted, computed } from 'vue'
import { predictionService } from '../services/prediction.service.js'

const history = ref([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    history.value = await predictionService.getHistory()
  } catch (error) {
    console.error('Failed to load history:', error)
  } finally {
    isLoading.value = false
  }
})

const stats = computed(() => {
  if (!history.value.length) {
    return {
      totalPoints: 0,
      accuracy: 0,
      bestStreak: 0,
      currentStreak: 0,
      totalPredictions: 0
    }
  }

  const settled = history.value.filter(h => h.result)
  const correct = settled.filter(h => h.isCorrect)

  let bestStreak = 0
  let currentStreak = 0
  let tempStreak = 0

  settled.forEach((h, idx) => {
    if (h.isCorrect) {
      tempStreak++
      if (idx === settled.length - 1) {
        currentStreak = tempStreak
      }
    } else {
      bestStreak = Math.max(bestStreak, tempStreak)
      tempStreak = 0
      currentStreak = 0
    }
  })
  bestStreak = Math.max(bestStreak, tempStreak)

  return {
    totalPoints: history.value.reduce((sum, h) => sum + (h.points || 0), 0),
    accuracy: settled.length > 0 ? Math.round((correct.length / settled.length) * 100) : 0,
    bestStreak,
    currentStreak,
    totalPredictions: settled.length
  }
})

const recentPicks = computed(() => {
  return history.value
    .filter(h => h.result)
    .slice(0, 10)
})

const groupedByMatchday = computed(() => {
  const grouped = {}
  history.value.forEach(h => {
    const day = h.matchday || 'Matchday 1'
    if (!grouped[day]) grouped[day] = []
    grouped[day].push(h)
  })
  return grouped
})
</script>

<template>
  <div class="history-view">
    <div class="container">
      <header class="page-header">
        <div>
          <h1 class="page-title">Your history</h1>
          <p class="page-subtitle">
            Every prediction you've made this tournament, with the score and the points it earned you.
          </p>
        </div>
      </header>

      <div v-if="isLoading" class="loading-state">
        Loading history...
      </div>

      <div v-else class="history-content">
        <!-- Stats cards -->
        <div class="stats-grid">
          <div class="stat-card">
            <p class="stat-label">Total Points</p>
            <p class="stat-value">{{ stats.totalPoints }}</p>
            <p class="stat-meta">across all boards</p>
          </div>

          <div class="stat-card">
            <p class="stat-label">Accuracy</p>
            <p class="stat-value">{{ stats.accuracy }}%</p>
            <p class="stat-meta">{{ stats.totalPredictions }} of {{ history.length }} correct</p>
          </div>

          <div class="stat-card">
            <p class="stat-label">Best Streak</p>
            <p class="stat-value">{{ stats.bestStreak }}</p>
            <p class="stat-meta">matchdays 1-{{ stats.bestStreak }}</p>
          </div>

          <div class="stat-card stat-card--highlight">
            <p class="stat-label">Current Streak</p>
            <p class="stat-value">{{ stats.currentStreak }}</p>
            <p class="stat-meta">in a row</p>
          </div>
        </div>

        <!-- Recent picks sparkline -->
        <div class="recent-section">
          <h3 class="section-title">Recent Picks</h3>
          <div class="recent-picks">
            <div
              v-for="(pick, idx) in recentPicks"
              :key="idx"
              class="recent-pick"
              :class="{ correct: pick.isCorrect, wrong: !pick.isCorrect }"
            >
              <span v-if="pick.isCorrect" class="pick-icon">✓</span>
              <span v-else class="pick-icon">✗</span>
            </div>
          </div>
        </div>

        <!-- Matchday log -->
        <div class="matchday-section">
          <h3 class="section-title">Matchday by Matchday</h3>
          <div
            v-for="(matches, matchday) in groupedByMatchday"
            :key="matchday"
            class="matchday-group"
          >
            <h4 class="matchday-title">{{ matchday }}</h4>
            <div class="matches-list">
              <div
                v-for="match in matches"
                :key="match.id"
                class="match-row"
              >
                <div class="match-info">
                  <span class="match-time">{{ new Date(match.kickoffAt).toLocaleDateString() }}</span>
                  <span class="match-teams">
                    {{ match.homeTeam }} {{ match.homeScore ?? '-' }} - {{ match.awayScore ?? '-' }} {{ match.awayTeam }}
                  </span>
                </div>
                <div class="match-prediction">
                  <span class="prediction-label">{{ match.predictionLabel || 'No pick' }}</span>
                  <span v-if="match.points !== null" class="prediction-points" :class="{ correct: match.isCorrect }">
                    {{ match.isCorrect ? '+' : '' }}{{ match.points }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-view {
  min-height: 100vh;
  padding: var(--space-8) var(--space-6);
}

.container {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: var(--space-10);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--text-4xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-3);
}

.page-subtitle {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
  max-width: 600px;
}

.loading-state {
  padding: var(--space-12);
  text-align: center;
  color: var(--color-text-secondary);
}

.history-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-10);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--space-4);
}

.stat-card {
  padding: var(--space-6);
  background: var(--color-surface);
  border: var(--border-hairline);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

[data-theme="dark"] .stat-card {
  border-radius: var(--radius-lg);
}

.stat-card--highlight {
  background: var(--color-accent);
  border-color: var(--color-accent);
}

[data-theme="dark"] .stat-card--highlight {
  box-shadow: var(--shadow-accent);
}

.stat-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.stat-card--highlight .stat-label {
  color: var(--color-accent-text);
  opacity: 0.8;
}

.stat-value {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: 1;
}

.stat-card--highlight .stat-value {
  color: var(--color-accent-text);
}

.stat-meta {
  font-size: var(--text-sm);
  color: var(--color-text-disabled);
}

.stat-card--highlight .stat-meta {
  color: var(--color-accent-text);
  opacity: 0.7;
}

/* Recent Picks */
.recent-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.recent-picks {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.recent-pick {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  border: var(--border-hairline);
}

[data-theme="dark"] .recent-pick {
  border-radius: var(--radius-md);
}

.recent-pick.correct {
  background: var(--color-correct);
  border-color: var(--color-correct);
  color: var(--color-accent-text);
}

[data-theme="dark"] .recent-pick.correct {
  background: rgba(63,185,80,0.2);
  box-shadow: 0 0 20px rgba(63,185,80,0.3);
}

.recent-pick.wrong {
  background: var(--color-surface);
  border-color: var(--color-border);
  color: var(--color-text-disabled);
  opacity: 0.5;
}

/* Matchday Log */
.matchday-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.matchday-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.matchday-title {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.matches-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.match-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4);
  background: var(--color-surface);
  border: var(--border-hairline);
  gap: var(--space-4);
}

[data-theme="dark"] .match-row {
  border-radius: var(--radius-md);
}

.match-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
}

.match-time {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
}

.match-teams {
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.match-prediction {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.prediction-label {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.prediction-points {
  font-family: var(--font-mono);
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-disabled);
  min-width: 32px;
  text-align: right;
}

.prediction-points.correct {
  color: var(--color-correct);
}

@media (max-width: 768px) {
  .history-view {
    padding: var(--space-6) var(--space-4);
  }

  .page-title {
    font-size: var(--text-3xl);
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-value {
    font-size: var(--text-4xl);
  }

  .match-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .match-prediction {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
