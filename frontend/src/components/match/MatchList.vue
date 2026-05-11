<script setup>
import MatchCard from '../common/MatchCard.vue'

const props = defineProps({
  matches: {
    type: Array,
    default: () => []
  }
})
</script>

<template>
  <div class="match-list">
    <div v-if="matches.length === 0" class="empty-state">
      No matches yet
    </div>
    <MatchCard
      v-for="match in matches"
      :key="match.id"
      :match="match"
    />
  </div>
</template>

<style scoped>
.match-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.empty-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}
</style>
