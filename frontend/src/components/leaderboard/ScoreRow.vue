<script setup>
import BaseAvatar from '../base/BaseAvatar.vue'

const props = defineProps({
  rank: {
    type: Number,
    required: true
  },
  user: {
    type: Object,
    required: true
  },
  points: {
    type: Number,
    required: true
  },
  streak: {
    type: Number,
    default: 0
  },
  trend: {
    type: String,
    default: 'same',
    validator: v => ['up', 'down', 'same'].includes(v)
  },
  isCurrentUser: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div class="score-row" :class="{ current: isCurrentUser }">
    <div class="rank">{{ rank }}</div>
    <BaseAvatar
      :src="user.avatarUrl"
      :name="user.name"
      size="sm"
    />
    <div class="user-name">{{ user.name }}</div>
    <div class="points">{{ points }}</div>
    <div class="streak">
      <span v-if="streak > 0" class="streak-value">{{ streak }}🔥</span>
      <span v-else class="streak-empty">—</span>
    </div>
    <div class="trend" :class="`trend--${trend}`">
      <span v-if="trend === 'up'" class="trend-icon">↑</span>
      <span v-else-if="trend === 'down'" class="trend-icon">↓</span>
      <span v-else class="trend-icon">—</span>
    </div>
  </div>
</template>

<style scoped>
.score-row {
  display: grid;
  grid-template-columns: auto auto 1fr auto auto auto;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-left: 2px solid transparent;
  transition: border-color var(--duration-fast) var(--ease-default);
}

.score-row.current {
  border-left-color: var(--color-accent);
  background: var(--color-surface);
}

[data-theme="light"] .score-row.current {
  background: transparent;
  border-left-color: var(--color-border-strong);
  border-left-width: 4px;
}

.rank {
  font-size: var(--text-base);
  font-weight: var(--font-bold);
  font-family: var(--font-mono);
  color: var(--color-text-secondary);
  min-width: var(--space-8);
  text-align: center;
}

.user-name {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.points {
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  font-family: var(--font-mono);
  color: var(--color-text-primary);
}

.streak {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  min-width: var(--space-12);
  text-align: center;
}

.streak-value {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-accent);
}

.streak-empty {
  color: var(--color-text-disabled);
}

.trend {
  font-size: var(--text-lg);
  min-width: var(--space-8);
  text-align: center;
}

.trend-icon {
  display: inline-block;
}

.trend--up .trend-icon {
  color: var(--color-correct);
}

.trend--down .trend-icon {
  color: var(--color-danger);
}

.trend--same .trend-icon {
  color: var(--color-text-disabled);
}

@media (max-width: 768px) {
  .score-row {
    grid-template-columns: auto auto 1fr auto;
    gap: var(--space-3);
  }

  .streak,
  .trend {
    display: none;
  }
}
</style>
