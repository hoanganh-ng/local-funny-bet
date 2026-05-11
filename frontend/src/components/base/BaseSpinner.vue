<script setup>
defineProps({
  size: {
    type: String,
    default: 'md',
    validator: v => ['sm', 'md', 'lg'].includes(v)
  }
})
</script>

<template>
  <div
    class="spinner"
    :class="`spinner--${size}`"
    role="status"
    aria-label="Loading"
  >
    <span class="spinner__circle" aria-hidden="true"></span>
    <span class="sr-only">Loading...</span>
  </div>
</template>

<style scoped>
.spinner {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.spinner__circle {
  display: block;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-accent);
  animation: spin 800ms linear infinite;
}

/* ─── Size variants ─── */
.spinner--sm .spinner__circle {
  width: 16px;
  height: 16px;
}

.spinner--md .spinner__circle {
  width: 24px;
  height: 24px;
}

.spinner--lg .spinner__circle {
  width: 40px;
  height: 40px;
  border-width: 3px;
}

/* ─── Dark theme ─── */
[data-theme="dark"] .spinner__circle {
  border-radius: var(--radius-full);
}

/* ─── Light theme ─── */
[data-theme="light"] .spinner__circle {
  border-radius: 0;
  border-color: #E5E5E5;
  border-top-color: #000;
}

/* ─── Animation ─── */
@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ─── Accessibility ─── */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
