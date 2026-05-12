<script setup>
import { computed } from 'vue'

const props = defineProps({
  predictions: {
    type: Array,
    default: () => []
  },
  homeTeam: {
    type: String,
    default: 'Home'
  },
  awayTeam: {
    type: String,
    default: 'Away'
  }
})

const totalPredictions = computed(() => props.predictions.length)

const distribution = computed(() => {
  if (totalPredictions.value === 0) {
    return { home: 0, draw: 0, away: 0 }
  }

  const counts = {
    home: props.predictions.filter(p => p.prediction === 'home').length,
    draw: props.predictions.filter(p => p.prediction === 'draw').length,
    away: props.predictions.filter(p => p.prediction === 'away').length
  }

  return {
    home: Math.round((counts.home / totalPredictions.value) * 100),
    draw: Math.round((counts.draw / totalPredictions.value) * 100),
    away: Math.round((counts.away / totalPredictions.value) * 100)
  }
})

const maxPct = computed(() => {
  return Math.max(distribution.value.home, distribution.value.draw, distribution.value.away)
})
</script>

<template>
  <div class="crowd-picks">
    <div class="crowd-header">
      <h3 class="crowd-title">Crowd Picks</h3>
      <p class="crowd-subtitle">{{ totalPredictions }} prediction{{ totalPredictions !== 1 ? 's' : '' }}</p>
    </div>

    <div v-if="totalPredictions === 0" class="empty-state">
      No predictions yet. Be the first!
    </div>

    <div v-else class="picks-grid">
      <!-- Home -->
      <div class="pick-item" :class="{ leading: distribution.home === maxPct && maxPct > 0 }">
        <div class="pick-label">
          <span class="pick-outcome">1</span>
          <span class="pick-team">{{ homeTeam }}</span>
        </div>
        <div class="pick-bar-container">
          <div class="pick-bar" :style="{ width: `${distribution.home}%` }"></div>
        </div>
        <div class="pick-pct">{{ distribution.home }}%</div>
      </div>

      <!-- Draw -->
      <div class="pick-item" :class="{ leading: distribution.draw === maxPct && maxPct > 0 }">
        <div class="pick-label">
          <span class="pick-outcome">X</span>
          <span class="pick-team">Draw</span>
        </div>
        <div class="pick-bar-container">
          <div class="pick-bar" :style="{ width: `${distribution.draw}%` }"></div>
        </div>
        <div class="pick-pct">{{ distribution.draw }}%</div>
      </div>

      <!-- Away -->
      <div class="pick-item" :class="{ leading: distribution.away === maxPct && maxPct > 0 }">
        <div class="pick-label">
          <span class="pick-outcome">2</span>
          <span class="pick-team">{{ awayTeam }}</span>
        </div>
        <div class="pick-bar-container">
          <div class="pick-bar" :style="{ width: `${distribution.away}%` }"></div>
        </div>
        <div class="pick-pct">{{ distribution.away }}%</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.crowd-picks {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  padding: var(--space-6);
  background: var(--color-surface);
  border: var(--border-hairline);
}

[data-theme="dark"] .crowd-picks {
  border-radius: var(--radius-lg);
}

.crowd-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.crowd-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.crowd-subtitle {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.empty-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.picks-grid {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.pick-item {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .pick-item {
  border-radius: var(--radius-md);
}

.pick-item.leading {
  background: var(--color-pick-bg);
  border-color: var(--color-accent);
}

[data-theme="dark"] .pick-item.leading {
  box-shadow: 0 0 20px rgba(94,106,210,0.2);
}

[data-theme="light"] .pick-item.leading {
  background: transparent;
  border-width: 2px;
}

.pick-label {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 120px;
}

.pick-outcome {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  background: var(--color-bg);
  border: var(--border-hairline);
  color: var(--color-text-primary);
}

[data-theme="dark"] .pick-outcome {
  border-radius: var(--radius-sm);
}

.pick-item.leading .pick-outcome {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: var(--color-accent-text);
}

.pick-team {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pick-bar-container {
  position: relative;
  height: 8px;
  background: var(--color-bg);
  border: var(--border-hairline);
  overflow: hidden;
}

[data-theme="dark"] .pick-bar-container {
  border-radius: var(--radius-sm);
}

.pick-bar {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: var(--color-accent);
  transition: width var(--duration-slow) var(--ease-default);
}

[data-theme="dark"] .pick-bar {
  box-shadow: var(--shadow-accent);
}

[data-theme="light"] .pick-bar {
  background: #000;
}

.pick-pct {
  font-family: var(--font-mono);
  font-size: var(--text-base);
  font-weight: var(--font-bold);
  color: var(--color-text-secondary);
  min-width: 48px;
  text-align: right;
}

.pick-item.leading .pick-pct {
  color: var(--color-accent);
}

@media (max-width: 768px) {
  .pick-label {
    min-width: 100px;
  }

  .pick-team {
    font-size: var(--text-xs);
  }

  .pick-pct {
    font-size: var(--text-sm);
    min-width: 40px;
  }
}
</style>
