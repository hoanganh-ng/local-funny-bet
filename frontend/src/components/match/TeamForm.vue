<script setup>
const props = defineProps({
  form: {
    type: Array,
    default: () => []
  }
})

function getResultClass(result) {
  const map = {
    'W': 'win',
    'D': 'draw',
    'L': 'loss'
  }
  return map[result] || 'unknown'
}

function getResultLabel(result) {
  return result === 'W' ? 'Win' : result === 'D' ? 'Draw' : result === 'L' ? 'Loss' : '?'
}
</script>

<template>
  <div class="team-form">
    <div
      v-for="(result, index) in form"
      :key="index"
      class="form-indicator"
      :class="getResultClass(result)"
      :title="getResultLabel(result)"
    >
      {{ result }}
    </div>
  </div>
</template>

<style scoped>
.team-form {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}

.form-indicator {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  border: var(--border-hairline);
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .form-indicator {
  border-radius: var(--radius-sm);
}

.form-indicator.win {
  background: var(--color-correct);
  border-color: var(--color-correct);
  color: var(--color-accent-text);
}

[data-theme="dark"] .form-indicator.win {
  background: rgba(63,185,80,0.2);
  box-shadow: 0 0 12px rgba(63,185,80,0.2);
}

.form-indicator.draw {
  background: var(--color-bg-elevated);
  border-color: var(--color-border);
  color: var(--color-text-secondary);
}

.form-indicator.loss {
  background: var(--color-surface);
  border-color: var(--color-border);
  color: var(--color-text-disabled);
  opacity: 0.6;
}

[data-theme="light"] .form-indicator.win {
  background: #000;
  border-color: #000;
  color: #fff;
}

[data-theme="light"] .form-indicator.draw {
  background: transparent;
  border-color: #000;
  color: #000;
}

[data-theme="light"] .form-indicator.loss {
  background: transparent;
  border-color: rgba(0,0,0,0.2);
  color: rgba(0,0,0,0.4);
}
</style>
