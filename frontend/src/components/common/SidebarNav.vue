<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth.store.js'
import { useLeaderboardStore } from '../../store/leaderboard.store.js'
import { useTheme } from '../../composables/useTheme.js'
import BaseButton from '../base/BaseButton.vue'
import { ref } from 'vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const leaderboardStore = useLeaderboardStore()
const { theme, toggle: toggleTheme } = useTheme()

function isActive(path) {
  return route.path === path || route.path.startsWith(path + '/')
}

function goToBoard(id) {
  router.push(`/leaderboards/${id}`)
}

const isLoggingOut = ref(false)

async function handleLogout() {
  if (isLoggingOut.value) return

  isLoggingOut.value = true
  try {
    await authStore.logout()
    router.push('/login')
  } catch (err) {
    console.error('Logout failed:', err)
    router.push('/login')  // Fail-safe: always redirect
  } finally {
    isLoggingOut.value = false
  }
}
</script>

<template>
  <aside class="sidebar">
    <!-- Logo + Tournament info -->
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">⚽</div>
        <div class="logo-text">
          <span class="logo-title">Predictor</span>
          <span class="logo-subtitle">Internal</span>
        </div>
      </div>

      <div class="tournament-info">
        <p class="tournament-stage">Group Stage · Matchday 2</p>
        <p class="tournament-name">Continental Cup 2026</p>
        <p class="tournament-progress">8 / 16 matches played</p>
      </div>
    </div>

    <!-- Main navigation -->
    <nav class="nav-main">
      <router-link to="/" class="nav-link" :class="{ active: isActive('/') && route.path === '/' }">
        <span class="nav-icon">⚽</span>
        <span class="nav-label">Matches</span>
      </router-link>

      <router-link to="/leaderboards" class="nav-link" :class="{ active: isActive('/leaderboards') }">
        <span class="nav-icon">🏆</span>
        <span class="nav-label">Leaderboards</span>
        <span v-if="leaderboardStore.leaderboards.length" class="nav-count">
          {{ leaderboardStore.leaderboards.length }}
        </span>
      </router-link>

      <router-link to="/history" class="nav-link" :class="{ active: isActive('/history') }">
        <span class="nav-icon">📊</span>
        <span class="nav-label">History</span>
      </router-link>
    </nav>

    <!-- Your boards list -->
    <div class="boards-section">
      <div class="boards-header">
        <span class="boards-title">Your Boards</span>
        <button class="boards-add" @click="$router.push('/leaderboards/new')">+</button>
      </div>

      <div v-if="leaderboardStore.leaderboards.length === 0" class="boards-empty">
        <p>No boards yet</p>
      </div>

      <div v-else class="boards-list">
        <button
          v-for="board in leaderboardStore.leaderboards"
          :key="board.id"
          class="board-item"
          :class="{ active: route.params.id === String(board.id) }"
          @click="goToBoard(board.id)"
        >
          <span class="board-indicator">●</span>
          <span class="board-name">{{ board.name }}</span>
          <span class="board-rank">#{{ board.userRank || 3 }}</span>
        </button>
      </div>
    </div>

    <!-- Profile footer -->
    <div class="sidebar-footer">
      <div class="profile-card">
        <div class="profile-avatar">
          {{ authStore.user?.name?.charAt(0) || 'A' }}
        </div>
        <div class="profile-info">
          <p class="profile-name">{{ authStore.user?.name || 'Alex Chen' }}</p>
          <p class="profile-email">{{ authStore.user?.email || 'alex.chen@acme.co' }}</p>
        </div>
        <button class="profile-menu" @click="toggleTheme">
          <span v-if="theme === 'dark'">☀️</span>
          <span v-else>🌙</span>
        </button>
      </div>

      <BaseButton
        variant="ghost"
        size="sm"
        :loading="isLoggingOut"
        @click="handleLogout"
        class="logout-btn"
      >
        Log out
      </BaseButton>
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
  overflow-y: auto;
  position: sticky;
  top: 0;
}

/* ──── Header ──── */
.sidebar-header {
  padding: var(--space-6);
  border-bottom: var(--border-hairline);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-2xl);
  background: var(--color-accent);
  border-radius: var(--radius-md);
}

[data-theme="light"] .logo-icon {
  border-radius: 0;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-title {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.logo-subtitle {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.tournament-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.tournament-stage {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.tournament-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.tournament-progress {
  font-size: var(--text-xs);
  color: var(--color-text-disabled);
}

/* ──── Main Nav ──── */
.nav-main {
  padding: var(--space-6) var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.nav-link {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-default);
  cursor: pointer;
  position: relative;
}

[data-theme="dark"] .nav-link {
  border-radius: var(--radius-md);
}

.nav-link:hover {
  color: var(--color-text-primary);
  background: var(--color-surface);
}

.nav-link.active {
  color: var(--color-text-primary);
  background: var(--color-accent);
  font-weight: var(--font-semibold);
}

[data-theme="dark"] .nav-link.active {
  box-shadow: var(--shadow-accent);
}

[data-theme="light"] .nav-link.active {
  background: var(--color-accent);
  color: var(--color-accent-text);
}

.nav-icon {
  font-size: var(--text-xl);
}

.nav-label {
  flex: 1;
}

.nav-count {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  padding: 2px var(--space-2);
  background: var(--color-surface);
  border-radius: var(--radius-pill);
}

[data-theme="dark"] .nav-link.active .nav-count {
  background: rgba(255,255,255,0.2);
}

/* ──── Boards Section ──── */
.boards-section {
  flex: 1;
  padding: var(--space-4);
  border-top: var(--border-hairline);
  overflow-y: auto;
}

.boards-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-2) var(--space-2) var(--space-4);
}

.boards-title {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.boards-add {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  background: transparent;
  border: var(--border-hairline);
  cursor: pointer;
  transition: all var(--duration-fast);
}

[data-theme="dark"] .boards-add {
  border-radius: var(--radius-sm);
}

.boards-add:hover {
  color: var(--color-text-primary);
  border-color: var(--color-border-hover);
  background: var(--color-surface);
}

.boards-empty {
  padding: var(--space-6) var(--space-4);
  text-align: center;
  font-size: var(--text-sm);
  color: var(--color-text-disabled);
}

.boards-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.board-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  text-align: left;
  cursor: pointer;
  transition: all var(--duration-fast);
  width: 100%;
}

[data-theme="dark"] .board-item {
  border-radius: var(--radius-sm);
}

.board-item:hover {
  color: var(--color-text-primary);
  background: var(--color-surface);
}

.board-item.active {
  color: var(--color-text-primary);
  background: var(--color-surface-hover);
}

.board-indicator {
  font-size: 8px;
  color: var(--color-accent);
}

.board-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.board-rank {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-text-disabled);
}

/* ──── Footer ──── */
.sidebar-footer {
  padding: var(--space-4);
  border-top: var(--border-hairline);
}

.profile-card {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface);
  cursor: pointer;
  transition: all var(--duration-fast);
}

[data-theme="dark"] .profile-card {
  border-radius: var(--radius-md);
}

.profile-card:hover {
  background: var(--color-surface-hover);
}

.profile-avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-accent-text);
  background: var(--color-accent);
  border-radius: var(--radius-pill);
}

[data-theme="light"] .profile-avatar {
  border-radius: 0;
}

.profile-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.profile-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-email {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-menu {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-lg);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast);
}

[data-theme="dark"] .profile-menu {
  border-radius: var(--radius-sm);
}

.profile-menu:hover {
  background: var(--color-bg-elevated);
}

.logout-btn {
  width: 100%;
  margin-top: var(--space-2);
  justify-content: center;
  color: var(--color-text-secondary);
  padding: var(--space-2);
}

.logout-btn:hover {
  color: var(--color-text-primary);
}
</style>
