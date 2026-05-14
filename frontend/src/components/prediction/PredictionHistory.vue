<script setup>
import BaseAvatar from '../base/BaseAvatar.vue'
import BaseBadge from '../base/BaseBadge.vue'

const props = defineProps({
  predictions: {
    type: Array,
    default: () => []
  }
})

const outcomeLabel = (outcome) => {
  if (outcome === 'home_win') return 'Home'
  if (outcome === 'draw') return 'Draw'
  if (outcome === 'away_win') return 'Away'
  return '-'
}

const outcomeClass = (outcome) => {
  return `outcome-${outcome}`
}

const avatarUrl = (p) => p.user_avatar_url || null
</script>

<template>
  <div class="prediction-history">
    <h3 class="history-title">Predictions</h3>
    <div v-if="predictions.length === 0" class="empty-state">
      No predictions yet
    </div>
    <div v-else class="predictions-list">
      <div
        v-for="prediction in predictions"
        :key="prediction.user_id"
        class="prediction-item"
      >
        <BaseAvatar
          :src="avatarUrl(prediction)"
          :name="prediction.user_name"
          size="sm"
        />
        <span class="user-name">{{ prediction.user_name }}</span>
        <BaseBadge :class="outcomeClass(prediction.outcome)">
          {{ outcomeLabel(prediction.outcome) }}
        </BaseBadge>
      </div>
    </div>
  </div>
</template>

<style scoped>
.prediction-history {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.history-title {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.empty-state {
  padding: var(--space-6);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.predictions-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.prediction-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface);
  border: var(--border-hairline);
  border-radius: var(--radius-md);
}

.user-name {
  flex: 1;
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.outcome-home_win {
  background: var(--color-home-win);
  color: var(--color-bg);
}

.outcome-draw {
  background: var(--color-draw);
  color: var(--color-bg);
}

.outcome-away_win {
  background: var(--color-away-win);
  color: var(--color-bg);
}

[data-theme="light"] .outcome-home_win,
[data-theme="light"] .outcome-draw,
[data-theme="light"] .outcome-away_win {
  background: transparent;
  color: var(--color-text-primary);
  border-color: var(--color-border-strong);
}
</style>
