<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()

function isActive(path) {
  if (path === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(path)
}
</script>

<template>
  <nav class="bottom-tab-bar">
    <router-link to="/" class="tab" :class="{ active: isActive('/') }">
      <span class="tab-icon">⚽</span>
      <span class="tab-label">Matches</span>
    </router-link>

    <router-link to="/leaderboards" class="tab" :class="{ active: isActive('/leaderboards') }">
      <span class="tab-icon">🏆</span>
      <span class="tab-label">Boards</span>
    </router-link>

    <router-link to="/history" class="tab" :class="{ active: isActive('/history') }">
      <span class="tab-icon">📊</span>
      <span class="tab-label">History</span>
    </router-link>
  </nav>
</template>

<style scoped>
.bottom-tab-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: var(--bottom-nav-h);
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  background: var(--color-bg-elevated);
  border-top: var(--border-hairline);
  z-index: 100;
}

[data-theme="dark"] .bottom-tab-bar {
  backdrop-filter: blur(20px);
  background: rgba(10,10,12,0.95);
}

.tab {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  text-decoration: none;
  color: var(--color-text-secondary);
  transition: all var(--duration-fast) var(--ease-default);
  position: relative;
}

.tab::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 3px;
  background: var(--color-accent);
  transition: width var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .tab::before {
  box-shadow: var(--shadow-accent);
}

.tab.active::before {
  width: 40px;
}

.tab.active {
  color: var(--color-accent);
}

[data-theme="light"] .tab.active {
  color: var(--color-text-primary);
  font-weight: var(--font-semibold);
}

.tab-icon {
  font-size: var(--text-2xl);
}

.tab-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

@media (min-width: 768px) {
  .bottom-tab-bar {
    display: none;
  }
}
</style>
