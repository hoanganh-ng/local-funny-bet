<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    default: null,
    validator: (val) => val === null || ['home_win', 'draw', 'away_win'].includes(val)
  },
  locked: {
    type: Boolean,
    default: false
  },
  result: {
    type: String,
    default: null,
    validator: (val) => val === null || ['home_win', 'draw', 'away_win'].includes(val)
  },
  points: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const selectOutcome = (outcome) => {
  if (props.locked) return
  emit('update:modelValue', outcome)
}

const getOutcomeState = (outcome) => {
  if (!props.result) return null
  if (!props.modelValue) return null
  if (props.modelValue === outcome) {
    return props.result === outcome ? 'correct' : 'wrong'
  }
  return null
}
</script>

<template>
  <div class="outcome-picker">
    <button
      type="button"
      class="outcome-btn outcome-btn--home"
      :class="{
        selected: modelValue === 'home_win',
        locked: locked,
        correct: getOutcomeState('home_win') === 'correct',
        wrong: getOutcomeState('home_win') === 'wrong'
      }"
      :disabled="locked"
      @click="selectOutcome('home_win')"
    >
      <span class="outcome-label">1</span>
      <span class="outcome-sublabel">HOME</span>
      <span v-if="getOutcomeState('home_win') === 'correct'" class="outcome-points">
        +{{ points }}
      </span>
    </button>

    <button
      type="button"
      class="outcome-btn outcome-btn--draw"
      :class="{
        selected: modelValue === 'draw',
        locked: locked,
        correct: getOutcomeState('draw') === 'correct',
        wrong: getOutcomeState('draw') === 'wrong'
      }"
      :disabled="locked"
      @click="selectOutcome('draw')"
    >
      <span class="outcome-label">X</span>
      <span class="outcome-sublabel">DRAW</span>
      <span v-if="getOutcomeState('draw') === 'correct'" class="outcome-points">
        +{{ points }}
      </span>
    </button>

    <button
      type="button"
      class="outcome-btn outcome-btn--away"
      :class="{
        selected: modelValue === 'away_win',
        locked: locked,
        correct: getOutcomeState('away_win') === 'correct',
        wrong: getOutcomeState('away_win') === 'wrong'
      }"
      :disabled="locked"
      @click="selectOutcome('away_win')"
    >
      <span class="outcome-label">2</span>
      <span class="outcome-sublabel">AWAY</span>
      <span v-if="getOutcomeState('away_win') === 'correct'" class="outcome-points">
        +{{ points }}
      </span>
    </button>
  </div>
</template>

<style scoped>
.outcome-picker {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
}

.outcome-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-4) var(--space-3);
  font-family: var(--font-body);
  background: var(--color-surface);
  border: var(--border-hairline);
  cursor: pointer;
  transition: background var(--duration-fast) var(--ease-default),
              border-color var(--duration-fast) var(--ease-default),
              box-shadow var(--duration-fast) var(--ease-default),
              transform var(--duration-fast) var(--ease-default);
  position: relative;
}

/* ──── Dark theme ──── */
[data-theme="dark"] .outcome-btn {
  border-radius: var(--radius-md);
}

[data-theme="dark"] .outcome-btn:hover:not(.locked) {
  background: var(--color-surface-hover);
  border-color: var(--color-border-hover);
  transform: translateY(-1px);
}

[data-theme="dark"] .outcome-btn.selected {
  background: var(--color-accent-glow);
  border-color: var(--color-border-accent);
  box-shadow: var(--shadow-accent);
}

[data-theme="dark"] .outcome-btn--home:hover:not(.locked):not(.selected) {
  background: var(--color-home-win);
}

[data-theme="dark"] .outcome-btn--draw:hover:not(.locked):not(.selected) {
  background: var(--color-draw);
}

[data-theme="dark"] .outcome-btn--away:hover:not(.locked):not(.selected) {
  background: var(--color-away-win);
}

[data-theme="dark"] .outcome-btn.correct {
  background: rgba(63, 185, 80, 0.15);
  border-color: var(--color-correct);
  box-shadow: 0 0 20px rgba(63, 185, 80, 0.2);
}

[data-theme="dark"] .outcome-btn.wrong {
  background: var(--color-surface);
  border-color: var(--color-border);
}

/* ──── Light theme ──── */
[data-theme="light"] .outcome-btn {
  border-radius: 0;
}

[data-theme="light"] .outcome-btn:hover:not(.locked) {
  border-color: var(--color-border-hover);
}

[data-theme="light"] .outcome-btn.selected {
  background: var(--color-accent-glow);
  color: var(--color-accent-text);
  border: var(--border-medium);
}

[data-theme="light"] .outcome-btn.correct {
  background: var(--color-correct);
  color: var(--color-accent-text);
  border: var(--border-medium);
  font-weight: var(--font-bold);
}

[data-theme="light"] .outcome-btn.wrong {
  color: var(--color-wrong);
  border-color: var(--color-border);
}

/* ──── Labels ──── */
.outcome-label {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

[data-theme="dark"] .outcome-btn.selected .outcome-label {
  color: var(--color-accent);
}

[data-theme="dark"] .outcome-btn.correct .outcome-label {
  color: var(--color-correct);
}

[data-theme="light"] .outcome-btn.selected .outcome-label,
[data-theme="light"] .outcome-btn.correct .outcome-label {
  color: inherit;
}

.outcome-sublabel {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

[data-theme="dark"] .outcome-btn.selected .outcome-sublabel {
  color: var(--color-text-secondary);
}

[data-theme="dark"] .outcome-btn.correct .outcome-sublabel {
  color: var(--color-correct);
}

[data-theme="light"] .outcome-btn.selected .outcome-sublabel,
[data-theme="light"] .outcome-btn.correct .outcome-sublabel {
  color: inherit;
  opacity: 0.7;
}

.outcome-points {
  position: absolute;
  top: var(--space-2);
  right: var(--space-2);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-correct);
}

[data-theme="light"] .outcome-points {
  color: inherit;
}

/* ──── Locked state ──── */
.outcome-btn.locked {
  cursor: not-allowed;
  opacity: 0.6;
}

.outcome-btn.wrong {
  opacity: 0.5;
}
</style>
