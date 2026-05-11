import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth.store.js'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import LeaderboardView from '../views/LeaderboardView.vue'
import MatchView from '../views/MatchView.vue'
import JoinView from '../views/JoinView.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView
  },
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: { requiresAuth: true }
  },
  {
    path: '/leaderboard/:id',
    name: 'Leaderboard',
    component: LeaderboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/match/:id',
    name: 'Match',
    component: MatchView,
    meta: { requiresAuth: true }
  },
  {
    path: '/join',
    name: 'Join',
    component: JoinView,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  if (to.path === '/auth/callback') {
    const token = to.query.token
    if (token) {
      const authStore = useAuthStore()

      try {
        const response = await fetch('/api/auth/me', {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        })

        if (response.ok) {
          const user = await response.json()
          authStore.setAuth(token, user)
        } else {
          authStore.setAuth(token, {})
        }
      } catch (err) {
        console.error('Failed to fetch user info:', err)
        authStore.setAuth(token, {})
      }

      next('/')
      return
    }
  }

  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
