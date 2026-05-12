<script setup>
import ScoreRow from './ScoreRow.vue'
import { useAuthStore } from '../../store/auth.store.js'

const props = defineProps({
  scores: {
    type: Array,
    default: () => []
  }
})

const authStore = useAuthStore()
</script>

<template>
  <div class="leaderboard-table">
    <div v-if="scores.length === 0" class="empty-state">
      No scores yet — predictions lock after kickoff
    </div>
    <div v-else class="table-container">
      <div class="table-header">
        <div class="header-cell header-rank">#</div>
        <div class="header-cell header-avatar"></div>
        <div class="header-cell header-name">Player</div>
        <div class="header-cell header-points">Points</div>
        <div class="header-cell header-streak">Streak</div>
        <div class="header-cell header-trend">Trend</div>
      </div>
      <div class="scores-list">
        <ScoreRow
          v-for="(score, index) in scores"
          :key="score.user.id"
          :rank="index + 1"
          :user="score.user"
          :points="score.points"
          :streak="score.streak || 0"
          :trend="score.trend || 'same'"
          :is-current-user="authStore.user?.id === score.user.id"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard-table {
  display: flex;
  flex-direction: column;
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.empty-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.table-container {
  display: flex;
  flex-direction: column;
}

.table-header {
  display: grid;
  grid-template-columns: auto auto 1fr auto auto auto;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  background: var(--color-bg-elevated);
  border-bottom: var(--border-hairline);
  position: sticky;
  top: 0;
  z-index: 1;
}

.header-cell {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.header-rank {
  min-width: var(--space-8);
  text-align: center;
}

.header-avatar {
  width: 32px;
}

.header-points,
.header-streak,
.header-trend {
  text-align: center;
}

.header-streak {
  min-width: var(--space-12);
}

.header-trend {
  min-width: var(--space-8);
}

.scores-list {
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .table-header {
    grid-template-columns: auto auto 1fr auto;
    gap: var(--space-3);
  }

  .header-streak,
  .header-trend {
    display: none;
  }
}
</style>
