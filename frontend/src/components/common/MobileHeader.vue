<script setup>
import { useUiStore } from '../../store/ui.store.js'
import { useTournament } from '../../composables/useTournament.js'

const uiStore = useUiStore()
const { tournament } = useTournament()
</script>

<template>
  <header class="mobile-header">
    <div class="tournament-name">
      {{ tournament?.name || 'Continental Cup 2026' }}
    </div>
    <button
      class="toggle-btn"
      @click="uiStore.toggleMobileNav"
      :aria-label="uiStore.showSidebarOnMobile ? 'Close sidebar' : 'Open sidebar'"
    >
      <span class="toggle-icon">{{ uiStore.showSidebarOnMobile ? '✕' : '☰' }}</span>
    </button>
  </header>
</template>

<style scoped>
.mobile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--mobile-header-h);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-4);
  background: var(--color-bg-elevated);
  border-bottom: var(--border-hairline);
  z-index: 99;
}

[data-theme="dark"] .mobile-header {
  backdrop-filter: blur(20px);
  background: rgba(10,10,12,0.95);
}

.tournament-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toggle-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast);
  flex-shrink: 0;
}

[data-theme="dark"] .toggle-btn {
  border-radius: var(--radius-sm);
}

.toggle-btn:hover {
  background: var(--color-surface);
}

.toggle-btn:active {
  background: var(--color-surface-hover);
}

.toggle-icon {
  font-size: var(--text-2xl);
  color: var(--color-text-primary);
}

/* Hide on desktop */
@media (min-width: 769px) {
  .mobile-header {
    display: none;
  }
}
</style>
