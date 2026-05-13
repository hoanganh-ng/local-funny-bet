<script setup>
import { computed } from 'vue'
import MatchStatusBadge from './MatchStatusBadge.vue'

const props = defineProps({
  match: {
    type: Object,
    required: true
  }
})

const showScore = computed(() => {
  return props.match.status === 'finished' || props.match.status === 'live'
})

const kickoffDisplay = computed(() => {
  const date = new Date(props.match.kickoffAt)
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})
</script>

<template>
  <div class="match-hero-card">
    <div class="hero-header">
      <MatchStatusBadge
        :status="match.status"
        :kickoff-at="match.kickoffAt"
      />
      <div class="kickoff-time">{{ kickoffDisplay }}</div>
    </div>

    <div class="hero-body">
      <!-- Home Team -->
      <div class="team-section home">
        <div class="team-badge">
          <span class="team-flag">{{ match.homeTeam.slice(0, 3).toUpperCase() }}</span>
        </div>
        <h2 class="team-name">{{ match.homeTeam }}</h2>
        <div v-if="showScore" class="team-score">{{ match.homeScore ?? '-' }}</div>
      </div>

      <!-- Divider -->
      <div class="hero-divider">
        <span class="vs-text">VS</span>
      </div>

      <!-- Away Team -->
      <div class="team-section away">
        <div class="team-badge">
          <span class="team-flag">{{ match.awayTeam.slice(0, 3).toUpperCase() }}</span>
        </div>
        <h2 class="team-name">{{ match.awayTeam }}</h2>
        <div v-if="showScore" class="team-score">{{ match.awayScore ?? '-' }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.match-hero-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
  padding: var(--space-10) var(--space-8);
  background: var(--color-surface);
  border: var(--border-hairline);
  position: relative;
}

[data-theme="dark"] .match-hero-card {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.hero-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.kickoff-time {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.hero-body {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: var(--space-8);
  align-items: center;
}

.team-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  text-align: center;
}

.team-badge {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
}

[data-theme="dark"] .team-badge {
  border-radius: var(--radius-lg);
}

[data-theme="light"] .team-badge {
  border-width: 2px;
}

.team-flag {
  font-family: var(--font-mono);
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.team-name {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: var(--leading-tight);
}

.team-score {
  font-family: var(--font-display);
  font-size: var(--text-6xl);
  font-weight: var(--font-black);
  color: var(--color-accent);
  line-height: 1;
  margin-top: var(--space-2);
}

[data-theme="light"] .team-score {
  color: var(--color-text-primary);
}

.hero-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: var(--space-8) 0;
}

.vs-text {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  font-weight: var(--font-bold);
  color: var(--color-text-disabled);
  text-transform: uppercase;
  letter-spacing: var(--tracking-widest);
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
}

[data-theme="dark"] .vs-text {
  border-radius: var(--radius-sm);
}

.hero-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding-top: var(--space-4);
  border-top: var(--border-hairline);
}

.venue-icon {
  font-size: var(--text-base);
}

.venue-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

@media (max-width: 768px) {
  .match-hero-card {
    padding: var(--space-6);
  }

  .hero-body {
    grid-template-columns: 1fr;
    gap: var(--space-6);
  }

  .hero-divider {
    order: 2;
    padding: var(--space-4) 0;
  }

  .team-section.home {
    order: 1;
  }

  .team-section.away {
    order: 3;
  }

  .team-badge {
    width: 64px;
    height: 64px;
  }

  .team-flag {
    font-size: var(--text-base);
  }

  .team-name {
    font-size: var(--text-2xl);
  }

  .team-score {
    font-size: var(--text-5xl);
  }
}
</style>
