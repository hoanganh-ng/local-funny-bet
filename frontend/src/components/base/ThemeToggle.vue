<script setup>
import { ref, onMounted } from 'vue'

const theme = ref('dark')

onMounted(() => {
  theme.value = localStorage.getItem('theme') || 'dark'
  document.documentElement.setAttribute('data-theme', theme.value)
})

function toggle() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', theme.value)
  localStorage.setItem('theme', theme.value)
}
</script>

<template>
  <button
    class="theme-toggle"
    @click="toggle"
    :aria-label="`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`"
  >
    <span class="theme-toggle__icon">{{ theme === 'dark' ? '○' : '●' }}</span>
  </button>
</template>

<style scoped>
.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: color var(--duration-fast) var(--ease-hover);
  padding: 0;
}

.theme-toggle:hover {
  color: var(--color-text-primary);
}

.theme-toggle__icon {
  font-size: var(--text-2xl);
  line-height: 1;
}

/* ─── Dark theme ─── */
[data-theme="dark"] .theme-toggle {
  border-radius: var(--radius-md);
}

[data-theme="dark"] .theme-toggle:hover {
  background: var(--color-surface);
}

[data-theme="dark"] .theme-toggle:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}

/* ─── Light theme ─── */
[data-theme="light"] .theme-toggle {
  border-radius: 0;
}

[data-theme="light"] .theme-toggle:hover {
  background: transparent;
}

[data-theme="light"] .theme-toggle:focus-visible {
  outline: 2px solid #000;
  outline-offset: 2px;
}
</style>
