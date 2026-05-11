<script setup>
import { computed } from 'vue'
import BaseCard from '../base/BaseCard.vue'
import MatchStatusBadge from '../match/MatchStatusBadge.vue'

const props = defineProps({
  match: {
    type: Object,
    required: true
  }
})

const showScore = computed(() => {
  return props.match.status === 'finished' || props.match.status === 'live'
})
</script>

<template>
  <BaseCard class="match-card">
    <div class="match-header">
      <MatchStatusBadge
        :status="match.status"
        :kickoff-at="match.kickoffAt"
      />
    </div>

    <div class="match-teams">
      <div class="team">
        <span class="team-name">{{ match.homeTeam }}</span>
        <span v-if="showScore" class="score">{{ match.homeScore ?? '-' }}</span>
      </div>

      <div class="vs">vs</div>

      <div class="team">
        <span class="team-name">{{ match.awayTeam }}</span>
        <span v-if="showScore" class="score">{{ match.awayScore ?? '-' }}</span>
      </div>
    </div>
  </BaseCard>
</template>

<style scoped>
.match-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-4);
}

.match-header {
  display: flex;
  justify-content: flex-end;
}

.match-teams {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.team {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.team-name {
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
}

.score {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  min-width: var(--space-8);
  text-align: center;
}

.vs {
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-align: center;
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}
</style>
