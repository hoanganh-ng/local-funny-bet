<script setup>
import { computed } from 'vue'
import BaseBadge from '../base/BaseBadge.vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (val) => ['scheduled', 'live', 'finished'].includes(val)
  },
  kickoffAt: {
    type: String,
    default: null
  }
})

const badgeText = computed(() => {
  if (props.status === 'scheduled' && props.kickoffAt) {
    const date = new Date(props.kickoffAt)
    const month = date.toLocaleDateString('en-US', { month: 'short' })
    const day = date.getDate()
    const time = date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
    return `${month} ${day} · ${time}`
  }
  if (props.status === 'live') return 'LIVE'
  if (props.status === 'finished') return 'FT'
  return props.status
})

const badgeClass = computed(() => {
  if (props.status === 'live') return 'live'
  if (props.status === 'finished') return 'finished'
  return 'scheduled'
})
</script>

<template>
  <BaseBadge :class="badgeClass">{{ badgeText }}</BaseBadge>
</template>

<style scoped>
.scheduled {
  color: var(--color-status-scheduled);
  border-color: var(--color-border);
}

.live {
  color: var(--color-status-live);
  border-color: var(--color-status-live);
  animation: pulse var(--duration-slow) var(--ease-default) infinite;
}

[data-theme="light"] .live {
  font-weight: var(--font-bold);
  border-width: 2px;
  animation: none;
}

.finished {
  color: var(--color-status-finished);
  border-color: var(--color-border);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    box-shadow: 0 0 0 0 var(--color-accent-glow);
  }
  50% {
    opacity: 0.8;
    box-shadow: 0 0 20px var(--color-accent-glow);
  }
}
</style>
