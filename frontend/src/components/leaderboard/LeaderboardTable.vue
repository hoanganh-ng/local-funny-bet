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
    <div v-else class="scores-list">
      <ScoreRow
        v-for="(score, index) in scores"
        :key="score.user.id"
        :rank="index + 1"
        :user="score.user"
        :points="score.points"
        :is-current-user="authStore.user?.id === score.user.id"
      />
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

.scores-list {
  display: flex;
  flex-direction: column;
}
</style>
