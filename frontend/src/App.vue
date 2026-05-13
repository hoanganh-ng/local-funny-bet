<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MobileHeader from './components/common/MobileHeader.vue'
import SidebarNav from './components/common/SidebarNav.vue'
import BottomTabBar from './components/common/BottomTabBar.vue'
import { useUiStore } from './store/ui.store.js'
import { useAuthStore } from './store/auth.store.js'
import { useLeaderboardStore } from './store/leaderboard.store.js'
import { useTournament } from './composables/useTournament.js'
import { useTheme } from './composables/useTheme.js'

const route = useRoute()
const router = useRouter()
const uiStore = useUiStore()
const authStore = useAuthStore()
const leaderboardStore = useLeaderboardStore()
const { tournament } = useTournament()
const { toggle: toggleTheme } = useTheme()

const showLayout = computed(() => {
  return route.meta.public !== true
})

const activeRoute = computed(() => {
  if (route.path === '/' || route.path.startsWith('/matches')) return 'matches'
  if (route.path.startsWith('/leaderboards')) return 'leaderboards'
  if (route.path.startsWith('/history')) return 'history'
  return ''
})

const sidebarUser = computed(() => ({
  name: authStore.user?.name || 'User',
  email: authStore.user?.email || '',
  initial: authStore.user?.name?.charAt(0) || 'U'
}))

function handleSwitchTournament() {
  console.log('TODO(global): tournament switcher')
}

function handleToggleTheme() {
  toggleTheme()
}

async function handleLogout() {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (err) {
    console.error('Logout failed:', err)
    router.push('/login')
  }
}
</script>

<template>
  <div id="app">
    <!-- Public routes: no layout -->
    <router-view v-if="!showLayout" />

    <!-- Authenticated routes: with sidebar/bottom nav -->
    <div v-else class="app-layout">
      <MobileHeader />

      <!-- Backdrop overlay for mobile sidebar -->
      <div
        v-if="uiStore.showSidebarOnMobile"
        class="mobile-sidebar-backdrop"
        @click="uiStore.closeMobileSidebar"
      />

      <SidebarNav
        class="app-sidebar"
        :tournament="tournament?.name"
        :active-route="activeRoute"
        :leaderboard-count="leaderboardStore.leaderboards.length"
        :user="sidebarUser"
        :mobile-visible="uiStore.showSidebarOnMobile"
        @switch-tournament="handleSwitchTournament"
        @toggle-theme="handleToggleTheme"
        @logout="handleLogout"
        @nav-click="uiStore.closeMobileSidebar"
      />
      <main class="app-main">
        <router-view />
      </main>
      <BottomTabBar />
    </div>
  </div>
</template>

<style>
@import './assets/css/variables.css';
@import './assets/css/global.css';

#app {
  min-height: 100vh;
  background: var(--color-bg);
}

/* ──── App Layout ──── */
.app-layout {
  display: grid;
  grid-template-columns: var(--sidebar-w) 1fr;
  min-height: 100vh;
}

.app-sidebar {
  position: sticky;
  top: 0;
  height: 100vh;
}

.app-main {
  min-height: 100vh;
  background: var(--color-bg);
}

/* Mobile: stack vertically, show bottom nav */
@media (max-width: 768px) {
  .app-layout {
    grid-template-columns: 1fr;
  }

  .app-sidebar {
    display: none;
  }

  .app-main {
    padding-top: var(--mobile-header-h);
    padding-bottom: var(--bottom-nav-h);
  }
}

/* Mobile sidebar backdrop */
@media (max-width: 768px) {
  .mobile-sidebar-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 98;
    animation: fadeIn var(--duration-fast) var(--ease-default);
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
}
</style>

