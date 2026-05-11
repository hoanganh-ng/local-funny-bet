<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    default: null,
    validator: (val) => val === null || ['home_win', 'draw', 'away_win'].includes(val)
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const selectOutcome = (outcome) => {
  if (props.disabled) return
  emit('update:modelValue', outcome)
}
</script>

<template>
  <div class="outcome-picker">
    <button
      type="button"
      class="outcome-btn"
      :class="{
        selected: modelValue === 'home_win',
        disabled: disabled
      }"
      :disabled="disabled"
      @click="selectOutcome('home_win')"
    >
      Home Win
    </button>
    <button
      type="button"
      class="outcome-btn"
      :class="{
        selected: modelValue === 'draw',
        disabled: disabled
      }"
      :disabled="disabled"
      @click="selectOutcome('draw')"
    >
      Draw
    </button>
    <button
      type="button"
      class="outcome-btn"
      :class="{
        selected: modelValue === 'away_win',
        disabled: disabled
      }"
      :disabled="disabled"
      @click="selectOutcome('away_win')"
    >
      Away Win
    </button>
  </div>
</template>

<style scoped>
.outcome-picker {
  display: flex;
  gap: var(--space-2);
}

.outcome-btn {
  flex: 1;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  font-family: var(--font-body);
  color: var(--color-text-primary);
  background: transparent;
  border: var(--border-thin);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.outcome-btn:hover:not(.disabled) {
  background: var(--color-surface-hover);
  border-color: var(--color-border-hover);
}

.outcome-btn.selected {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: var(--color-bg);
}

[data-theme="light"] .outcome-btn.selected {
  background: var(--color-bg);
  border-color: var(--color-border-strong);
  color: var(--color-text-primary);
  border-width: 2px;
}

.outcome-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  color: var(--color-text-disabled);
  border-color: var(--color-border);
}
</style>
