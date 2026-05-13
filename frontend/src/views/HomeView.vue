<script setup>
import { ref, computed, onMounted } from 'vue'
import { matchService } from '../services/match.service.js'
import MatchCard from '../components/common/MatchCard.vue'
import AdSlot from '../components/base/AdSlot.vue'

const matches = ref([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    const [scheduled, live] = await Promise.all([
      matchService.list('scheduled'),
      matchService.list('live')
    ])
    matches.value = [...live, ...scheduled]
  } catch (error) {
    console.error('Failed to load matches:', error)
  } finally {
    isLoading.value = false
  }
})

const now = new Date()
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
const tomorrow = new Date(today)
tomorrow.setDate(tomorrow.getDate() + 1)

const liveMatches = computed(() => {
  return matches.value.filter(m => m.status === 'live')
})

const todayMatches = computed(() => {
  return matches.value.filter(m => {
    if (m.status === 'live') return false
    const kickoff = new Date(m.kickoffAt)
    return kickoff >= today && kickoff < tomorrow
  })
})

const tomorrowMatches = computed(() => {
  return matches.value.filter(m => {
    const kickoff = new Date(m.kickoffAt)
    const dayAfter = new Date(tomorrow)
    dayAfter.setDate(dayAfter.getDate() + 1)
    return kickoff >= tomorrow && kickoff < dayAfter
  })
})

const laterMatches = computed(() => {
  return matches.value.filter(m => {
    const kickoff = new Date(m.kickoffAt)
    const dayAfter = new Date(tomorrow)
    dayAfter.setDate(dayAfter.getDate() + 1)
    return kickoff >= dayAfter
  })
})

const statsToday = computed(() => {
  const all = [...liveMatches.value, ...todayMatches.value]
  const predicted = all.filter(m => m.userPrediction)
  return {
    live: liveMatches.value.length,
    upcoming: todayMatches.value.length,
    predicted: predicted.length,
    total: all.length
  }
})
</script>

<template>
  <div class="home-view">
    <div class="container">
      <AdSlot position="top" />

      <!-- Header -->
      <header class="page-header">
        <div class="header-content">
          <p class="header-time">Today · {{ new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }) }}</p>
          <h1 class="page-title">Today's slate</h1>
          <p class="page-subtitle">
            {{ statsToday.live }} live · {{ statsToday.upcoming }} in the next 24h · {{ statsToday.predicted }}/{{ statsToday.total }} predicted
          </p>
        </div>
      </header>

      <div v-if="isLoading" class="loading-state">
        Loading matches...
      </div>

      <div v-else class="matches-content">
        <!-- LIVE NOW -->
        <section v-if="liveMatches.length > 0" class="matches-section">
          <div class="section-header">
            <h2 class="section-title">
              <span class="live-indicator">● LIVE</span>
              NOW
            </h2>
          </div>
          <div class="matches-list">
            <MatchCard
              v-for="match in liveMatches"
              :key="match.id"
              :match="match"
            />
          </div>
        </section>

        <!-- TODAY -->
        <section v-if="todayMatches.length > 0" class="matches-section">
          <div class="section-header">
            <h2 class="section-title">Later Today</h2>
          </div>
          <div class="matches-list">
            <MatchCard
              v-for="match in todayMatches"
              :key="match.id"
              :match="match"
            />
          </div>
        </section>

        <!-- TOMORROW -->
        <section v-if="tomorrowMatches.length > 0" class="matches-section">
          <div class="section-header">
            <h2 class="section-title">Tomorrow</h2>
          </div>
          <div class="matches-list">
            <MatchCard
              v-for="match in tomorrowMatches"
              :key="match.id"
              :match="match"
            />
          </div>
        </section>

        <!-- UPCOMING -->
        <section v-if="laterMatches.length > 0" class="matches-section">
          <div class="section-header">
            <h2 class="section-title">Upcoming</h2>
          </div>
          <div class="matches-list">
            <MatchCard
              v-for="match in laterMatches.slice(0, 10)"
              :key="match.id"
              :match="match"
            />
          </div>
        </section>

        <!-- Empty state -->
        <div v-if="matches.length === 0" class="empty-state">
          <div class="empty-icon">⚽</div>
          <h2 class="empty-title">No matches scheduled</h2>
          <p class="empty-subtitle">
            Check back later for upcoming fixtures.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  min-height: 100vh;
  padding: var(--space-8) var(--space-6);
}

.container {
  max-width: 900px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-10);
}

.page-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.header-time {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: var(--leading-tight);
}

.page-subtitle {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

.loading-state {
  padding: var(--space-12);
  text-align: center;
  color: var(--color-text-secondary);
}

.matches-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-10);
}

.matches-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.live-indicator {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-accent);
  background: var(--color-pick-bg);
  border: var(--border-accent);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
  animation: pulse 2s ease-in-out infinite;
}

[data-theme="dark"] .live-indicator {
  border-radius: var(--radius-pill);
  box-shadow: var(--shadow-accent);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

.matches-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.empty-state {
  padding: var(--space-16) var(--space-8);
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
}

.empty-icon {
  font-size: 64px;
  opacity: 0.3;
}

.empty-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.empty-subtitle {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
}

@media (max-width: 768px) {
  .home-view {
    padding: var(--space-6) var(--space-4);
  }

  .page-title {
    font-size: var(--text-4xl);
  }
}
</style>
