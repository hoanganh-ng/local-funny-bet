<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './store/auth.store.js'
import SidebarNav from './components/common/SidebarNav.vue'
import BottomTabBar from './components/common/BottomTabBar.vue'

const route = useRoute()
const authStore = useAuthStore()

const showLayout = computed(() => {
  return route.meta.public !== true
})

onMounted(async () => {
  // Try to rehydrate auth from httpOnly cookie on app boot
  try {
    await authStore.refresh()
  } catch (err) {
    // Refresh failed - user not logged in, silent fail
    console.log('No valid session found')
  }
})
</script>

<template>
  <div id="app">
    <!-- Public routes: no layout -->
    <router-view v-if="!showLayout" />

    <!-- Authenticated routes: with sidebar/bottom nav -->
    <div v-else class="app-layout">
      <SidebarNav class="app-sidebar" />
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
    padding-bottom: var(--bottom-nav-h);
  }
}
</style>

