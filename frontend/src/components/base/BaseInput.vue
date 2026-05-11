<script setup>
defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  error: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  type: { type: String, default: 'text' }
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="input-wrapper">
    <label v-if="label" class="input-label" :for="`input-${label}`">
      {{ label }}
    </label>
    <input
      :id="`input-${label}`"
      class="input"
      :class="{ 'input--error': error }"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :aria-invalid="!!error"
      :aria-describedby="error ? `error-${label}` : undefined"
      @input="$emit('update:modelValue', $event.target.value)"
    />
    <span v-if="error" :id="`error-${label}`" class="input-error">
      {{ error }}
    </span>
  </div>
</template>

<style scoped>
.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input-label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.input {
  width: 100%;
  padding: var(--space-3) var(--space-4);
  font-family: var(--font-body);
  font-size: var(--text-base);
  color: var(--color-text-primary);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  transition: all var(--duration-fast) var(--ease-hover);
}

.input::placeholder {
  color: var(--color-text-disabled);
}

.input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-error {
  font-size: var(--text-sm);
  color: var(--color-danger);
  font-family: var(--font-mono);
}

/* ─── Dark theme ─── */
[data-theme="dark"] .input {
  border-radius: var(--radius-md);
}

[data-theme="dark"] .input:focus {
  outline: none;
  border-color: var(--color-border-accent);
  box-shadow: 0 0 0 3px var(--color-accent-glow);
}

[data-theme="dark"] .input--error {
  border-color: var(--color-danger);
}

[data-theme="dark"] .input--error:focus {
  box-shadow: 0 0 0 3px rgba(229,83,75,0.2);
}

/* ─── Light theme ─── */
[data-theme="light"] .input {
  border-radius: 0;
  border: none;
  border-bottom: 2px solid var(--color-border);
  background: transparent;
  padding-left: 0;
  padding-right: 0;
}

[data-theme="light"] .input:focus {
  outline: none;
  border-bottom-width: 4px;
  border-bottom-color: var(--color-accent);
}

[data-theme="light"] .input--error {
  border-bottom-color: var(--color-danger);
}

[data-theme="light"] .input--error:focus {
  border-bottom-width: 4px;
}

[data-theme="light"] .input-error {
  color: var(--color-text-primary);
  font-weight: var(--font-semibold);
}
</style>
