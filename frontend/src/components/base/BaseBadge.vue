<script setup>
defineProps({
  status: {
    type: String,
    required: true,
    validator: v => ['scheduled', 'live', 'finished'].includes(v)
  }
})
</script>

<template>
  <span
    class="badge"
    :class="`badge--${status}`"
    :aria-label="`Status: ${status}`"
  >
    <slot>{{ status }}</slot>
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-1) var(--space-3);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
  line-height: var(--leading-none);
}

/* ─── Dark theme ─── */
[data-theme="dark"] .badge {
  border-radius: var(--radius-full);
}

[data-theme="dark"] .badge--scheduled {
  background: rgba(138,143,152,0.15);
  color: var(--color-status-scheduled);
  border: 1px solid rgba(138,143,152,0.3);
}

[data-theme="dark"] .badge--live {
  background: rgba(94,106,210,0.15);
  color: var(--color-status-live);
  border: 1px solid var(--color-border-accent);
  box-shadow: 0 0 10px rgba(94,106,210,0.3);
  animation: pulse-dark 2s ease-in-out infinite;
}

[data-theme="dark"] .badge--live::before {
  content: '';
  width: 6px;
  height: 6px;
  background: var(--color-status-live);
  border-radius: 50%;
  box-shadow: 0 0 8px currentColor;
}

[data-theme="dark"] .badge--finished {
  background: rgba(255,255,255,0.05);
  color: var(--color-status-finished);
  border: 1px solid rgba(255,255,255,0.1);
}

/* ─── Light theme ─── */
[data-theme="light"] .badge {
  border-radius: 0;
}

[data-theme="light"] .badge--scheduled {
  background: transparent;
  color: var(--color-status-scheduled);
  border: 1px solid var(--color-border);
}

[data-theme="light"] .badge--live {
  background: var(--color-status-live);
  color: #fff;
  border: 2px solid var(--color-status-live);
  font-weight: var(--font-bold);
}

[data-theme="light"] .badge--live::before {
  content: '●';
  font-size: var(--text-xs);
}

[data-theme="light"] .badge--finished {
  background: transparent;
  color: var(--color-status-finished);
  border: 1px solid var(--color-border);
}

/* ─── Animations ─── */
@keyframes pulse-dark {
  0%, 100% {
    box-shadow: 0 0 10px rgba(94,106,210,0.3);
    border-color: rgba(94,106,210,0.3);
  }
  50% {
    box-shadow: 0 0 20px rgba(94,106,210,0.5);
    border-color: rgba(94,106,210,0.5);
  }
}
</style>
