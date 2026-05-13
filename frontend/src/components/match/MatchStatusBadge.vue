<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (val) => ['scheduled', 'live', 'finished'].includes(val)
  }
})

const badgeClass = computed(() => `badge--${props.status}`)
</script>

<template>
  <span class="badge" :class="badgeClass" :aria-label="`Status: ${status}`">
    <span v-if="status === 'live'" class="live-dot" aria-hidden="true" />
    <slot>{{ status === 'finished' ? 'FT' : status }}</slot>
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: 3px 10px;
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  line-height: var(--leading-none);
  border-radius: var(--radius-full);
  border: none;
}

.badge--scheduled {
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-text-secondary);
}

[data-theme="light"] .badge--scheduled {
  background: rgba(0, 0, 0, 0.06);
  color: var(--color-text-secondary);
}

.badge--live {
  background: rgba(94, 106, 210, 0.15);
  color: var(--color-status-live);
}

[data-theme="light"] .badge--live {
  background: rgba(0, 0, 0, 0.08);
  color: var(--color-status-live);
}

.badge--finished {
  background: transparent;
  color: var(--color-text-disabled);
}

/* Live dot */
.live-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  background: var(--color-status-live);
  flex-shrink: 0;
  animation: dot-pulse 1.5s ease-in-out infinite;
}

[data-theme="light"] .live-dot {
  background: var(--color-status-live);
}

@keyframes dot-pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.4; }
}
</style>
