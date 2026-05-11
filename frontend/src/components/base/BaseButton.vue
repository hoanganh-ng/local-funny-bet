<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: v => ['primary', 'secondary', 'ghost', 'danger'].includes(v)
  },
  size: {
    type: String,
    default: 'md',
    validator: v => ['sm', 'md', 'lg'].includes(v)
  },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false }
})
</script>

<template>
  <button
    class="btn"
    :class="[`btn--${variant}`, `btn--${size}`]"
    :disabled="disabled || loading"
    :aria-busy="loading"
  >
    <span v-if="loading" class="btn__spinner" aria-hidden="true"></span>
    <span :class="{ 'btn__content--loading': loading }">
      <slot />
    </span>
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-family: var(--font-body);
  font-weight: var(--font-medium);
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-hover);
  position: relative;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

/* ─── Size variants ─── */
.btn--sm {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
}

.btn--md {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
}

.btn--lg {
  padding: var(--space-4) var(--space-8);
  font-size: var(--text-lg);
}

/* ─── Dark theme styles ─── */
[data-theme="dark"] .btn--primary {
  background: var(--color-accent);
  color: #fff;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-accent);
}

[data-theme="dark"] .btn--primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
  transform: translateY(-2px);
}

[data-theme="dark"] .btn--primary:active:not(:disabled) {
  transform: translateY(0);
}

[data-theme="dark"] .btn--secondary {
  background: var(--color-surface);
  color: var(--color-text-primary);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
}

[data-theme="dark"] .btn--secondary:hover:not(:disabled) {
  background: var(--color-surface-hover);
  border-color: var(--color-border-hover);
  transform: translateY(-1px);
}

[data-theme="dark"] .btn--ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border-radius: var(--radius-md);
}

[data-theme="dark"] .btn--ghost:hover:not(:disabled) {
  color: var(--color-text-primary);
  background: var(--color-surface);
}

[data-theme="dark"] .btn--danger {
  background: var(--color-danger);
  color: #fff;
  border-radius: var(--radius-md);
  box-shadow: 0 0 20px rgba(229,83,75,0.3);
}

[data-theme="dark"] .btn--danger:hover:not(:disabled) {
  filter: brightness(1.1);
  transform: translateY(-2px);
}

/* ─── Light theme styles ─── */
[data-theme="light"] .btn--primary {
  background: #000;
  color: #fff;
  border-radius: 0;
}

[data-theme="light"] .btn--primary:hover:not(:disabled) {
  background: #fff;
  color: #000;
  outline: 2px solid #000;
  outline-offset: -2px;
}

[data-theme="light"] .btn--secondary {
  background: #fff;
  color: #000;
  border: 1px solid #000;
  border-radius: 0;
}

[data-theme="light"] .btn--secondary:hover:not(:disabled) {
  background: #000;
  color: #fff;
}

[data-theme="light"] .btn--ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border-radius: 0;
}

[data-theme="light"] .btn--ghost:hover:not(:disabled) {
  color: var(--color-text-primary);
  text-decoration: underline;
  text-underline-offset: 4px;
}

[data-theme="light"] .btn--danger {
  background: #000;
  color: #fff;
  border: 2px solid #000;
  border-radius: 0;
}

[data-theme="light"] .btn--danger:hover:not(:disabled) {
  background: #fff;
  color: #000;
}

/* ─── Loading state ─── */
.btn__spinner {
  width: 1em;
  height: 1em;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 600ms linear infinite;
}

[data-theme="light"] .btn__spinner {
  border-radius: 0;
}

.btn__content--loading {
  opacity: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
