<script setup>
defineProps({
  hoverable: { type: Boolean, default: false }
})
</script>

<template>
  <div class="card" :class="{ 'card--hoverable': hoverable }">
    <slot />
  </div>
</template>

<style scoped>
.card {
  padding: var(--space-6);
  transition: all var(--duration-base) var(--ease-default);
}

/* ─── Dark theme ─── */
[data-theme="dark"] .card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  backdrop-filter: blur(8px);
}

[data-theme="dark"] .card--hoverable:hover {
  border-color: var(--color-border-hover);
  box-shadow: var(--shadow-lg);
  transform: translateY(-4px);
}

/* ─── Light theme ─── */
[data-theme="light"] .card {
  background: var(--color-surface);
  border: var(--border-thin);
  border-radius: 0;
  box-shadow: none;
}

[data-theme="light"] .card--hoverable {
  transition: background var(--duration-fast), color var(--duration-fast);
}

[data-theme="light"] .card--hoverable:hover {
  background: var(--color-surface-hover);
  color: #fff;
}

/* Override child text color in light theme hover */
[data-theme="light"] .card--hoverable:hover :deep(*) {
  color: inherit;
}
</style>
