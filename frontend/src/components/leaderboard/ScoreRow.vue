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
  </div>
</template>

<style scoped>
.score-row {
  display: grid;
  grid-template-columns: auto auto 1fr auto;
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
</style>
