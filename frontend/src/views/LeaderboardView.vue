<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useLeaderboard } from '../composables/useLeaderboard.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import { useAuthStore } from '../store/auth.store.js'
import { useLeaderboardStore } from '../store/leaderboard.store.js'
import LeaderboardHeader from '../components/leaderboard/LeaderboardHeader.vue'
import LeaderboardPodium from '../components/leaderboard/LeaderboardPodium.vue'
import LeaderboardTable from '../components/leaderboard/LeaderboardTable.vue'
import MatchList from '../components/match/MatchList.vue'
import AdSlot from '../components/base/AdSlot.vue'

const route = useRoute()
const authStore = useAuthStore()
const leaderboardStore = useLeaderboardStore()
const leaderboardId = computed(() => route.params.id)

const { scores, isLoading: isLoadingScores } = useLeaderboard(leaderboardId)

const leaderboard = ref(null)
const isOwner = computed(() => !!authStore.user && authStore.user.id === leaderboard.value?.created_by)
const isLoadingLeaderboard = ref(false)

onMounted(async () => {
  await Promise.all([
    loadLeaderboard(),
    leaderboardStore.loadMatches()
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

async function handleGetInviteLink() {
  const response = await leaderboardService.generateInviteLink(leaderboardId.value)
  const url = `${window.location.origin}/join?token=${response.token}`
  await navigator.clipboard.writeText(url)
}
</script>

<template>
  <div class="leaderboard-view">
    <div class="container">
      <div v-if="isLoadingLeaderboard" class="loading-state">
        Loading leaderboard...
      </div>
      <template v-else-if="leaderboard">
        <LeaderboardHeader
          :leaderboard="leaderboard"
          :is-owner="isOwner"
          :on-get-invite-link="handleGetInviteLink"
        />

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
          <h2 class="section-title">Matches</h2>
          <MatchList
            :matches="leaderboardStore.matches"
            :has-more="leaderboardStore.hasMore"
            :loading-more="leaderboardStore.loadingMore"
            @load-more="leaderboardStore.loadMoreMatches()"
          />
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

</style>
