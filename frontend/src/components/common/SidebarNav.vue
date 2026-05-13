<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  tournament: String,
  activeRoute: String,
  leaderboardCount: Number,
  user: Object,
  mobileVisible: Boolean
})

const emit = defineEmits(['switch-tournament', 'toggle-theme', 'logout', 'nav-click'])

const router = useRouter()
const userInitial = computed(() => props.user?.initial || props.user?.name?.charAt(0) || 'A')

function navigateTo(path) {
  router.push(path)
  emit('nav-click')
}
</script>

<template>
  <aside class="sidebar" :data-mobile-visible="mobileVisible">
    <!-- Logo -->
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">
          <i class="ti ti-ball-football" aria-hidden="true"></i>
        </div>
        <span class="logo-title">Predictor</span>
      </div>

      <!-- Tournament switcher pill -->
      <button class="tournament-pill" @click="emit('switch-tournament')" aria-label="Switch tournament">
        <i class="ti ti-world" aria-hidden="true"></i>
        <span class="tournament-name">{{ tournament || 'FIFA World Cup 2026' }}</span>
        <i class="ti ti-chevron-down" aria-hidden="true"></i>
      </button>
    </div>

    <div class="divider"></div>

    <!-- Main navigation -->
    <nav class="nav-main" aria-label="Main navigation">
      <div class="nav-item" :class="{ active: activeRoute === 'matches' }" @click="navigateTo('/')">
        <i class="ti ti-ball-football nav-icon" aria-hidden="true"></i>
        <span class="nav-label">Matches</span>
      </div>

      <div class="nav-item" :class="{ active: activeRoute === 'leaderboards' }" @click="navigateTo('/leaderboards')">
        <i class="ti ti-trophy nav-icon" aria-hidden="true"></i>
        <span class="nav-label">Leaderboards</span>
        <span v-if="leaderboardCount" class="nav-count">{{ leaderboardCount }}</span>
      </div>

      <div class="nav-item" :class="{ active: activeRoute === 'history' }" @click="navigateTo('/history')">
        <i class="ti ti-history nav-icon" aria-hidden="true"></i>
        <span class="nav-label">History</span>
      </div>
    </nav>

    <div class="spacer"></div>

    <!-- User footer -->
    <div class="sidebar-footer">
      <div class="divider"></div>
      <div class="user-row">
        <div class="user-avatar">{{ userInitial }}</div>
        <div class="user-info">
          <p class="user-name">{{ user?.name || 'User' }}</p>
          <p class="user-email">{{ user?.email || '' }}</p>
        </div>
        <div class="user-actions">
          <button class="icon-btn" @click="emit('toggle-theme')" aria-label="Toggle theme">
            <i class="ti ti-sun" aria-hidden="true"></i>
          </button>
          <button class="icon-btn" @click="emit('logout')" aria-label="Log out">
            <i class="ti ti-logout" aria-hidden="true"></i>
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: var(--sidebar-w);
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-elevated);
  border-right: var(--border-hairline);
  padding: var(--space-4) 0;
  box-sizing: border-box;
  position: sticky;
  top: 0;
}

/* ──── Header ──── */
.sidebar-header {
  padding: 0 var(--space-4) var(--space-4);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.logo-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-accent);
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.logo-icon i {
  color: white;
  font-size: 14px;
}

.logo-title {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

/* ──── Tournament Pill ──── */
.tournament-pill {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--color-surface);
  border: 0.5px solid var(--color-border-accent);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 150ms ease;
  width: 100%;
}

.tournament-pill:hover {
  background: var(--color-surface-hover);
  border-color: var(--color-accent-glow);
}

.tournament-pill > i:first-child {
  color: var(--color-accent);
  font-size: 13px;
  flex-shrink: 0;
}

.tournament-name {
  color: var(--color-text-primary);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-align: left;
}

.tournament-pill > i:last-child {
  color: var(--color-text-disabled);
  font-size: 11px;
  flex-shrink: 0;
}

/* ──── Divider ──── */
.divider {
  height: 0.5px;
  background: var(--color-border);
  margin: var(--space-3) var(--space-3);
}

/* ──── Main Nav ──── */
.nav-main {
  padding: 0 var(--space-2);
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-3);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background 120ms ease;
}

.nav-item:hover {
  background: var(--color-surface);
}

.nav-item:hover .nav-label,
.nav-item:hover .nav-icon {
  color: var(--color-text-primary);
}

.nav-item.active {
  background: var(--color-pick-bg);
}

.nav-icon {
  color: var(--color-text-secondary);
  font-size: 15px;
  flex-shrink: 0;
  transition: color 120ms ease;
}

.nav-item.active .nav-icon {
  color: var(--color-accent);
}

.nav-label {
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  flex: 1;
  transition: color 120ms ease;
}

.nav-item.active .nav-label {
  color: var(--color-text-primary);
}

.nav-count {
  background: var(--color-border);
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
  line-height: 1.5;
}

/* ──── Spacer ──── */
.spacer {
  flex: 1;
}

/* ──── Footer ──── */
.sidebar-footer {
  padding: 0 var(--space-2);
}

.user-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background 120ms ease;
}

.user-row:hover {
  background: var(--color-surface);
}

.user-avatar {
  width: 27px;
  height: 27px;
  border-radius: var(--radius-full);
  background: var(--color-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  color: white;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  margin: 0;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  color: var(--color-text-disabled);
  font-size: var(--text-xs);
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 1px;
  flex-shrink: 0;
}

.icon-btn {
  background: none;
  border: none;
  color: var(--color-text-disabled);
  cursor: pointer;
  padding: var(--space-1);
  border-radius: var(--radius-sm);
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 120ms ease;
}

.icon-btn:hover {
  background: var(--color-surface-hover);
  color: var(--color-text-primary);
}

.icon-btn i {
  font-size: 14px;
}

/* Mobile: hide sidebar by default, show via prop/class */
@media (max-width: 768px) {
  .sidebar {
    display: none;
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 100;
  }

  .sidebar[data-mobile-visible="true"] {
    display: flex;
  }
}
</style>
